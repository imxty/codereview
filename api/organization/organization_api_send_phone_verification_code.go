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

func (s *OrganizationAPIHandler) SendPhoneVerificationCode(ctx context.Context, req *pb.SendPhoneVerificationCodeRequest, rsp *pb.SendPhoneVerificationCodeResponse) error {
	err := validateSendPhoneVerificationCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	action, _ := toSvcTemplateAction(req.GetTemplateAction())
	// 发送短信
	sendRsp, err := s.notificationAPI.SendVerificationCode(ctx, &notificationv1.SendVerificationCodeRequest{
		Phone:          req.GetPhone(),
		Language:       notificationv1.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		TemplateAction: action,
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
	if req.GetTemplateAction() == pb.TemplateAction_TEMPLATE_ACTION_INVALID || req.GetTemplateAction() == pb.TemplateAction_TEMPLATE_ACTION_UNSET {
		return gerr.New("invalid template action")
	}
	return nil
}
