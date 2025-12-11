package customer

import (
	"context"
	gerr "errors"

	"github.com/jinmukeji/go-pkg/v2/age"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *CustomerAPIHandler) SearchCustomers(ctx context.Context, req *pb.SearchCustomersRequest, rsp *pb.SearchCustomersResponse) error {
	// 验证request
	err := validateSearchCustomersRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 从数据库获取常客列表
	dcustomers, err := s.customerStore.SearchCustomers(ctx, req.GetTenantId(), req.GetKeywords(), int(req.GetPagination().GetOffset()), int(req.GetPagination().GetSize()))
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 转换
	appCustomers := make([]*pb.Customer, len(dcustomers))
	for k, v := range dcustomers {
		appCustomers[k] = toProtoCustomer(v, getNicknameInitial(v.GetNickname()), int32(age.Age((v.GetBirthday()))))
	}

	// 返回结果
	rsp.Customers = appCustomers
	return nil
}

// 验证request
func validateSearchCustomersRequest(req *pb.SearchCustomersRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}
