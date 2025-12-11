package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取权限组组织名
func (u *UserAPIHandler) ListPrivilegeOrganizationNames(ctx context.Context, req *pb.ListPrivilegeOrganizationNamesRequest, rsp *pb.ListPrivilegeOrganizationNamesResponse) error {
	err := validateListPrivilegeOrganizationNames(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	privilege_id := req.GetPrivilegeId()
	// 获取数据权限组织 ID
	organizations, err := u.userStore.GetDataPrivilegeByID(ctx, privilege_id)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	var users []domain.UserIntf
	// 判断是否有全部权限
	if len(organizations) == 1 && organizations[0] == "all" {
		// 获取全部组织信息
		users, err = u.userStore.ListAllOrganizations(ctx)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	} else {
		// 根据组织 ID 获取组织信息
		users, err = u.userStore.ListOrganizations(ctx, organizations)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	results := make(map[string]string)
	for _, each := range users {
		results[each.GetUserID()] = each.GetNickname()
	}

	rsp.Organizations = results

	return nil
}

// 验证 request
func validateListPrivilegeOrganizationNames(req *pb.ListPrivilegeOrganizationNamesRequest) error {
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
