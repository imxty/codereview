package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 检查系统用户数据权限
func (u *UserAPIHandler) CheckSystemUserDataPrivilege(ctx context.Context, req *pb.CheckSystemUserDataPrivilegeRequest, rsp *pb.CheckSystemUserDataPrivilegeResponse) error {
	err := validateCheckSystemUserDataPrivilegeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取系统用户信息
	user, err := u.userStore.GetSystemUserByID(ctx, req.GetUserId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if user == nil {
		rsp.Allow = false
		return nil
	}

	// 获取用户数据权限
	data, err := u.userStore.GetDataPrivilegeByOrganizationID(ctx, user.GetPrivilegeID(), req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 如果数据不存在，返回 false，不能访问
	if data == nil {
		rsp.Allow = false
		return nil
	}

	rsp.Allow = true

	return nil
}

// 验证 request
func validateCheckSystemUserDataPrivilegeRequest(req *pb.CheckSystemUserDataPrivilegeRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
