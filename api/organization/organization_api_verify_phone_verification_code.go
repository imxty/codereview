package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	notificationv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *OrganizationAPIHandler) VerifyPhoneVerificationCode(ctx context.Context, req *pb.VerifyPhoneVerificationCodeRequest, rsp *pb.VerifyPhoneVerificationCodeResponse) error {

	err := validateVerifyPhoneVerificationCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	action, _ := toSvcTemplateAction(req.GetTemplateAction())
	_, err = s.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationv1.VerifyPhoneVerificationCodeRequest{
		Phone:   req.GetPhone(),
		SmsCode: req.GetSmsCode(),
		TxId:    req.GetTxId(),
		Action:  action,
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateVerifyPhoneVerificationCodeRequest(req *pb.VerifyPhoneVerificationCodeRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx id should not be empty")
	}
	if req.GetTemplateAction() == pb.TemplateAction_TEMPLATE_ACTION_INVALID || req.GetTemplateAction() == pb.TemplateAction_TEMPLATE_ACTION_UNSET {
		return gerr.New("invalid template action")
	}
	return nil
}
