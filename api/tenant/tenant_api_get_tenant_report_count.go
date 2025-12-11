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

// 获取商户报告统计
func (s *TenantAPIHandler) GetTenantReportCount(ctx context.Context, req *pb.GetTenantReportCountRequest, rsp *pb.GetTenantReportCountResponse) error {
	err := validateGetTenantReportCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取报告统计
	getRsp, err := s.reportAPI.GetTenantReportCount(ctx, &reportv1.GetTenantReportCountRequest{
		TenantId: req.GetTenantId(),
		StaffIds: req.GetStaffIds(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	rsp.TodayCustomerMeasurementCount = getRsp.GetTodayCustomerMeasurementCount()
	rsp.TodayTempCustomerMeasurementCount = getRsp.GetTodayTempCustomerMeasurementCount()
	rsp.CustomerMeasurementYearCount = getRsp.GetCustomerMeasurementYearCount()
	rsp.TempCustomerMeasurementYearCount = getRsp.GetTempCustomerMeasurementYearCount()
	rsp.StaffCustomerCount = rsp.GetStaffCustomerCount()

	return nil
}

// 验证request
func validateGetTenantReportCountRequest(req *pb.GetTenantReportCountRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
