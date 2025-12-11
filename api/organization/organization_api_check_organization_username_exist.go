package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 判断组织用户名是否存在
func (s *OrganizationAPIHandler) CheckOrganizationUsernameExist(ctx context.Context, req *pb.CheckOrganizationUsernameExistRequest, rsp *pb.CheckOrganizationUsernameExistResponse) error {
	err := validateCheckOrganizationUsernameExistRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 检查组织名是否存在
	checkRsp, err := s.userAPI.CheckOrganizationUsernameExist(ctx, &userv1.CheckOrganizationUsernameExistRequest{
		Username: req.GetUsername(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回响应
	rsp.Exist = checkRsp.GetExist()

	return nil
}

// 验证request
func validateCheckOrganizationUsernameExistRequest(req *pb.CheckOrganizationUsernameExistRequest) error {
	if req.GetUsername() == "" {
		return gerr.New("username should not be empty")
	}
	return nil
}
