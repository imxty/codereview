package report

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ReportAPIHandler) ModifyReportRemark(ctx context.Context, req *pb.ModifyReportRemarkRequest, rsp *pb.ModifyReportRemarkResponse) error {

	// 1.验证request
	err := validateModifyReportRemarkRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取报告
	report, err := s.reportStore.GetReportByID(ctx, req.GetTenantId(), req.GetReportId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if report == nil {
		return errors.Error(ErrReportNotFound, "report not found")
	}

	// 修改备注
	err = s.reportStore.ModifyReportRemark(ctx, req.GetTenantId(), req.GetReportId(), req.GetRemark(), report.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证request
func validateModifyReportRemarkRequest(req *pb.ModifyReportRemarkRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetReportId() == "" {
		return gerr.New("report_id should not be empty")
	}
	return nil
}
