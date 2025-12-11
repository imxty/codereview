package report

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ReportAPIHandler) GetCustomerConstitution(ctx context.Context, req *pb.GetCustomerConstitutionRequest, rsp *pb.GetCustomerConstitutionResponse) error {
	err := validateGetCustomerConstitutionRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取常客最近测量的5笔报告
	reports, err := s.reportStore.ListCustomerLastReports(ctx, req.GetTenantId(), req.GetCustomerId(), ReportMinCount)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.IsReportAvailable = true

	// 报告数量不足返回错误
	if reports == nil || len(reports) < ReportMinCount {
		rsp.IsReportAvailable = false
		return nil
	}

	// 返回体质map
	summary, err := getConstitution(reports)
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}

	rsp.ConstitutionSummary = summary
	return nil
}

// 验证request
func validateGetCustomerConstitutionRequest(req *pb.GetCustomerConstitutionRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	if req.GetCustomerId() == "" {
		return gerr.New("customer id should not be empty")
	}
	return nil
}
