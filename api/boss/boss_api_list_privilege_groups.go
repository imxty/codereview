package boss

import (
	"context"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
)

// 获取权限组列表
func (s *BossAPIHandler) ListPrivilegeGroups(ctx context.Context, req *pb.ListPrivilegeGroupsRequest, rsp *pb.ListPrivilegeGroupsResponse) error {
	groupsRsp, err := s.userAPI.ListPrivilegeGroups(ctx, &userv1.ListPrivilegeGroupsRequest{})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 转换权限列表
	apiGroups := make([]*pb.PrivilegeGroup, len(groupsRsp.GetPrivilegeGroups()))
	for k, v := range groupsRsp.GetPrivilegeGroups() {
		apiGroups[k] = toApiPrivilegeGroup(v)
	}

	// 返回响应
	rsp.PrivilegeGroups = apiGroups

	return nil
}
