package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	notificationv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 发送短信验证码
func (s *BossAPIHandler) SendPhoneVerificationCode(ctx context.Context, req *pb.SendPhoneVerificationCodeRequest, rsp *pb.SendPhoneVerificationCodeResponse) error {
	err := validateSendPhoneVerificationCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	sendRsp, err := s.notificationAPI.SendVerificationCode(ctx, &notificationv1.SendVerificationCodeRequest{
		Phone:          req.GetPhone(),
		TemplateAction: toSvcTemplateAction(req.GetTemplateAction()),
		Language:       toSvcLanguage(req.GetLanguage()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.TxId = sendRsp.GetTxId()

	return nil
}

// 验证 request
func validateSendPhoneVerificationCodeRequest(req *pb.SendPhoneVerificationCodeRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	return nil
}
