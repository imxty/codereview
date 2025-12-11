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

// 获取待复查客户列表
func (s *TenantAPIHandler) SearchOverdueCustomers(ctx context.Context, req *pb.SearchOverdueCustomersRequest, rsp *pb.SearchOverdueCustomersResponse) error {
	err := validateSearchOverdueCustomersRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchRsp, err := s.customerAPI.SearchOverdueCustomers(ctx, &customerv1.SearchOverdueCustomersRequest{
		TenantId:   req.GetTenantId(),
		Pagination: toCustomerPagination(req.GetPagination()),
		Key:        req.GetKey(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	results := make([]*pb.OverdueCustomer, len(searchRsp.GetCustomers()))
	for k, v := range searchRsp.GetCustomers() {
		results[k] = &pb.OverdueCustomer{
			Customer:             toAppCustomer(v.GetCustomer()),
			LastMeasurementTime:  v.GetLastMeasurementTime(),
			SinceLastMeasurement: v.GetSinceLastMeasurement(),
			StaffName:            v.GetStaffName(),
		}
	}

	rsp.Customers = results
	rsp.TotalCount = searchRsp.GetTotalCount()

	return nil
}

// 验证request
func validateSearchOverdueCustomersRequest(req *pb.SearchOverdueCustomersRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
