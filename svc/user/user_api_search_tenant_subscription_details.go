package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/jinmukeji/plat-pkg/v4/micro/meta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 查询商户续期详情
func (u *UserAPIHandler) SearchTenantSubscriptionDetails(ctx context.Context, req *pb.SearchTenantSubscriptionDetailsRequest, rsp *pb.SearchTenantSubscriptionDetailsResponse) error {
	err := validateSearchTenantSubscriptionDetailsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 如果 organization_id 是空值，查询所有组织
	// 如果不是空，只查询该组织
	var oids []string
	if req.GetOrganizationId() == "" {
		// 获取权限组
		user_id, ok := meta.Get(ctx, "user_id")
		if ok {
			user, err := u.userStore.GetSystemUserByID(ctx, user_id)
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
			if user == nil {
				return errors.Errorf(ErrSystemUserNotFound, "user [%s] not found", user_id)
			}
			datas, err := u.userStore.GetDataPrivilegeByID(ctx, user.GetPrivilegeID())
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
			if len(datas) == 1 && datas[0] == "all" {
				// 查询所有组织
				organizations, err := u.userStore.ListAllOrganizations(ctx)
				if err != nil {
					return errors.Error(codes.DataAccessFailed, err.Error())
				}
				for _, o := range organizations {
					oids = append(oids, o.GetOrganizationID())
				}
			} else {
				oids = datas
			}
		} else {
			// 查询所有组织
			organizations, err := u.userStore.ListAllOrganizations(ctx)
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
			for _, o := range organizations {
				oids = append(oids, o.GetOrganizationID())
			}
		}
	} else {
		oids = []string{req.GetOrganizationId()}
	}

	if len(oids) < 1 {
		rsp.TotalYears = 0
		rsp.Subscriptions = []*pb.SubscriptionPeriod{}
		rsp.Total = 0
		return nil
	}

	// 1. 如果商户 ID 为空，查询组织下所有商户
	// 2. 如果商户 ID 不为空，只查询该商户
	var tenant_ids []string
	if req.GetTenantId() == "" {
		tenant_ids, err = u.userStore.ListAllOrganizationTenantIdsIncludeDeleted(ctx, oids)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	} else {
		tenant_ids = []string{req.GetTenantId()}
	}

	// 如果商户id为空，直接返回空数据
	if len(tenant_ids) < 1 {
		rsp.TotalYears = 0
		rsp.Subscriptions = []*pb.SubscriptionPeriod{}
		rsp.Total = 0
		return nil
	}

	// 获取订阅信息
	periods, total_count, total_years, err := u.userStore.BatchGetTenantsSubscriptionPeriodsDetails(ctx, tenant_ids, req.GetContactName(), req.GetContactPhone(), req.GetStartTime().AsTime(), req.GetEndTime().AsTime(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 整合数据
	results := make([]*pb.SubscriptionPeriod, len(periods))
	for k, v := range periods {
		results[k] = &pb.SubscriptionPeriod{
			TenantId:         v.GetTenantID(),
			CreatedAt:        timestamppb.New(v.GetCreatedAt()),
			OrganizationName: v.GetOrganizationName(),
			TenantName:       v.GetTenantName(),
			ContactName:      v.GetContactName(),
			ContactPhone:     v.GetContactPhone(),
			Years:            v.GetYears(),
			ExpiredAt:        timestamppb.New(v.GetExpiredTime()),
		}
	}

	rsp.Subscriptions = results
	rsp.Total = int32(total_count)
	rsp.TotalYears = total_years

	return nil
}

// 验证 request
func validateSearchTenantSubscriptionDetailsRequest(req *pb.SearchTenantSubscriptionDetailsRequest) error {
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be nil")
	}
	if req.GetEndTime() == nil {
		return gerr.New("end_time should not be nil")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
