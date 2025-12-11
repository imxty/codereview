package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 删除系统用户（离职）
func (u *UserAPIHandler) DeleteSystemUser(ctx context.Context, req *pb.DeleteSystemUserRequest, rsp *pb.DeleteSystemUserResponse) error {
	err := validateDeleteSystemUser(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 检查用户
	user, err := u.userStore.GetSystemUserByID(ctx, req.GetUserId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if user == nil {
		return errors.Errorf(ErrSystemUserNotFound, "system user[%s] not found", req.GetUserId())
	}

	// 删除
	err = u.userStore.DeleteSystemUser(ctx, user.GetUserID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	return nil
}

// 验证 request
func validateDeleteSystemUser(req *pb.DeleteSystemUserRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	return nil
}
