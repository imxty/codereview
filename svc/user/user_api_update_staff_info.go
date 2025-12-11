package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) UpdateStaffInfo(ctx context.Context, req *pb.UpdateStaffInfoRequest, rsp *pb.UpdateStaffInfoResponse) error {
	err := validateUpdateStaffInfo(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	staff, err := u.userStore.GetSystemUserByID(ctx, req.GetStaffId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if staff == nil {
		return errors.Errorf(ErrSystemUserNotFound, "system user[%s] not found", req.GetStaffId())
	}

	err = u.userStore.UpdateStaffInfo(ctx, req.GetStaffId(), req.GetStaffName(), req.GetStaffPhone(), req.GetPrivilegeId(), req.GetRemark())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

func validateUpdateStaffInfo(req *pb.UpdateStaffInfoRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetStaffName() == "" {
		return gerr.New("staff_name should not be empty")
	}
	if req.GetStaffPhone() == "" {
		return gerr.New("staff_phone should not be empty")
	}
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
