package user

import (
	"context"
	gerr "errors"

	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取组织 vip 客户概况
func (u *UserAPIHandler) GetOrganizationCustomerCompareSummary(ctx context.Context, req *pb.GetOrganizationCustomerCompareSummaryRequest, rsp *pb.GetOrganizationCustomerCompareSummaryResponse) error {
	err := validateGetOrganizationCustomerCompareSummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := u.customerAPI.GetOrganizationCustomerCompareSummary(ctx, &customerv1.GetOrganizationCustomerCompareSummaryRequest{
		OrganizationId: req.GetOrganizationId(),
		Date:           req.GetDate(),
	})
	if err != nil {
		return err
	}

	rsp.MonthAddedCustomerCount = getRsp.GetMonthAddedCustomerCount()
	rsp.YearOnYear = getRsp.GetYearOnYear()
	rsp.MonthOnMonth = getRsp.GetMonthOnMonth()

	return nil
}

// 验证 request
func validateGetOrganizationCustomerCompareSummaryRequest(req *pb.GetOrganizationCustomerCompareSummaryRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetDate() == nil {
		return gerr.New("date should not be nil")
	}
	return nil
}
