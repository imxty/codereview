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

func (s *AppAPIHandler) GetCustomerConstitution(ctx context.Context, req *pb.GetCustomerConstitutionRequest, rsp *pb.GetCustomerConstitutionResponse) error {
	err := validateGetCustomerConstitutionRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 发送获取常客的体质报告请求
	getRsp, err := s.reportAPI.GetCustomerConstitution(ctx, &reportpb.GetCustomerConstitutionRequest{
		TenantId:   req.GetTenantId(),
		CustomerId: req.GetCustomerId(),
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.IsReportAvailable = getRsp.GetIsReportAvailable()
	rsp.ConstitutionSummary = getRsp.GetConstitutionSummary()
	return nil
}

// 验证request
func validateGetCustomerConstitutionRequest(req *pb.GetCustomerConstitutionRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomerId() == "" {
		return gerr.New("customer_id should not be empty")
	}
	return nil
}
