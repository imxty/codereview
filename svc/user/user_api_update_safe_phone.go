package user

import (
	"context"
	gerr "errors"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) UpdateSafePhone(ctx context.Context, req *pb.UpdateSafePhoneRequest, rsp *pb.UpdateSafePhoneResponse) error {
	// 1.验证request
	err := validateUpdateSafePhoneRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询租户
	tenant, err := u.userStore.GetTenantUser(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenant == nil {
		return errors.Error(ErrTenantNotFound, "tenant not found")
	}

	// 查询新手机是否已经注册
	nu, err := u.userStore.GetTenantUserByPhone(ctx, req.GetNewSafePhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if nu != nil {
		return errors.Errorf(ErrPhoneHasBeenUsed, "tenant phone [%s] has been used", req.GetNewSafePhone())
	}

	// 验证短信
	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   tenant.GetPhone(),
		SmsCode: req.GetOldPhoneSmsCode(),
		TxId:    req.GetOldPhoneSmsTxId(),
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE,
	})
	if err != nil {
		return err
	}

	// 获取新手机号最新的验证码
	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   req.GetNewSafePhone(),
		SmsCode: req.GetNewPhoneSmsCode(),
		TxId:    req.GetNewPhoneSmsTxId(),
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE,
	})
	if err != nil {
		return err
	}

	// 查询租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Error(ErrTenantNotFound, "tenant not found")
	}

	// 修改新手机号
	err = u.userStore.UpdateTenantUserSafePhone(ctx, req.GetTenantId(), req.GetNewSafePhone(), tenant.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	err = u.userStore.UpdateTenantSafePhone(ctx, req.GetTenantId(), req.GetNewSafePhone(), tenantEntity.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

func validateUpdateSafePhoneRequest(req *pb.UpdateSafePhoneRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
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
