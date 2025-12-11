package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取权限组信息
func (u *UserAPIHandler) GetPrivilegeGroup(ctx context.Context, req *pb.GetPrivilegeGroupRequest, rsp *pb.GetPrivilegeGroupResponse) error {
	err := validateGetPrivilegeGroupRequest(req)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获取权限组信息
	group, err := u.userStore.GetPrivilegeGroupById(ctx, req.GetPrivilegeId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获取数据权限
	datas, err := u.userStore.GetDataPrivilegeByID(ctx, req.GetPrivilegeId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if len(datas) == 1 && datas[0] == "all" {
		os, err := u.userStore.ListAllOrganizations(ctx)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		datas = make([]string, len(os))
		for k, v := range os {
			datas[k] = v.GetOrganizationID()
		}
	}

	// 获取页面权限
	pages, err := u.userStore.GetPagePrivilegeByID(ctx, req.GetPrivilegeId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	upages := make([]pb.PagePrivilege, len(pages))
	for k, v := range pages {
		upages[k] = toPrivilegePage(v)
	}

	rsp.PrivilegeGroup = &pb.PrivilegeGroup{
		PrivilegeId:    group.GetPrivilegeID(),
		PrivilegeName:  group.GetPrivilegeName(),
		Remark:         group.GetRemark(),
		PagePrivileges: upages,
		DataPrivileges: datas,
	}

	return nil
}

// 验证 request
func validateGetPrivilegeGroupRequest(req *pb.GetPrivilegeGroupRequest) error {
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
