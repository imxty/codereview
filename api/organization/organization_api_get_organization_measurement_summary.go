package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	reportv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取组织测量概况
func (s *OrganizationAPIHandler) GetOrganizationMeasurementSummary(ctx context.Context, req *pb.GetOrganizationMeasurementSummaryRequest, rsp *pb.GetOrganizationMeasurementSummaryResponse) error {
	err := validateGetOrganizationMeasurementSummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取测量概况
	getRsp, err := s.reportAPI.GetOrganizationReportCount(ctx, &reportv1.GetOrganizationReportCountRequest{
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回响应
	rsp.TodayCustomerMeasurementCount = getRsp.GetTodayCustomerMeasurementCount()
	rsp.TodayTempCustomerMeasurementCount = getRsp.GetTodayTempCustomerMeasurementCount()
	rsp.YearCustomerMeasurementCount = getRsp.GetCustomerMeasurementYearCount()
	rsp.YearTempCustomerMeasurementCount = getRsp.GetTempCustomerMeasurementYearCount()

	return nil
}

// 验证request
func validateGetOrganizationMeasurementSummaryRequest(req *pb.GetOrganizationMeasurementSummaryRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
