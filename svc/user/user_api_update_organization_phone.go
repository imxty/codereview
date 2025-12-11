package user

import (
	"context"
	gerr "errors"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 修改组织手机号
func (u *UserAPIHandler) UpdateOrganizationPhone(ctx context.Context, req *pb.UpdateOrganizationPhoneRequest, rsp *pb.UpdateOrganizationPhoneResponse) error {

	// 1.验证request
	err := validateUpdateOrganizationPhoneRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询租户
	organization, err := u.userStore.GetOrganizationByOrganizationId(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if organization == nil {
		return errors.Errorf(ErrOrganizationNotExist, "organization[%s] not found", req.GetOrganizationId())
	}

	// 查询新手机是否已经注册
	nu, err := u.userStore.GetOrganizationByPhone(ctx, req.GetNewSafePhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if nu != nil {
		return errors.Errorf(ErrPhoneHasBeenUsed, "organization phone [%s] has been used", req.GetNewSafePhone())
	}

	// 验证短信
	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   organization.GetPhone(),
		SmsCode: req.GetOldPhoneSmsCode(),
		TxId:    req.GetOldPhoneSmsTxId(),
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE,
	})
	if err != nil {
		return err
	}

	// 获取新手机号最新的验证码
	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   req.GetNewSafePhone(),
		SmsCode: req.GetNewPhoneSmsCode(),
		TxId:    req.GetNewPhoneSmsTxId(),
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE,
	})
	if err != nil {
		return err
	}

	// 修改新手机号
	err = u.userStore.UpdateOrganizationPhone(ctx, req.GetOrganizationId(), req.GetNewSafePhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证request
func validateUpdateOrganizationPhoneRequest(req *pb.UpdateOrganizationPhoneRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetOldPhoneSmsCode() == "" {
		return gerr.New("old_phone_sms_code should not be empty")
	}
	// 旧手机短信凭证ID
	if req.GetOldPhoneSmsTxId() == "" {
		return gerr.New("old_phone_sms_tx_id should not be empty")
	}
	// 新手机号
	if req.GetNewSafePhone() == "" {
		return gerr.New("new_safe_phone should not be empty")
	}
	// 新安全手机验证码
	if req.GetNewPhoneSmsCode() == "" {
		return gerr.New("new_phone_sms_code should not be empty")
	}
	// 新安全手机短信凭证ID
	if req.GetNewPhoneSmsTxId() == "" {
		return gerr.New("new_phone_sms_tx_id should not be empty")
	}
	return nil
}
