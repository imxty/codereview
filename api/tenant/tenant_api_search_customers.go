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

// 搜索常客
func (s *TenantAPIHandler) SearchCustomers(ctx context.Context, req *pb.SearchCustomersRequest, rsp *pb.SearchCustomersResponse) error {
	err := validateSearchCustomersRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchRsp, err := s.customerAPI.SearchCustomers(ctx, &customerv1.SearchCustomersRequest{
		TenantId:   req.GetTenantId(),
		Keywords:   req.GetKeywords(),
		Pagination: toSvcCustomerPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 转换
	customers := make([]*pb.Customer, len(searchRsp.GetCustomers()))
	for k, v := range searchRsp.GetCustomers() {
		customers[k] = toAppCustomer(v)
	}
	rsp.Customers = customers
	rsp.TotalCount = searchRsp.GetTotalCount()

	return nil
}

// 验证request
func validateSearchCustomersRequest(req *pb.SearchCustomersRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
