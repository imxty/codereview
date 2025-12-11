package user

import (
	"context"
	gerr "errors"
	"slices"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

const (
	YearDays = 365
)

// 提交订阅
func (u *UserAPIHandler) SubmitSubscription(ctx context.Context, req *pb.SubmitSubscriptionRequest, rsp *pb.SubmitSubscriptionResponse) error {
	err := validateSubmitSubscriptionRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商户信息
	tenants, err := u.userStore.BatchGetTenantNamesByIds(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 查询结果数量和输入不一致，说明有错误商户 ID
	if len(tenants) != len(req.GetTenantId()) {
		return errors.Error(ErrTenantNotFound, "has tenant not found")
	}

	// 获取组织信息
	organization_tenants, err := u.userStore.BatchGetOrganizationTenants(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	for _, each := range organization_tenants {
		if !each.GetIsActivated() {
			return errors.Errorf(ErrTenantUnAuth, "tenant[%s] unauth", each.GetTenantID())
		}
	}

	organization_tenants_map := make(map[string]string)
	oids := make([]string, len(organization_tenants))
	for k, v := range organization_tenants {
		oids[k] = v.GetOrganizationID()
		organization_tenants_map[v.GetTenantID()] = v.GetOrganizationID()
	}
	organizations, err := u.userStore.BatchGetOrganizationNamesByIDs(ctx, oids)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	organizations_map := make(map[string]string)
	for _, each := range organizations {
		organizations_map[each.GetOrganizationID()] = each.GetUsername()
	}

	// 获取商户时间线信息
	ts, err := u.userStore.ListSubscriptionTimeline(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 有时间线的商户
	ts_map := make(map[string]domain.SubscriptionTimelineIntf)
	tenantsHasTimeline := make([]string, 0, len(ts))
	for _, v := range ts {
		tenantsHasTimeline = append(tenantsHasTimeline, v.GetTenantID())
		ts_map[v.GetTenantID()] = v
	}
	// 没有时间线的商户
	tenantsNoTimeline := make([]string, 0, len(tenants)-len(ts))
	for _, v := range tenants {
		if !slices.Contains(tenantsHasTimeline, v.GetTenantID()) {
			tenantsNoTimeline = append(tenantsNoTimeline, v.GetTenantID())
		}
	}

	// 构建商户周期
	periods := make([]domain.SubscriptionPeriodIntf, len(tenants))
	var timeNow, endTime time.Time
	for k, v := range tenants {
		// 如果没有时间线
		if slices.Contains(tenantsNoTimeline, v.GetTenantID()) {
			timeNow = time.Now()
			endTime = timeNow.AddDate(0, 0, int(YearDays*req.GetTimeCount()))
		} else {
			// 有时间线且过期，重新算时间线
			if time.Now().After(ts_map[v.GetTenantID()].GetEndTime()) {
				timeNow = time.Now()
				endTime = timeNow.AddDate(0, 0, int(YearDays*req.GetTimeCount()))
			} else {
				// 有时间线且没有过期，增加时间线
				timeNow = ts_map[v.GetTenantID()].GetEndTime()
				endTime = timeNow.AddDate(0, 0, int(YearDays*req.GetTimeCount()))
			}
		}
		oid := organization_tenants_map[v.GetTenantID()]
		periods[k] = &domain.SubscriptionPeriod{
			SubscriptionPeriodID: xid.New().String(),
			TenantID:             v.GetTenantID(),
			OrganizationID:       oid,
			OrganizationName:     organizations_map[oid],
			TenantName:           v.GetStoreName(),
			ExpiredTime:          endTime,
			UserID:               req.GetUserId(),
			ContactName:          v.GetContactName(),
			ContactPhone:         v.GetContactPhone(),
			Years:                req.GetTimeCount(),
			Rev:                  0,
			CreatedAt:            time.Now(),
		}
	}

	// 构建时间线
	timelines := make([]domain.SubscriptionTimelineIntf, len(tenantsNoTimeline))
	for k, v := range tenantsNoTimeline {
		timelines[k] = &domain.SubscriptionTimeline{
			SubscriptionTimelineID: xid.New().String(),
			TenantID:               v,
			StartTime:              timeNow,
			EndTime:                endTime,
			Rev:                    0,
			Years:                  req.GetTimeCount(),
		}
	}

	ctx = u.userStore.BeginTx(ctx)
	// 创建 subscription_period 记录
	err = u.userStore.BatchCreateTenantSubscriptionPeriod(ctx, periods)
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 对所有已经存在时间线且已过期的商户做修改
	if len(tenantsHasTimeline) > 0 {
		tenantsHasExpired := make([]string, 0)
		for _, v := range tenantsHasTimeline {
			oldTimeline := ts_map[v]
			if time.Now().After(oldTimeline.GetEndTime().UTC()) {
				tenantsHasExpired = append(tenantsHasExpired, v)
			}
		}
		err = u.userStore.BatchAddTenantTimelines(ctx, tenantsHasExpired, endTime, req.GetTimeCount())
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}
	// 对所有已经存在时间线但未过期的商户做修改
	if len(tenantsHasTimeline) > 0 {
		tenantsHasNotExpired := make([]string, 0)
		for _, v := range tenantsHasTimeline {
			oldTimeline := ts_map[v]
			if !time.Now().After(oldTimeline.GetEndTime().UTC()) {
				tenantsHasNotExpired = append(tenantsHasNotExpired, v)
			}
		}
		err = u.userStore.BatchUpdateTenantTimelines(ctx, tenantsHasNotExpired, req.GetTimeCount(), req.GetTimeCount())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}
	// 为所有不存在时间线的商户创建时间线
	if len(tenantsNoTimeline) > 0 {
		err = u.userStore.BatchCreateTenantSubscriptionTimeline(ctx, timelines)
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	u.userStore.CommitTx(ctx)

	return nil
}

// 验证 request
func validateSubmitSubscriptionRequest(req *pb.SubmitSubscriptionRequest) error {
	if req.GetTenantId() == nil {
		return gerr.New("tenant_id should not be nil")
	}
	if len(req.GetTenantId()) < 1 {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetTimeCount() < 0 {
		return gerr.New("invalid time_count")
	}
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	return nil
}
