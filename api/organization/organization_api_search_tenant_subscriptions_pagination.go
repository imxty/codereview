package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 分页查询商户续费明细
func (s *OrganizationAPIHandler) SearchTenantSubscriptionsPagination(ctx context.Context, req *pb.SearchTenantSubscriptionsPaginationRequest, rsp *pb.SearchTenantSubscriptionsPaginationResponse) error {
	err := validateSearchTenantSubscriptionsPaginationRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchRsp, err := s.userAPI.SearchTenantSubscriptionsPagination(ctx, &userv1.SearchTenantSubscriptionsPaginationRequest{
		OrganizationId: req.GetOrganizationId(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
		Pagination:     toSvcPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	results := make([]*pb.SubscriptionPeriod, len(searchRsp.GetSubcriptions()))
	for k, v := range searchRsp.GetSubcriptions() {
		results[k] = &pb.SubscriptionPeriod{
			TenantId:         v.GetTenantId(),
			CreatedAt:        v.GetCreatedAt(),
			OrganizationName: v.GetOrganizationName(),
			TenantName:       v.GetTenantName(),
			ContactPhone:     v.GetContactPhone(),
			Years:            v.GetYears(),
			ExpiredAt:        v.GetExpiredAt(),
		}
	}

	rsp.Subcriptions = results
	rsp.TotalCount = searchRsp.GetTotalCount()

	return nil
}

// 验证request
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
		return gerr.New("pagination should not be nil")
	}
	return nil
}
