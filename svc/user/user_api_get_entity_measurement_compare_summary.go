package user

import (
	"context"
	gerr "errors"

	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取实体测量概况
func (u *UserAPIHandler) GetEntityMeasurementCompareSummary(ctx context.Context, req *pb.GetEntityMeasurementCompareSummaryRequest, rsp *pb.GetEntityMeasurementCompareSummaryResponse) error {
	err := validateGetEntityMeasurementCompareSummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 调用 report 接口
	getRsp, err := u.reportAPI.GetTenantMonthlyReportCount(ctx, &reportpb.GetTenantMonthlyReportCountRequest{
		TenantId: req.GetTenantId(),
		MeasurementTime: &reportpb.Date{
			Year:  int32(req.GetDate().AsTime().Year()),
			Month: int32(req.GetDate().AsTime().Month()),
			Day:   1,
		},
	})
	if err != nil {
		return err
	}

	// 返回响应
	rsp.MonthMeasurementCount = getRsp.GetMonthlyCustomerMeasurementCount()
	rsp.YearOnYear = getRsp.GetYearOnYear()
	rsp.MonthOnMonth = getRsp.GetMonthOnMonth()

	return nil
}

// 验证 request
func validateGetEntityMeasurementCompareSummaryRequest(req *pb.GetEntityMeasurementCompareSummaryRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetDate() == nil {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
