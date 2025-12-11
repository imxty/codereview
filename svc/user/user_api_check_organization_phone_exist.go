package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 判断组织手机号是否存在
func (u *UserAPIHandler) CheckOrganizationPhoneExist(ctx context.Context, req *pb.CheckOrganizationPhoneExistRequest, rsp *pb.CheckOrganizationPhoneExistResponse) error {
	// 验证request
	err := validateCheckOrganizationPhoneExistRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查找组织
	organization, err := u.userStore.GetOrganizationByPhone(ctx, req.GetPhone())
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
func validateCheckOrganizationPhoneExistRequest(req *pb.CheckOrganizationPhoneExistRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	return nil
}
