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

func (s *AppAPIHandler) SearchCustomers(ctx context.Context, req *pb.SearchCustomersRequest, rsp *pb.SearchCustomersResponse) error {
	// 验证request
	err := validateSearchCustomersRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 发送请求 获取响应
	searchCustomersRsp, err := s.customerAPI.SearchCustomers(ctx, &customerpb.SearchCustomersRequest{
		TenantId:   req.GetTenantId(),
		Keywords:   req.GetKeywords(),
		Pagination: toSvcCustomerPagination(req.GetPagination()),
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据给app
	appCustomers := make([]*pb.Customer, len(searchCustomersRsp.GetCustomers()))
	for k, v := range searchCustomersRsp.GetCustomers() {
		appCustomers[k], err = toAppCustomer(v)
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
	}

	rsp.Customers = appCustomers
	return nil
}

// 验证request
func validateSearchCustomersRequest(req *pb.SearchCustomersRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}
