package customer

import (
	"context"
	gerr "errors"

	"github.com/jinmukeji/go-pkg/v2/age"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 通过 ID 批量获取常客信息
func (c *CustomerAPIHandler) BatchGetCustomersByIds(ctx context.Context, req *pb.BatchGetCustomersByIdsRequest, rsp *pb.BatchGetCustomersByIdsResponse) error {
	err := validateBatchGetCustomersByIds(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	cs, err := c.customerStore.BatchGetCustomers(ctx, req.GetCustomerIds(), req.GetIncludeDeleted())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	results := make(map[string]*pb.Customer)
	for _, each := range cs {
		results[each.GetCustomerID()] = toProtoCustomer(each, getNicknameInitial(each.GetNickname()), int32(age.Age(each.GetBirthday())))
	}

	rsp.Customers = results

	return nil
}

// 验证 request
func validateBatchGetCustomersByIds(req *pb.BatchGetCustomersByIdsRequest) error {
	if req.GetCustomerIds() == nil {
		return gerr.New("customer_ids should not be empty")
	}
	return nil
}
