package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取权限组信息
func (s *BossAPIHandler) GetPrivilegeGroup(ctx context.Context, req *pb.GetPrivilegeGroupRequest, rsp *pb.GetPrivilegeGroupResponse) error {
	err := validateGetPrivilegeGroupRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询权限组信息
	groupRsp, err := s.userAPI.GetPrivilegeGroup(ctx, &userv1.GetPrivilegeGroupRequest{
		PrivilegeId: req.GetPrivilegeId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 查询组织信息
	organizationRsp, err := s.userAPI.BatchGetOrganizationNamesByIDs(ctx, &userv1.BatchGetOrganizationNamesByIDsRequest{
		OrganizationId: groupRsp.GetPrivilegeGroup().GetDataPrivileges(),
	})
	if err != nil {
		return err
	}

	// 返回数据
	rsp.PrivilegeGroup = toApiPrivilegeGroup(groupRsp.GetPrivilegeGroup())
	rsp.PrivilegeGroup.DataPrivileges = organizationRsp.GetOrganizations()

	return nil
}

// 验证request
func validateGetPrivilegeGroupRequest(req *pb.GetPrivilegeGroupRequest) error {
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
