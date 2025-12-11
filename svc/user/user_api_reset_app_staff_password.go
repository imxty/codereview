package user

import (
	"context"
	gerr "errors"

	notificationv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) ResetAppStaffPassword(ctx context.Context, req *pb.ResetAppStaffPasswordRequest, rsp *pb.ResetAppStaffPasswordResponse) error {
	err := validateResetAppStaffPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 验证验证码
	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationv1.VerifyPhoneVerificationCodeRequest{
		Phone:   req.GetPhone(),
		SmsCode: req.GetSmsCode(),
		TxId:    req.GetTxId(),
		Action:  notificationv1.TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD,
	})
	if err != nil {
		return err
	}

	// 获取员工信息
	staff, err := u.userStore.GetStaffByPhone(ctx, "+86", req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if staff == nil {
		return errors.Error(ErrStaffNotExist, "staff not found")
	}

	// 员工不属于任何租户
	if !staff.GetIsActivated() {
		return errors.Error(ErrStaffNotBelongTenant, "staff not belong to any tenants")
	}

	err = u.userStore.ResetAppStaffPassword(ctx, staff.GetUserID(), generateHashSHA256(req.GetNewPlainPassword()))
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

func validateResetAppStaffPasswordRequest(req *pb.ResetAppStaffPasswordRequest) error {
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
