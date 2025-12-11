package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) UpdatePrivilegeGroup(ctx context.Context, req *pb.UpdatePrivilegeGroupRequest, rsp *pb.UpdatePrivilegeGroupResponse) error {
	err := validateUpdatePrivilegeGroupRequest(req)
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
		return errors.Errorf(ErrSystemUserNotFound, "system user [%s] not found", req.GetUserId())
	}

	// 获取权限组信息
	group, err := u.userStore.GetPrivilegeGroupById(ctx, req.GetPrivilegeId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 权限组不存在
	if group == nil {
		return errors.Errorf(ErrPrivilegeGroupNotExist, "privilege group [%s] not found", req.GetPrivilegeId())
	}

	// 开启事务
	ctx = u.userStore.BeginTx(ctx)

	// 修改权限组信息
	err = u.userStore.UpdatePrivilegeGroup(ctx, req.GetPrivilegeId(), req.GetPrivilegeName(), req.GetRemark())
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 删除页面权限
	err = u.userStore.DeletePrivilegePages(ctx, req.GetPrivilegeId())
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 删除目录
	err = u.userStore.DeletePrivilegeMenus(ctx, req.GetPrivilegeId())
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 删除数据权限
	datas, err := u.userStore.BatchGetDataPrivileges(ctx, req.GetPrivilegeId())
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	ids := make([]string, len(datas))
	for k, v := range datas {
		ids[k] = v.GetDataPrivilegeID()
	}
	err = u.userStore.DeleteDataPrivileges(ctx, ids)
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 创建新的页面权限
	pbPages := req.GetPagePrivileges()
	if len(pbPages) > 0 {
		pageIds := make([]string, len(pbPages))
		for i, page := range pbPages {
			pageIds[i] = toPagePrivilege(page)
		}
		err = u.userStore.CreatePagePrivileges(ctx, req.GetPrivilegeId(), pageIds)
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

		err = u.userStore.CreateMenus(ctx, req.GetPrivilegeId(), menuIds)
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	// 数据权限
	pbOrganizations := req.GetDataPrivileges()
	if len(pbOrganizations) > 0 {
		err = u.userStore.CreateDataPrivileges(ctx, req.GetPrivilegeId(), pbOrganizations)
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	// 提交事务
	u.userStore.CommitTx(ctx)
	return nil
}

func validateUpdatePrivilegeGroupRequest(req *pb.UpdatePrivilegeGroupRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetPrivilegeName() == "" {
		return gerr.New("privilege_name should not be empty")
	}
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
