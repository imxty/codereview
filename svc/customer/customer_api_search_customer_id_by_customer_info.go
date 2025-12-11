package customer

import (
	"context"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 根据客户名称 + 手机号 + 类型搜索客户 ID
func (c *CustomerAPIHandler) SearchCustomerIDByCustomerInfo(ctx context.Context, req *pb.SearchCustomerIDByCustomerInfoRequest, rsp *pb.SearchCustomerIDByCustomerInfoResponse) error {
	id, err := c.customerStore.SearchCustomerIDByPhoneAndName(ctx, req.GetPhone(), req.GetName())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.CustomerId = id

	return nil
}
