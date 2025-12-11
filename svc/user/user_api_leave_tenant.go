package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 员工从商户离职
func (u *UserAPIHandler) LeaveTenant(ctx context.Context, req *pb.LeaveTenantRequest, rsp *pb.LeaveTenantResponse) error {
	// 验证 request
	err := validateLeaveTenantRequest(req)
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
	// 离开租户
	err = u.userStore.LeaveTenant(ctx, req.GetStaffId(), req.GetTenantId(), staff.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	return nil
}

// 验证 request
func validateLeaveTenantRequest(req *pb.LeaveTenantRequest) error {
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
