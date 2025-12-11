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

// 获取组织商户概况
func (s *OrganizationAPIHandler) GetOrganizationTenantSummary(ctx context.Context, req *pb.GetOrganizationTenantSummaryRequest, rsp *pb.GetOrganizationTenantSummaryResponse) error {
	err := validateGetOrganizationTenantSummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetOrganizationTenantSummary(ctx, &userv1.GetOrganizationTenantSummaryRequest{
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	overdueRsp, err := s.userAPI.GetOrganizationOverdueSummary(ctx, &userv1.GetOrganizationOverdueSummaryRequest{
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.TotalTenantAmounts = getRsp.GetTenantTotalCount()
	rsp.UsingTeantAmounts = getRsp.GetTenantUsingCount()
	rsp.TodayAddedTenantCount = getRsp.GetTodayAddedTenantCount()
	rsp.TenantsExpiredAmounts = getRsp.GetTenantExpiringSoonCount()
	rsp.CustomerOverdueCount = overdueRsp.GetOverdueCustomerTotalCount()

	return nil
}

// 验证request
func validateGetOrganizationTenantSummaryRequest(req *pb.GetOrganizationTenantSummaryRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
