package user

import (
	"context"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取权限组列表
func (u *UserAPIHandler) ListPrivilegeGroups(ctx context.Context, req *pb.ListPrivilegeGroupsRequest, rsp *pb.ListPrivilegeGroupsResponse) error {
	// 获取权限组列表
	groups, err := u.userStore.ListPrivilegeGroups(ctx)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回响应
	results := make([]*pb.PrivilegeGroup, len(groups))
	for i, each := range groups {
		results[i] = toProtoPrivilegeGroup(each)
	}
	rsp.PrivilegeGroups = results

	return nil
}
