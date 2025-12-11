package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 验证手机号是否存在
func (u *UserAPIHandler) CheckSystemUserPhoneExist(ctx context.Context, req *pb.CheckSystemUserPhoneExistRequest, rsp *pb.CheckSystemUserPhoneExistResponse) error {
	err := validateCheckSystemUserPhoneExistRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 通过手机号查询用户
	user, err := u.userStore.GetSystemUserByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.Exist = user != nil

	return nil
}

// 验证 request
func validateCheckSystemUserPhoneExistRequest(req *pb.CheckSystemUserPhoneExistRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	return nil
}
