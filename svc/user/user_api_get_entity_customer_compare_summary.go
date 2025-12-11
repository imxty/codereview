package user

import (
	"context"
	gerr "errors"

	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取实体客户概况
func (u *UserAPIHandler) GetEntityCustomerCompareSummary(ctx context.Context, req *pb.GetEntityCustomerCompareSummaryRequest, rsp *pb.GetEntityCustomerCompareSummaryResponse) error {
	err := validateGetEntityCustomerCompareSummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := u.customerAPI.GetEntityCustomerCompareSummary(ctx, &customerv1.GetEntityCustomerCompareSummaryRequest{
		TenantId: req.GetTenantId(),
		Date:     req.GetDate(),
	})
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.MonthAddedCustomerCount = getRsp.GetMonthAddedCustomerCount()
	rsp.YearOnYear = getRsp.GetYearOnYear()
	rsp.MonthOnMonth = getRsp.GetMonthOnMonth()

	return nil
}

// 验证 request
func validateGetEntityCustomerCompareSummaryRequest(req *pb.GetEntityCustomerCompareSummaryRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetDate() == nil {
		return gerr.New("date should not be empty")
	}
	return nil
}
