package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrStaffHasNotActivated
	ErrStaffHasNotActivated = 5010
	// ErrPhoneHasBeenUsed
	ErrPhoneHasBeenUsed = 5011
	// ErrStaffNotExist
	ErrStaffNotExist = 5008
)

// 更新员工信息
func (u *UserAPIHandler) UpdateStaff(ctx context.Context, req *pb.UpdateStaffRequest, rsp *pb.UpdateStaffResponse) error {
	// 验证 request
	err := validateUpdateStaffRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取员工信息
	staff, err := u.userStore.GetStaffByID(ctx, req.GetStaffId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if staff == nil {
		return errors.Error(ErrStaffNotExist, "staff not found")
	}

	// 如果手机号码不一样，则判断新手机号是否注册
	if staff.GetPhone() != req.GetPhone() {
		// 查询用户
		user, err := u.userStore.GetStaffByPhone(ctx, "+86", req.GetPhone())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 员工存在
		if user != nil {
			return errors.Error(ErrPhoneHasBeenUsed, "phone has been used")
		}
	}

	// hash 密码
	hashedPassword := staff.GetHashedPassword()
	if req.GetPlainPassword() != "" {
		hashedPassword = generateHashSHA256(req.GetPlainPassword())
	}
	err = u.userStore.UpdateStaff(ctx, req.GetStaffId(), req.GetPhone(), req.GetName(), hashedPassword, staff.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证 request
func validateUpdateStaffRequest(req *pb.UpdateStaffRequest) error {
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetName() == "" {
		return gerr.New("staff_name should not be empty")
	}
	return nil
}
