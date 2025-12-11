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

func (s *AppAPIHandler) ListCustomers(ctx context.Context, req *pb.ListCustomersRequest, rsp *pb.ListCustomersResponse) error {

	// 验证request
	err := validateListCustomersRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 发送请求 获取响应
	getCustomersRsp, err := s.customerAPI.ListCustomers(ctx, &customerpb.ListCustomersRequest{
		TenantId:   req.GetTenantId(),
		Pagination: toSvcCustomerPagination(req.GetPagination()),
		StaffId:    req.GetStaffId(),
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据给app
	appCustomers := make([]*pb.Customer, len(getCustomersRsp.GetCustomers()))
	for k, v := range getCustomersRsp.GetCustomers() {
		appCustomers[k], err = toAppCustomer(v)
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
	}

	rsp.Customers = appCustomers
	rsp.TotalCount = getCustomersRsp.GetTotalCount()

	return nil
}

// 验证request
func validateListCustomersRequest(req *pb.ListCustomersRequest) error {
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	if req.GetPagination().GetSize() <= 0 {
		return gerr.New("invalid pagination")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
