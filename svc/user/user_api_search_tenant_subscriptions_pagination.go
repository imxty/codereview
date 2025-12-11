package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 分页查询商户续费明细
func (u *UserAPIHandler) SearchTenantSubscriptionsPagination(ctx context.Context, req *pb.SearchTenantSubscriptionsPaginationRequest, rsp *pb.SearchTenantSubscriptionsPaginationResponse) error {
	err := validateSearchTenantSubscriptionsPaginationRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取组织下的所有商户（包括已删除的）
	tenant_ids, err := u.userStore.ListAllOrganizationTenantIdsIncludeDeleted(ctx, []string{req.GetOrganizationId()})
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 查询订阅明细
	periods, total_count, err := u.userStore.BatchGetTenantsSubscriptionPeriods(ctx, tenant_ids, req.GetStartTime().AsTime(), req.GetEndTime().AsTime(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 数据转换
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

	rsp.Subcriptions = results
	rsp.TotalCount = int32(total_count)

	return nil
}

// 验证 request
func validateSearchTenantSubscriptionsPaginationRequest(req *pb.SearchTenantSubscriptionsPaginationRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be empty")
	}
	if req.GetEndTime() == nil {
		return gerr.New("end_time should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}
