package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
)

// 复查人员列表
func (s *AppAPIHandler) SearchCustomerLastStatus(ctx context.Context, req *pb.SearchCustomerLastStatusRequest, rsp *pb.SearchCustomerLastStatusResponse) error {

	// 1.验证request
	err := validateSearchCustomerLastStatusRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	searchRsp, err := s.customerAPI.SearchAppCustomerLastStatus(ctx, &customerv1.SearchAppCustomerLastStatusRequest{
		TenantId:   req.GetTenantId(),
		Key:        req.GetKey(),
		Pagination: toSvcCustomerPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	results := make([]*pb.CustomerLastStatus, len(searchRsp.GetCustomers()))
	for k, v := range searchRsp.GetCustomers() {
		results[k] = toAppCustomerLastStatus(v)
	}

	rsp.Customers = results
	rsp.TotalCount = searchRsp.GetTotalCount()

	return nil
}

// 验证request
func validateSearchCustomerLastStatusRequest(req *pb.SearchCustomerLastStatusRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}
