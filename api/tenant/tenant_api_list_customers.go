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

func (s *TenantAPIHandler) ListCustomers(ctx context.Context, req *pb.ListCustomersRequest, rsp *pb.ListCustomersResponse) error {
	err := validateListCustomersRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	getGetCustomerRsp, err := s.customerAPI.ListCustomers(ctx, &customerpb.ListCustomersRequest{
		// 租户ID
		TenantId: req.GetTenantId(),
		// 常客ID
		Pagination: toSvcCustomerPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	// 转换
	customers := make([]*pb.Customer, len(getGetCustomerRsp.GetCustomers()))
	for i, v := range getGetCustomerRsp.GetCustomers() {
		customers[i] = toAppCustomer(v)
	}
	rsp.Customers = customers
	rsp.TotalCount = getGetCustomerRsp.GetTotalCount()
	return nil
}

// 验证request
func validateListCustomersRequest(req *pb.ListCustomersRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	if req.GetPagination().GetSize() <= 0 {
		return gerr.New("invalid pagination")
	}
	return nil
}
