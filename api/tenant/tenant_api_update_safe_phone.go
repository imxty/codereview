package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) UpdateSafePhone(ctx context.Context, req *pb.UpdateSafePhoneRequest, rsp *pb.UpdateSafePhoneResponse) error {
	err := validateUpdateSafePhoneRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 更新安全手机号请求
	_, err = s.userAPI.UpdateSafePhone(ctx, &userv1.UpdateSafePhoneRequest{
		// 租户ID
		TenantId: req.GetTenantId(),
		// 旧安全手机验证码
		OldPhoneSmsCode: req.GetOldPhoneSmsCode(),
		// 旧手机短信凭证ID
		OldPhoneSmsTxId: req.GetOldPhoneSmsTxId(),
		// 新手机号
		NewSafePhone: req.GetNewSafePhone(),
		// 新安全手机验证码
		NewPhoneSmsCode: req.GetNewPhoneSmsCode(),
		// 新安全手机短信凭证ID
		NewPhoneSmsTxId: req.GetNewPhoneSmsTxId(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
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
