package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 查询商户名是否存在
func (u *UserAPIHandler) CheckTenantNameExists(ctx context.Context, req *pb.CheckTenantNameExistsRequest, rsp *pb.CheckTenantNameExistsResponse) error {
	err := validateCheckTenantNameExists(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	tenant, err := u.userStore.GetTenantByName(ctx, req.GetTenantName())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.IsExists = tenant != nil

	return nil
}

// 验证 request
func validateCheckTenantNameExists(req *pb.CheckTenantNameExistsRequest) error {
	if req.GetTenantName() == "" {
		return gerr.New("tenant_name should not be empty")
	}
	return nil
}
