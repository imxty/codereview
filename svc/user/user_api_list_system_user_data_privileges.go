package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取系统用户的数据权限
func (u *UserAPIHandler) ListSystemUserDataPrivileges(ctx context.Context, req *pb.ListSystemUserDataPrivilegesRequest, rsp *pb.ListSystemUserDataPrivilegesResponse) error {
	err := validateListSystemUserDataPrivilegesRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取用户信息
	user, err := u.userStore.GetStaffByID(ctx, req.GetUserId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if user == nil {
		return errors.Errorf(ErrSystemUserNotFound, "user[%s] not found", req.GetUserId())
	}

	// 获取用户数据权限
	organization_ids, err := u.userStore.GetDataPrivilegeByID(ctx, user.GetPrivilegeID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 全部数据
	if len(organization_ids) == 1 && organization_ids[0] == "all" {
		os, err := u.userStore.ListAllOrganizations(ctx)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		organization_ids = make([]string, len(os))
		for k, v := range os {
			organization_ids[k] = v.GetOrganizationID()
		}
	}

	// 返回响应
	rsp.OrganizationIds = organization_ids

	return nil
}

// 验证 request
func validateListSystemUserDataPrivilegesRequest(req *pb.ListSystemUserDataPrivilegesRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	return nil
}
