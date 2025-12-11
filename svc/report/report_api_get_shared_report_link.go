package report

import (
	"context"
	gerr "errors"
	"fmt"

	"github.com/google/uuid"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// GetSharedReportLink 获取分享报告连接
func (s *ReportAPIHandler) GetSharedReportLink(ctx context.Context, req *pb.GetSharedReportLinkRequest, rsp *pb.GetSharedReportLinkResponse) error {

	// 1.验证request
	err := validateGetSharedReportLinkRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 1. 获取报告
	report, err := s.reportStore.GetReportByID(ctx, req.GetTenantId(), req.GetReportId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if report == nil {
		return errors.Errorf(ErrReportNotFound, "report not found by report_id %s", req.GetReportId())
	}
	// 生成报告的token(uuid)
	// 生成uuid
	token, err := uuid.NewUUID()
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}
	shareToken := &domain.SharedReport{
		Token:    token.String(),
		ReportID: req.GetReportId(),
	}
	// 保存生成的报告token
	err = s.reportStore.CreateSharedReport(ctx, shareToken)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回连接
	rsp.ReportLink = getReportLink(s.reportLink, token.String(), req.GetTenantId())

	return nil
}

func getReportLink(baseUrl, token, tid string) string {
	return fmt.Sprintf("%s%s?t=%s", baseUrl, token, tid)
}

// 验证request
func validateGetSharedReportLinkRequest(req *pb.GetSharedReportLinkRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetReportId() == "" {
		return gerr.New("report_id should not be empty")
	}
	return nil
}
