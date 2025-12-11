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
func (s *OrganizationAPIHandler) GetOrganizationMonthlyReportCount(ctx context.Context, req *pb.GetOrganizationMonthlyReportCountRequest, rsp *pb.GetOrganizationMonthlyReportCountResponse) error {
	err := validateGetOrganizationMonthlyReportCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取测量概况
	getRsp, err := s.reportAPI.GetOrganizationMonthlyReportCount(ctx, &reportv1.GetOrganizationMonthlyReportCountRequest{
		OrganizationId:  req.GetOrganizationId(),
		MeasurementTime: toReportDate(req.GetMeasurementTime()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.MonthOnMonth = getRsp.GetMonthOnMonth()
	rsp.YearOnYear = getRsp.GetYearOnYear()
	rsp.MonthlyCustomerMeasurementCount = getRsp.GetMonthlyCustomerMeasurementCount()

	return nil
}

// 验证request
func validateGetOrganizationMonthlyReportCountRequest(req *pb.GetOrganizationMonthlyReportCountRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetMeasurementTime() == nil {
		return gerr.New("measurement time should not be empty")
	}
	return nil
}
