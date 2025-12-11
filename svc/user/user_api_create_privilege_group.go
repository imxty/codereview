package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

// 创建权限组
func (u *UserAPIHandler) CreatePrivilegeGroup(ctx context.Context, req *pb.CreatePrivilegeGroupRequest, rsp *pb.CreatePrivilegeGroupResponse) error {
	err := validateCreatePrivilegeGroupRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 开启事务
	// 1. 创建权限组
	// 2. 创建页面权限
	// 3. 创建数据权限
	ctx = u.userStore.BeginTx(ctx)
	privilege_id := xid.New().String()

	// 权限组
	err = u.userStore.CreatePrivilegeGroup(ctx, privilege_id, req.GetPrivilegeName(), req.GetRemark())
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 页面权限
	pbPages := req.GetPagePrivileges()
	if len(pbPages) > 0 {
		pageIds := make([]string, len(pbPages))
		for i, page := range pbPages {
			pageIds[i] = toPagePrivilege(page)
		}
		err = u.userStore.CreatePagePrivileges(ctx, privilege_id, pageIds)
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}

		// 设置目录
		mPaths := make([]string, len(pageIds))
		for k, v := range pageIds {
			mPaths[k] = toMenuPath(v)
		}

		menus, err := u.userStore.GetMenusByPaths(ctx, mPaths)
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}

		menuIds := make([]string, len(menus))
		for k, v := range menus {
			menuIds[k] = v.GetMenuID().String()
		}

		err = u.userStore.CreateMenus(ctx, privilege_id, menuIds)
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	// 数据权限
	pbOrganizations := req.GetDataPrivileges()
	if len(pbOrganizations) > 0 {
		err = u.userStore.CreateDataPrivileges(ctx, privilege_id, pbOrganizations)
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	u.userStore.CommitTx(ctx)

	return nil
}

// 验证 request
func validateCreatePrivilegeGroupRequest(req *pb.CreatePrivilegeGroupRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetPrivilegeName() == "" {
		return gerr.New("privilege_name should not be empty")
	}
	return nil
}
