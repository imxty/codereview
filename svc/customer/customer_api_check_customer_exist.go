package customer

import (
	"context"
	gerr "errors"

	"github.com/jinmukeji/go-pkg/v2/age"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *CustomerAPIHandler) CheckCustomerExist(ctx context.Context, req *pb.CheckCustomerExistRequest, rsp *pb.CheckCustomerExistResponse) error {
	err := validateCheckCustomerExistRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 查询常客
	customer, err := s.customerStore.GetCustomer(ctx, req.GetTenantId(), req.GetCustomerId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if customer == nil {
		rsp.CustomerExisting = false
		// 如果不存在则去获取第一位常客
		customers, err := s.customerStore.ListCustomers(ctx, req.GetTenantId(), 0, 1)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		if len(customers) == 0 {
			// 直接return即可
			return nil
		}
		rsp.AlternativeCustomer = toProtoCustomer(customers[0], customers[0].GetInitial(), int32(age.Age(customers[0].GetBirthday())))
		return nil
	}
	// 如果存在直接返回即可
	rsp.CustomerExisting = true
	rsp.CurrentCustomer = toProtoCustomer(customer, getNicknameInitial(customer.GetNickname()), int32(age.Age(customer.GetBirthday())))
	return nil
}

// 验证request
func validateCheckCustomerExistRequest(req *pb.CheckCustomerExistRequest) error {
	if req.GetCustomerId() == "" {
		return gerr.New("customer should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
