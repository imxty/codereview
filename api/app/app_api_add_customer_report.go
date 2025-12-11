package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) AddCustomerReport(ctx context.Context, req *pb.AddCustomerReportRequest, rsp *pb.AddCustomerReportResponse) error {
	err := validateAddCustomerReportRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	customer := toReportSvcCustomer(req.GetCustomer())

	// 添加常客报告请求
	addRsp, err := s.reportAPI.AddCustomerReport(ctx, &reportpb.AddCustomerReportRequest{
		TenantId: req.GetTenantId(),
		ReportId: req.GetReportId(),
		Customer: customer,
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.CustomerId = addRsp.GetCustomerId()
	return nil
}

// 验证request
func validateAddCustomerReportRequest(req *pb.AddCustomerReportRequest) error {
	if req.GetCustomer() == nil {
		return gerr.New("customer should not be nil")
	}
	if req.GetReportId() == "" {
		return gerr.New("report id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
