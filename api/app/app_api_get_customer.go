package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest, rsp *pb.GetCustomerResponse) error {
	// 验证request
	err := validateGetCustomerRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 发送请求 获取响应
	getCustomerRsp, err := s.customerAPI.GetCustomer(ctx, &customerpb.GetCustomerRequest{
		TenantId:   req.GetTenantId(),
		CustomerId: req.GetCustomerId(),
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据给app
	appCustomer, err := toAppCustomer(getCustomerRsp.GetCustomer())
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}

	rsp.Customer = appCustomer
	rsp.OverdueCount = getCustomerRsp.GetOverdueCount()
	return nil
}

// 验证request
func validateGetCustomerRequest(req *pb.GetCustomerRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomerId() == "" {
		return gerr.New("customer_id should not be empty")
	}
	return nil
}
