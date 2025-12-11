package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取组织 vip 客户概况
func (s *OrganizationAPIHandler) GetOrganizationCustomerCompareSummary(ctx context.Context, req *pb.GetOrganizationCustomerCompareSummaryRequest, rsp *pb.GetOrganizationCustomerCompareSummaryResponse) error {
	err := validateGetOrganizationCustomerCompareSummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.customerAPI.GetOrganizationCustomerCompareSummary(ctx, &customerv1.GetOrganizationCustomerCompareSummaryRequest{
		OrganizationId: req.GetOrganizationId(),
		Date:           req.GetDate(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.MonthAddedCustomerCount = getRsp.GetMonthAddedCustomerCount()
	rsp.YearOnYear = getRsp.GetYearOnYear()
	rsp.MonthOnMonth = getRsp.GetMonthOnMonth()

	return nil
}

// 验证request
func validateGetOrganizationCustomerCompareSummaryRequest(req *pb.GetOrganizationCustomerCompareSummaryRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetDate() == nil {
		return gerr.New("date should not be empty")
	}
	return nil
}
