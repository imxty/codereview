package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 删除员工
func (u *UserAPIHandler) DeleteStaff(ctx context.Context, req *pb.DeleteStaffRequest, rsp *pb.DeleteStaffResponse) error {
	// 验证 request
	err := validateDeleteStaffRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询员工
	staff, err := u.userStore.GetStaffByID(ctx, req.GetStaffId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if staff == nil {
		return errors.Errorf(ErrStaffNotExist, "staff[%s] not found", req.GetStaffId())
	}

	// 删除员工
	err = u.userStore.DeleteStaff(ctx, req.GetStaffId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证 request
func validateDeleteStaffRequest(req *pb.DeleteStaffRequest) error {
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	return nil
}
