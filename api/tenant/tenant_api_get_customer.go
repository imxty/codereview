package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest, rsp *pb.GetCustomerResponse) error {
	err := validateGetCustomerRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 请求常客数据
	getRsp, err := s.customerAPI.GetCustomer(ctx, &customerv1.GetCustomerRequest{
		CustomerId: req.GetCustomerId(),
		TenantId:   req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.Customer = toAppCustomer(getRsp.GetCustomer())
	return nil
}

// 验证request
func validateGetCustomerRequest(req *pb.GetCustomerRequest) error {
	if req.GetCustomerId() == "" {
		return gerr.New("customer id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
