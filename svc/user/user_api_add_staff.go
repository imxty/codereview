package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

const (
	// ErrStaffExceedLimit
	ErrStaffExceedLimit = 5003
	// ErrStaffPhoneHasBeenUsed
	ErrStaffPhoneHasBeenUsed = 5005
)

// 添加员工
func (u *UserAPIHandler) AddStaff(ctx context.Context, req *pb.AddStaffRequest, rsp *pb.AddStaffResponse) error {
	// 验证 request
	err := validateAddStaffRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Errorf(ErrTenantNotFound, "tenant_entity[%s] not found", req.GetTenantId())
	}

	// 查询用户
	user, err := u.userStore.GetActivateStaffByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
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

	// hash 密码
	hashedPassword := generateHashSHA256(req.GetPlainPassword())
	// 员工存在且已经加入租户
	if user != nil {
		return errors.Errorf(ErrStaffPhoneHasBeenUsed, "staff phone has been used [%s]", req.GetPhone())
	} else {
		// 查看当前租户是否有过该用户注册过
		// 查询用户
		uu, err := u.userStore.GetTenantStaffByPhone(ctx, req.GetTenantId(), req.GetPhone())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		if uu != nil {
			if uu.GetIsActivated() {
				return errors.Errorf(ErrStaffPhoneHasBeenUsed, "staff phone has been used [%s]", req.GetPhone())
			}
		}
	}
	// 初始化员工
	staff := &domain.User{
		UserID:         xid.New().String(),
		TenantID:       req.GetTenantId(),
		Username:       req.GetName(),
		RoleType:       domain.RoleTypeStaff,
		Nickname:       req.GetName(),
		Phone:          req.GetPhone(),
		HashedPassword: hashedPassword,
		IsActivated:    true,
	}

	// 创建员工
	err = u.userStore.CreateStaff(ctx, staff)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证 request
func validateAddStaffRequest(req *pb.AddStaffRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetName() == "" {
		return gerr.New("name should not be empty")
	}
	if req.GetPlainPassword() == "" {
		return gerr.New("plain_password should not be empty")
	}
	if len(req.GetPlainPassword()) < 6 || len(req.GetPlainPassword()) > 20 {
		return gerr.New("password length must be 6-20")
	}
	return nil
}
