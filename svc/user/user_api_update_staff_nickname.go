package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 更新员工的昵称请求
func (u *UserAPIHandler) UpdateStaffNickname(ctx context.Context, req *pb.UpdateStaffNicknameRequest, rsp *pb.UpdateStaffNicknameResponse) error {
	err := validateUpdateStaffNicknameRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 查询用户
	staff, err := u.userStore.GetStaffByID(ctx, req.GetStaffId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if staff == nil {
		return errors.Error(ErrStaffNotExist, "staff not exist")
	}
	// 如果没有激活则无法更新
	if !staff.GetIsActivated() {
		return errors.Errorf(ErrStaffHasNotActivated, "staff[%s] is not activated", req.GetStaffId())
	}

	// 更新员工的昵称
	err = u.userStore.UpdateStaffNickname(ctx, req.GetStaffId(), req.GetNewNickname(), staff.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证 request
func validateUpdateStaffNicknameRequest(req *pb.UpdateStaffNicknameRequest) error {
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetNewNickname() == "" {
		return gerr.New("nickname should not be empty")
	}
	return nil
}
