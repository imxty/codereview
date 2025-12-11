package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) CheckCustomerExist(ctx context.Context, req *pb.CheckCustomerExistRequest, rsp *pb.CheckCustomerExistResponse) error {
	err := validateCheckCustomerExistRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 发送请求
	checkRsp, err := s.customerAPI.CheckCustomerExist(ctx, &customerv1.CheckCustomerExistRequest{
		CustomerId: req.GetCustomerId(),
		TenantId:   req.GetTenantId(),
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.CustomerExisting = checkRsp.GetCustomerExisting()
	alCustomer, err := toAppCustomer(checkRsp.GetAlternativeCustomer())
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	cuCustomer, err := toAppCustomer(checkRsp.GetCurrentCustomer())
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	rsp.CurrentCustomer = cuCustomer
	rsp.AlternativeCustomer = alCustomer
	return nil
}

// 验证request
func validateCheckCustomerExistRequest(req *pb.CheckCustomerExistRequest) error {
	if req.GetCustomerId() == "" {
		return gerr.New("customer id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
