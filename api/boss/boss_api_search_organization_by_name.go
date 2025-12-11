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

// 根据组织名查询
func (s *BossAPIHandler) SearchOrganizationByName(ctx context.Context, req *pb.SearchOrganizationByNameRequest, rsp *pb.SearchOrganizationByNameResponse) error {
	err := validateSearchOrganizationByName(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchRsp, err := s.userAPI.SearchOrganizationByName(ctx, &userv1.SearchOrganizationByNameRequest{
		OrganizationName: req.GetOrganizationName(),
		PrivilegeId:      req.GetPrivilegeId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回响应
	rsp.Organizations = searchRsp.GetOrganizations()

	return nil
}

func validateSearchOrganizationByName(req *pb.SearchOrganizationByNameRequest) error {
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
