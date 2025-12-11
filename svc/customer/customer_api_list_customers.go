package customer

import (
	"context"
	gerr "errors"

	"github.com/jinmukeji/go-pkg/v2/age"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *CustomerAPIHandler) ListCustomers(ctx context.Context, req *pb.ListCustomersRequest, rsp *pb.ListCustomersResponse) error {

	// 验证request
	err := validateListCustomersRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 从数据库获取常客列表
	dCustomers, err := s.customerStore.ListCustomers(ctx, req.GetTenantId(), int(req.GetPagination().GetOffset()), int(req.GetPagination().GetSize()))
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 获取常客总数
	count, err := s.customerStore.GetCustomerCount(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 转换
	appCustomers := make([]*pb.Customer, len(dCustomers))
	for k, v := range dCustomers {
		appCustomers[k] = toProtoCustomer(v, v.GetInitial(), int32(age.Age(v.GetBirthday())))
	}

	rsp.Customers = appCustomers
	rsp.TotalCount = int32(count)
	return nil
}

// 验证request
func validateListCustomersRequest(req *pb.ListCustomersRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	if req.GetPagination().GetSize() <= 0 {
		return gerr.New("invalid pagination size")
	}
	if req.GetPagination().GetOffset() < 0 {
		return gerr.New("invalid pagination offset")
	}
	return nil
}
