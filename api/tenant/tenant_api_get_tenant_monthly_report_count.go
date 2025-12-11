package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	reportv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取商户每月报告统计
func (s *TenantAPIHandler) GetTenantMonthlyReportCount(ctx context.Context, req *pb.GetTenantMonthlyReportCountRequest, rsp *pb.GetTenantMonthlyReportCountResponse) error {
	err := validateGetTenantMonthlyReportCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.reportAPI.GetTenantMonthlyReportCount(ctx, &reportv1.GetTenantMonthlyReportCountRequest{
		TenantId:        req.GetTenantId(),
		MeasurementTime: toReportDate(req.GetMeasurementTime()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.MonthlyCustomerMeasurementCount = getRsp.GetMonthlyCustomerMeasurementCount()
	rsp.MonthOnMonth = getRsp.GetMonthOnMonth()
	rsp.YearOnYear = getRsp.GetYearOnYear()

	return nil
}

// 验证request
func validateGetTenantMonthlyReportCountRequest(req *pb.GetTenantMonthlyReportCountRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetMeasurementTime() == nil {
		return gerr.New("measurement_time should not be nil")
	}
	return nil
}
