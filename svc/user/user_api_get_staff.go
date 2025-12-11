package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 通过员工 ID 查询员工信息
func (u *UserAPIHandler) GetStaff(ctx context.Context, req *pb.GetStaffRequest, rsp *pb.GetStaffResponse) error {
	err := validateGetStaffRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询 staff
	staff, err := u.userStore.GetStaffByID(ctx, req.GetStaffId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if staff == nil {
		return errors.Error(ErrStaffNotExist, "staff not found")
	}

	// 员工被删除
	if !staff.GetIsActivated() {
		return errors.Errorf(ErrStaffHasNotActivated, "staff[%s] not activated", req.GetStaffId())
	}

	// 获取租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Error(ErrTenantNotFound, "tenant_entity not found")
	}

	sf := &pb.Staff{
		// 员工 ID
		StaffId: req.GetStaffId(),
		// 员工姓名
		Name: staff.GetNickname(),
		// 员工手机号
		Phone: staff.GetPhone(),
		// 员工状态
		IsActivated: staff.GetIsActivated(),
		IsDeleted:   !staff.GetIsActivated(),
		// 商户信息
		Tenant: toProtoTenant(tenantEntity),
	}
	rsp.Staff = sf
	return nil
}

// 验证 request
func validateGetStaffRequest(req *pb.GetStaffRequest) error {
	if req.GetStaffId() == "" {
		return gerr.New("staff id should not be empty")
	}
	return nil
}
