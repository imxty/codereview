package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) SendPhoneVerificationCode(ctx context.Context, req *pb.SendPhoneVerificationCodeRequest, rsp *pb.SendPhoneVerificationCodeResponse) error {
	err := validateSendPhoneVerificationCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	templateAction, err := toSvcTemplateAction(req.GetTemplateAction())
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 发送短信
	sendRsp, err := s.notificationAPI.SendVerificationCode(ctx, &notificationpb.SendVerificationCodeRequest{
		Phone:          req.GetPhone(),
		Language:       toSvcLanguage(),
		TemplateAction: templateAction,
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.TxId = sendRsp.GetTxId()
	return nil
}

// 验证request
func validateSendPhoneVerificationCodeRequest(req *pb.SendPhoneVerificationCodeRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetLanguage() == pb.Language_LANGUAGE_UNSET || req.GetLanguage() == pb.Language_LANGUAGE_INVALID {
		return gerr.New("linvalid language")
	}
	if req.GetTemplateAction() == pb.TemplateAction_TEMPLATE_ACTION_INVALID || req.GetTemplateAction() == pb.TemplateAction_TEMPLATE_ACTION_UNSET {
		return gerr.New("invalid template_action")
	}
	return nil
}
