package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrSamePasswords
	ErrSamePasswords = 5013
)

// 修改员工密码
func (u *UserAPIHandler) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest, rsp *pb.UpdatePasswordResponse) error {
	err := validateUpdatePasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询用户
	staff, err := u.userStore.GetStaffByID(ctx, req.GetStaffId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if staff == nil {
		return errors.Errorf(ErrStaffNotExist, "staff[%s] not exist", req.GetStaffId())
	}

	if !staff.GetIsActivated() {
		return errors.Errorf(ErrStaffHasNotActivated, "staff[%s] is not activated", req.GetStaffId())
	}

	// 新旧密码不能一致
	if req.GetNewPlainPassword() == req.GetOldPlainPassword() {
		return errors.Errorf(ErrSamePasswords, "new password should not be same as the old one.[%s]", req.GetNewPlainPassword())
	}
	// 旧密码不正确
	if generateHashSHA256(req.GetOldPlainPassword()) != staff.GetHashedPassword() {
		return errors.Error(ErrWrongOldPassword, "wrong old password")
	}
	// 更新密码
	err = u.userStore.UpdateStaffPassword(ctx, req.GetStaffId(), generateHashSHA256(req.GetNewPlainPassword()), staff.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	return nil
}

// 验证 request
func validateUpdatePasswordRequest(req *pb.UpdatePasswordRequest) error {
	if req.GetNewPlainPassword() == "" {
		return gerr.New("new_password should not be empty")
	}
	if req.GetOldPlainPassword() == "" {
		return gerr.New("old_password should not be empty")
	}
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
