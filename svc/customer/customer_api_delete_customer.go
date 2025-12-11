package customer

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrCustomerNotFound 常客不存在
	ErrCustomerNotFound = 5035
)

func (s *CustomerAPIHandler) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest, rsp *pb.DeleteCustomerResponse) error {
	// 验证request
	err := validateDeleteCustomerRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取当前常客版本号
	customer, err := s.customerStore.GetCustomer(ctx, req.GetTenantId(), req.GetCustomerId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if customer == nil {
		return errors.Errorf(ErrCustomerNotFound, "customer[%s] not found", req.GetCustomerId())
	}

	rev := customer.GetRev()

	// 删除常客
	err = s.customerStore.DeleteCustomer(ctx, req.GetTenantId(), req.GetCustomerId(), rev)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, "delete customer failed")
	}
	return nil
}

// 验证request
func validateDeleteCustomerRequest(req *pb.DeleteCustomerRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetCustomerId() == "" {
		return gerr.New("customer_id should not be empty")
	}
	return nil
}
