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

// 获取权限组组织名
func (s *BossAPIHandler) ListPrivilegeOrganizationNames(ctx context.Context, req *pb.ListPrivilegeOrganizationNamesRequest, rsp *pb.ListPrivilegeOrganizationNamesResponse) error {
	err := validateListPrivilegeOrganizationNames(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	listRsp, err := s.userAPI.ListPrivilegeOrganizationNames(ctx, &userv1.ListPrivilegeOrganizationNamesRequest{
		PrivilegeId: req.GetPrivilegeId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.Organizations = listRsp.GetOrganizations()

	return nil
}

// 验证 request
func validateListPrivilegeOrganizationNames(req *pb.ListPrivilegeOrganizationNamesRequest) error {
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
