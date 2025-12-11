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
	// ErrPrivilegeGroupNotExist
	ErrPrivilegeGroupNotExist = 5401
)

// 创建系统用户
func (u *UserAPIHandler) CreateSystemUser(ctx context.Context, req *pb.CreateSystemUserRequest, rsp *pb.CreateSystemUserResponse) error {
	err := validateCreateSystemUserRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 检查权限组
	privilegeGroup, err := u.userStore.GetPrivilegeGroupById(ctx, req.GetPrivilegeId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if privilegeGroup == nil {
		return errors.Errorf(ErrPrivilegeGroupNotExist, "privilege group[%s] not exist", req.GetPrivilegeId())
	}

	user, err := u.userStore.GetSystemUserByPhone(ctx, req.GetStaffPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 如果手机号存在，不能添加
	if user != nil {
		return errors.Errorf(ErrPhoneHasBeenUsed, "system_user phone[%s] has been used", req.GetStaffPhone())
	}

	// 创建系统用户
	systemUser := &domain.User{
		UserID:      xid.New().String(),
		Phone:       req.GetStaffPhone(),
		PrivilegeID: req.GetPrivilegeId(),
		Remark:      req.GetRemark(),
		Nickname:    req.GetStaffName(),
		Username:    req.GetStaffName(),
		RoleType:    domain.RoleTypeJmAdmin,
		IsActivated: true,
	}
	err = u.userStore.CreateSystemUser(ctx, systemUser)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证 request
func validateCreateSystemUserRequest(req *pb.CreateSystemUserRequest) error {
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
