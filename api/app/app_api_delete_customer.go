package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
)

func (s *AppAPIHandler) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest, rsp *pb.DeleteCustomerResponse) error {

	// 验证request
	err := validateDeleteCustomerRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 发送请求 获取响应
	_, err = s.customerAPI.DeleteCustomer(ctx, &customerpb.DeleteCustomerRequest{
		TenantId:   req.GetTenantId(),
		CustomerId: req.GetCustomerId(),
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

// 验证request
func validateDeleteCustomerRequest(req *pb.DeleteCustomerRequest) error {
	if req.GetCustomerId() == "" {
		return gerr.New("customer_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
