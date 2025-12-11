package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	ChinaAreaCode = "+86"
)

func (s *AppAPIHandler) SendVerificationCode(ctx context.Context, req *pb.SendVerificationCodeRequest, rsp *pb.SendVerificationCodeResponse) error {
	err := validateSendVerificationCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 发送请求
	sendRsp, err := s.notificationAPI.SendVerificationCode(ctx, &notificationpb.SendVerificationCodeRequest{
		Phone:          req.GetPhone(),
		TemplateAction: toSvcTemplate(req.GetTemplateAction()),
		Language:       toSvcLanguage(req.GetLanguage()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	rsp.TxId = sendRsp.GetTxId()
	return nil
}

// 验证request
func validateSendVerificationCodeRequest(req *pb.SendVerificationCodeRequest) error {
	// 目前只有简体中文
	if req.GetLanguage() != pb.Language_LANGUAGE_SIMPLIFIED_CHINESE {
		return gerr.New("language must be simplified_chinese")
	}
	if req.GetPhone() == "" {
		return gerr.New("invalid phone")
	}
	if req.GetTemplateAction() == pb.TemplateAction_TEMPLATE_ACTION_UNSET || req.GetTemplateAction() == pb.TemplateAction_TEMPLATE_ACTION_INVALID {
		return gerr.New("invalid template_action")
	}
	return nil
}
