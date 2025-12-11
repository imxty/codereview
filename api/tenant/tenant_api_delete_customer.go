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

func (s *TenantAPIHandler) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest, rsp *pb.DeleteCustomerResponse) error {
	err := validateDeleteCustomerRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.customerAPI.DeleteCustomer(ctx, &customerv1.DeleteCustomerRequest{
		CustomerId: req.GetCustomerId(),
		TenantId:   req.GetTenantId(),
	})
	if err != nil {
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
