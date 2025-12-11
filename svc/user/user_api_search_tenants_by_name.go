package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 通过商户名称查询商户 ID
func (u *UserAPIHandler) SearchTenantsByName(ctx context.Context, req *pb.SearchTenantsByNameRequest, rsp *pb.SearchTenantsByNameResponse) error {
	err := validateSearchTenantsByNameRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 通过名称查询 ids
	tenantIds, err := u.userStore.SearchTenantIdsByName(ctx, req.GetTenantName())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.TenantIds = tenantIds

	return nil
}

// 验证 request
func validateSearchTenantsByNameRequest(req *pb.SearchTenantsByNameRequest) error {
	if req.GetTenantName() == "" {
		return gerr.New("tenant_name should not be empty")
	}
	return nil
}
