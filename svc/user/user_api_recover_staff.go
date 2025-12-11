package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrStaffHasBeenAdded
	ErrStaffHasBeenAdded = 5019
	// ErrStaffHasActivated
	ErrStaffHasActivated = 5020
)

func (u *UserAPIHandler) RecoverStaff(ctx context.Context, req *pb.RecoverStaffRequest, rsp *pb.RecoverStaffResponse) error {
	err := validateRecoverStaff(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Error(ErrTenantNotFound, "tenant_entity not found")
	}

	// 获取所有在职员工
	staffs, err := u.userStore.ListActivatedStaffs(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 查看人数是否超过最大值
	if int32(len(staffs)) >= tenantEntity.GetStaffCountQuota() {
		return errors.Error(ErrStaffExceedLimit, "staff number exceed limit")
	}

	// 获取员工信息
	staff, err := u.userStore.GetStaffByID(ctx, req.GetStaffId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if staff == nil {
		return errors.Error(ErrStaffNotExist, "staff not found")
	}

	// 如果当前员工不是该租户，或者该员工已经激活则无法恢复
	if staff.GetIsActivated() {
		return errors.Error(ErrStaffHasActivated, "staff has been activated")
	}
	if staff.GetTenantID() != req.GetTenantId() {
		return errors.Error(ErrStaffHasBeenAdded, "staff has joined other company")
	}

	// 查询当前手机号是否已经被注册，如果被注册则无法恢复
	user, err := u.userStore.GetStaffByPhone(ctx, "+86", staff.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if user != nil {
		return errors.Error(ErrPhoneHasBeenUsed, "phone has been used")
	}

	// 恢复员工
	err = u.userStore.RecoverStaff(ctx, req.GetStaffId(), req.GetTenantId(), staff.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

func validateRecoverStaff(req *pb.RecoverStaffRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	return nil
}
