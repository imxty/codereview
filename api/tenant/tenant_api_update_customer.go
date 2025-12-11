package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest, rsp *pb.UpdateCustomerResponse) error {
	err := validateUpdateCustomerRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	_, err = s.customerAPI.UpdateCustomer(ctx, &customerpb.UpdateCustomerRequest{
		// 租户ID
		TenantId: req.GetTenantId(),
		// 常客信息
		Customer: toSvcCustomer(req.GetCustomer()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateUpdateCustomerRequest(req *pb.UpdateCustomerRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomer().GetCustomerId() == "" {
		return gerr.New("customer_id should not be empty")
	}
	return nil
}
