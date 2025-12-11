package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 查询系统用户
func (u *UserAPIHandler) SearchSystemUser(ctx context.Context, req *pb.SearchSystemUserRequest, rsp *pb.SearchSystemUserResponse) error {
	err := validateSearchSystemUserRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 搜索权限组
	groups, err := u.userStore.GetPrivilegeGroupsByName(ctx, req.GetPrivilegeGroupName())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	group_infos := make(map[string]domain.PrivilegePolicyIntf)
	group_ids := make([]string, len(groups))
	for k, v := range groups {
		group_ids[k] = v.GetPrivilegeID()
		group_infos[v.GetPrivilegeID()] = v
	}

	// 搜索系统用户
	users, count, err := u.userStore.SearchSystemUser(ctx, req.GetNickname(), req.GetPhone(), group_ids, req.GetIsActivated(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	results := make([]*pb.SystemUser, len(users))
	// 获取权限组信息
	for i, each := range users {
		results[i] = toProtoSystemUser(each)
		results[i].PrivilegeGroup = toProtoPrivilegeGroup(group_infos[each.GetPrivilegeID()])
	}

	// 返回响应
	rsp.Users = results
	rsp.TotalCount = count

	return nil
}

// 验证 request
func validateSearchSystemUserRequest(req *pb.SearchSystemUserRequest) error {
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}
