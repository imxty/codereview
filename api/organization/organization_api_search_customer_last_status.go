package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 查询客户最新测量状态
func (s *OrganizationAPIHandler) SearchCustomerLastStatus(ctx context.Context, req *pb.SearchCustomerLastStatusRequest, rsp *pb.SearchCustomerLastStatusResponse) error {
	err := validateSearchCustomerLastStatusRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchRsp, err := s.customerAPI.SearchCustomerLastStatus(ctx, &customerv1.SearchCustomerLastStatusRequest{
		OrganizationId: req.GetOrganizationId(),
		TenantName:     req.GetTenantName(),
		CustomerName:   req.GetCustomerName(),
		CustomerPhone:  req.GetCustomerPhone(),
		Pagination:     toCustomerPagination(req.GetPagination()),
		OverdueMask:    req.GetOverdueMask(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	results := make([]*pb.CustomerLastStatus, len(searchRsp.GetCustomers()))
	for k, v := range searchRsp.GetCustomers() {
		results[k] = &pb.CustomerLastStatus{
			TenantName:    v.GetTenantName(),
			CustomerName:  v.GetCustomerName(),
			CustomerPhone: v.GetCustomerPhone(),
			OverdueCount:  v.GetOverdueCount(),
		}
	}

	rsp.Customers = results
	rsp.TotalCount = searchRsp.GetTotalCount()

	return nil
}

func validateSearchCustomerLastStatusRequest(req *pb.SearchCustomerLastStatusRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination ")
	}
	return nil
}
