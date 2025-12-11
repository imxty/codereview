package user

import (
	"context"
	gerr "errors"

	notificationv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 重置密码
func (u *UserAPIHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest, rsp *pb.ResetPasswordResponse) error {
	err := validateResetPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询系统用户
	user, err := u.userStore.GetSystemUserByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if user == nil {
		return errors.Errorf(ErrSystemUserNotFound, "system user by phone[%s] not found", req.GetPhone())
	}

	// 验证手机号
	_, err = u.notificationAPI.CheckTxIdUsage(ctx, &notificationv1.CheckTxIdUsageRequest{
		TxId: req.GetTxId(),
	})
	if err != nil {
		return err
	}

	// 重置密码
	err = u.userStore.ResetPassword(ctx, user.GetUserID(), generateHashSHA256(req.GetNewPlainPassword()))
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证 request
func validateResetPasswordRequest(req *pb.ResetPasswordRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetNewPlainPassword() == "" {
		return gerr.New("new_plain_password should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms_code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx_id should not be empty")
	}
	return nil
}
