package report

import (
	"context"
	gerr "errors"
	"strconv"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 添加常客报告请求
func (s *ReportAPIHandler) AddCustomerReport(ctx context.Context, req *pb.AddCustomerReportRequest, rsp *pb.AddCustomerReportResponse) error {

	// 1.验证request
	err := validateAddCustomerReportRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取报告
	report, err := s.reportStore.GetReportByID(ctx, req.GetTenantId(), req.GetReportId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if report == nil {
		return errors.Errorf(ErrReportNotFound, "report not found by report_id %s", req.GetReportId())
	}
	// 如果报告已经是常客报告则返回错误
	if report.GetIsCustomer() {
		return errors.Error(codes.InvalidRequest, "current report is already a customer-report")
	}

	// 添加常客信息
	addRsp, err := s.customerAPI.AddCustomer(ctx, &customerpb.AddCustomerRequest{
		TenantId: req.GetTenantId(),
		Customer: toCustomer(req.GetCustomer()),
	})
	if err != nil {
		return err
	}
	customer := addRsp.GetCustomer()
	// 将散客报告变为常客的报告
	age := strconv.Itoa(int(customer.GetAge()))
	// 将报告变为常客报告
	err = s.reportStore.AddCustomerReport(ctx, req.GetTenantId(), req.GetReportId(), customer.GetCustomerId(), age, report.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.CustomerId = customer.GetCustomerId()
	return nil
}

// 验证request
func validateAddCustomerReportRequest(req *pb.AddCustomerReportRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetReportId() == "" {
		return gerr.New("report_id should not be empty")
	}
	if req.GetCustomer() == nil {
		return gerr.New("customer should not be nil")
	}
	customer := req.GetCustomer()
	if customer.GetNickname() == "" {
		return gerr.New("nickname should not be empty")
	}
	if customer.GetBirthday() == nil {
		return gerr.New("birthday should not be empty")
	}
	return nil
}
