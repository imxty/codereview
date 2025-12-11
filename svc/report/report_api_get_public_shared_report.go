package report

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrTokenNotFound token不存在
	ErrTokenNotFound = 5024
)

func (s *ReportAPIHandler) GetPublicSharedReport(ctx context.Context, req *pb.GetPublicSharedReportRequest, rsp *pb.GetPublicSharedReportResponse) error {

	// 1.验证request
	err := validateGetPublicSharedReportRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 通过token获取报告ID
	sharedReport, err := s.reportStore.GetReportIDByToken(ctx, req.GetTenantId(), req.GetToken())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// token不存在
	if sharedReport == nil {
		return errors.Error(ErrTokenNotFound, "token not found")
	}

	getReportRsp := &pb.GetReportResponse{}
	err = s.GetReport(ctx, &pb.GetReportRequest{
		TenantId: req.GetTenantId(),
		ReportId: sharedReport.GetReportID(),
	}, getReportRsp)
	// 直接返回service层返回的错误即可，不需要二次包装
	if err != nil {
		return err
	}
	// 如果报告还没有完成就报错
	if !getReportRsp.GetIsCompleteReport() {
		return errors.Error(codes.InvalidRequest, "report has not complete")
	}

	// 返回报告
	rsp.Report = getReportRsp.GetReport()
	rsp.TongueFaceReport = getReportRsp.GetTongueFaceReport()
	// 员工备注不显示
	rsp.Report.StaffRemarks = ""

	return nil
}

// 验证request
func validateGetPublicSharedReportRequest(req *pb.GetPublicSharedReportRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetToken() == "" {
		return gerr.New("token should not be empty")
	}
	return nil
}
