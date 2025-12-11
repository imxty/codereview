package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// UpdateOrganizationPassword 修改组织密码
func (u *UserAPIHandler) UpdateOrganizationPassword(ctx context.Context, req *pb.UpdateOrganizationPasswordRequest, rsp *pb.UpdateOrganizationPasswordResponse) error {

	// 1.验证request
	err := validateUpdateOrganizationPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询租户
	organization, err := u.userStore.GetOrganizationByOrganizationId(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if organization == nil {
		return errors.Errorf(ErrOrganizationNotExist, "organization[%s] not found", req.GetOrganizationId())
	}

	// 对比旧密码是否正确
	if generateHashSHA256(req.GetOldPlainPassword()) != organization.GetHashedPassword() {
		return errors.Error(ErrWrongOldPassword, "old password not right")
	}

	// 修改密码
	err = u.userStore.UpdateOrganizationUserPassword(ctx, req.GetOrganizationId(), generateHashSHA256(req.GetNewPlainPassword()), organization.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	return nil
}

// 验证request
func validateUpdateOrganizationPasswordRequest(req *pb.UpdateOrganizationPasswordRequest) error {
	if req.GetOldPlainPassword() == "" {
		return gerr.New("old_plain_password should not be empty")
	}
	if req.GetNewPlainPassword() == "" {
		return gerr.New("new_plain_password should not be empty")
	}
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	return nil
}
