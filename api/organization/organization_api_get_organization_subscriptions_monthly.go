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

func (s *OrganizationAPIHandler) GetOrganizationSubscriptionsMonthly(ctx context.Context, req *pb.GetOrganizationSubscriptionsMonthlyRequest, rsp *pb.GetOrganizationSubscriptionsMonthlyResponse) error {
	err := validateGetOrganizationSubscriptionsMonthlyRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	if req.GetOrganizationId() == "" {
		rsp.Subscriptions = make([]*pb.TenantSubscriptionMonthly, 0)
		rsp.TotalCount = 0
		return nil
	}

	getRsp, err := s.userAPI.GetOrganizationSubscriptionsMonthly(ctx, &userv1.GetOrganizationSubscriptionsMonthlyRequest{
		OrganizationId: req.GetOrganizationId(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
		Pagination:     toSvcPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	results := make([]*pb.TenantSubscriptionMonthly, len(getRsp.GetSubscriptions()))
	for k, v := range getRsp.GetSubscriptions() {
		results[k] = &pb.TenantSubscriptionMonthly{
			Date:        v.GetDate(),
			TenantCount: v.GetTenantCount(),
			YearCount:   v.GetYearCount(),
		}
	}

	rsp.Subscriptions = results
	rsp.TotalCount = getRsp.GetTotalCount()

	return nil
}

func validateGetOrganizationSubscriptionsMonthlyRequest(req *pb.GetOrganizationSubscriptionsMonthlyRequest) error {
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
