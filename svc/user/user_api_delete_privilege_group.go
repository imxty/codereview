package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	ErrPrivilegeGroupUsing = 5347
)

func (u *UserAPIHandler) DeletePrivilegeGroup(ctx context.Context, req *pb.DeletePrivilegeGroupRequest, rsp *pb.DeletePrivilegeGroupResponse) error {
	err := validateDeletePrivilegeGroupRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取操作人员信息
	user, err := u.userStore.GetSystemUserByID(ctx, req.GetUserId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 操作人员不存在
	if user == nil {
		return errors.Errorf(ErrSystemUserNotFound, "system user[%s] not found", req.GetUserId())
	}

	// 获取权限组信息
	group, err := u.userStore.GetPrivilegeGroupById(ctx, req.GetPrivilegeId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 找不到权限组
	if group == nil {
		return errors.Errorf(ErrPrivilegeGroupNotExist, "privilege group[%s] not found", req.GetPrivilegeId())
	}

	// 检查权限组是否和用户绑定
	bindUser, err := u.userStore.GetUserByPrivilegeID(ctx, req.GetPrivilegeId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 如果权限组正在使用，不能删除
	if bindUser != nil {
		return errors.Errorf(ErrPrivilegeGroupUsing, "privilege group[%s] using", req.GetPrivilegeId())
	}

	// 删除权限组
	err = u.userStore.DeletePrivilegeGroup(ctx, req.GetPrivilegeId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

func validateDeletePrivilegeGroupRequest(req *pb.DeletePrivilegeGroupRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
