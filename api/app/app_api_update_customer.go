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

func (s *AppAPIHandler) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest, rsp *pb.UpdateCustomerResponse) error {

	// 验证request
	err := validateUpdateCustomerRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 数据转化为Svc的数据
	svcCustomer, err := toSvcCustomer(req.GetCustomer())
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}

	// 发送请求 获取响应
	_, err = s.customerAPI.UpdateCustomer(ctx, &customerpb.UpdateCustomerRequest{
		TenantId: req.GetTenantId(),
		Customer: svcCustomer,
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateUpdateCustomerRequest(req *pb.UpdateCustomerRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomer() == nil {
		return gerr.New("customer should not be nil")
	}
	return nil
}
