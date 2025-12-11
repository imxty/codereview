package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 搜索商户详情
func (u *UserAPIHandler) SearchTenantDetails(ctx context.Context, req *pb.SearchTenantDetailsRequest, rsp *pb.SearchTenantDetailsResponse) error {
	err := validateSearchTenantDetails(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取组织下所有商户 ID
	tids, err := u.userStore.ListAllOrganizationTenantIds(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 搜索商户
	tenants, totalCount, err := u.userStore.SearchTenantsByNameAndTime(ctx, tids, req.GetTenantName(), req.GetContactName(), req.GetContactPhone(), int32(req.GetStatus()), req.GetPagination().GetOffset(), req.GetPagination().GetSize(), req.GetStartTime(), req.GetEndTime())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 没有搜索到商户，返回空数据
	if req.GetTenantName() != "" && len(tenants) == 0 {
		rsp.Subscriptions = []*pb.SubscriptionPeriod{}
		rsp.Total = 0
		return nil
	}

	tenantIds := make([]string, len(tenants))
	for k, v := range tenants {
		tenantIds[k] = v.GetTenantID()
	}

	// 获取组织信息
	o, err := u.userStore.GetOrganizationByOrganizationId(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 搜索订阅信息
	subscriptions, _, err := u.userStore.BatchGetLatestTenantSubscriptionTimelinesDetails(ctx, tenantIds)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 整合数据
	subMap := make(map[string]domain.SubscriptionTimelineIntf)
	for _, each := range subscriptions {
		subMap[each.GetTenantID()] = each
	}

	results := make([]*pb.SubscriptionPeriod, len(tenants))
	for k, v := range tenants {
		sub, ok := subMap[v.GetTenantID()]
		var startTime *timestamppb.Timestamp
		if ok {
			startTime = timestamppb.New(sub.GetStartTime())
		}
		var endTime *timestamppb.Timestamp
		if ok {
			endTime = timestamppb.New(sub.GetEndTime())
		}
		var years int32
		if ok {
			years = sub.GetYears()
		}

		// 未认证，获取最新审核信息
		if v.GetSafePhone() == "" {
			lastRevision, err := u.userStore.GetLatestTenantEntityRevision(ctx, v.GetTenantID())
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
			results[k] = &pb.SubscriptionPeriod{
				TenantId:         v.GetTenantID(),
				CreatedAt:        startTime,
				OrganizationName: o.GetUsername(),
				TenantName:       lastRevision.GetStoreName(),
				ContactName:      lastRevision.GetContactName(),
				ContactPhone:     lastRevision.GetContactPhone(),
				Years:            -1,
				ExpiredAt:        endTime,
			}
		} else {
			results[k] = &pb.SubscriptionPeriod{
				TenantId:         v.GetTenantID(),
				CreatedAt:        startTime,
				OrganizationName: o.GetUsername(),
				TenantName:       v.GetStoreName(),
				ContactName:      v.GetContactName(),
				ContactPhone:     v.GetContactPhone(),
				Years:            years,
				ExpiredAt:        endTime,
			}
		}
	}

	rsp.Subscriptions = results
	rsp.Total = int32(totalCount)

	return nil
}

// 验证 request
func validateSearchTenantDetails(req *pb.SearchTenantDetailsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	if req.GetStatus() != pb.TenantStatus_TENANT_STATUS_USING && req.GetStatus() != pb.TenantStatus_TENANT_STATUS_PENDING && req.GetStatus() != pb.TenantStatus_TENANT_STATUS_UNSET {
		return gerr.New("status invalid")
	}
	return nil
}
