package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 判断组织用户名是否存在
func (u *UserAPIHandler) CheckOrganizationUsernameExist(ctx context.Context, req *pb.CheckOrganizationUsernameExistRequest, rsp *pb.CheckOrganizationUsernameExistResponse) error {
	// 验证request
	err := validateCheckOrganizationUsernameExistRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 通过用户名查找组织
	organization, err := u.userStore.GetOrganizationByUsername(ctx, req.GetUsername())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if organization != nil {
		rsp.Exist = true
		return nil
	}

	rsp.Exist = false

	return nil
}

// 验证request
func validateCheckOrganizationUsernameExistRequest(req *pb.CheckOrganizationUsernameExistRequest) error {
	if req.GetUsername() == "" {
		return gerr.New("username should not be empty")
	}
	return nil
}
