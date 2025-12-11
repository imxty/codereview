package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrWrongOldPassword
	ErrWrongOldPassword = 5012
)

func (u *UserAPIHandler) UpdateAdminPassword(ctx context.Context, req *pb.UpdateAdminPasswordRequest, rsp *pb.UpdateAdminPasswordResponse) error {
	// 验证 request
	err := validateUpdateAdminPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询租户
	tenant, err := u.userStore.GetTenantUser(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenant == nil {
		return errors.Error(ErrTenantNotFound, "tenant not found")
	}

	// 对比旧密码是否正确
	if generateHashSHA256(req.GetOldPlainPassword()) != tenant.GetHashedPassword() {
		return errors.Error(ErrWrongOldPassword, "old password not right")
	}

	// 修改密码
	err = u.userStore.UpdateTenantUserPassword(ctx, req.GetTenantId(), generateHashSHA256(req.GetNewPlainPassword()), tenant.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	return nil
}

// 验证 request
func validateUpdateAdminPasswordRequest(req *pb.UpdateAdminPasswordRequest) error {
	if req.GetOldPlainPassword() == "" {
		return gerr.New("old_plain_password should not be empty")
	}
	if req.GetNewPlainPassword() == "" {
		return gerr.New("new_plain_password should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	return nil
}
