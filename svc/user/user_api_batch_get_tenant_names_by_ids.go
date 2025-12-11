package user

import (
	"context"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) BatchGetTenantNamesByIDs(ctx context.Context, req *pb.BatchGetTenantNamesByIDsRequest, rsp *pb.BatchGetTenantNamesByIDsResponse) error {
	if req.GetTenantIds() == nil {
		rsp.TenantNames = make(map[string]string)
		return nil
	}

	// 没有值返回空切片
	if len(req.GetTenantIds()) == 0 {
		rsp.TenantNames = map[string]string{}
		return nil
	}

	// 批量获取商户信息
	tenants, err := u.userStore.BatchGetTenantNamesByIds(ctx, req.GetTenantIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回响应
	results := make(map[string]string)
	for _, each := range tenants {
		results[each.GetTenantID()] = each.GetStoreName()
	}
	rsp.TenantNames = results

	return nil
}
