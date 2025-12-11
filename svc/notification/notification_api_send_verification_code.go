package notification

import (
	"context"
	"encoding/json"
	gerr "errors"
	"fmt"
	"time"

	"github.com/jinmukeji/go-pkg/v2/crypto/rand"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/svc/notification/client"
	"github.com/jinmukeji/huimaibao-service/svc/notification/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

const (
	// AliyunPlatform 阿里云平台
	AliyunPlatform = "Aliyun"
)

const (
	// ErrSmsRepeatedSending
	ErrSmsRepeatedSending = 5041
	// ErrMessageLimitControl
	ErrMessageLimitControl = 5306
)

const (
	MessageLimitControl = "isv.BUSINESS_LIMIT_CONTROL"
)

func (s *NotificationAPIHandler) SendVerificationCode(ctx context.Context, req *pb.SendVerificationCodeRequest, rsp *pb.SendVerificationCodeResponse) error {
	// 最新时间
	now := time.Now().UTC()
	err := validateSendVerificationCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取当前用户最近发送的一条短信的时间
	sms, err := s.notificationStore.GetLatestSmsByAction(ctx, req.GetPhone(), req.GetTemplateAction().String())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 如果之前有发送过短信，就去判断是否在一分钟重复发送
	if sms != nil {
		if sms.GetCreatedAt().Add(time.Minute).After(now) {
			return errors.Error(ErrSmsRepeatedSending, "sms is sent repeatedly within one minute")
		}
	}

	txId, err := s.SendAliyunMessage(ctx, req)
	if err != nil {
		if err.Error() == MessageLimitControl {
			return errors.Errorf(ErrMessageLimitControl, "phone[%s] : %s", req.GetPhone(), err.Error())
		}
		return errors.Error(codes.InvalidOperation, err.Error())
	}

	rsp.TxId = txId
	return nil
}

// 验证 request
func validateSendVerificationCodeRequest(req *pb.SendVerificationCodeRequest) error {
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

// 发送阿里云消息
func (s *NotificationAPIHandler) SendAliyunMessage(ctx context.Context, req *pb.SendVerificationCodeRequest) (string, error) {
	txId := xid.New().String()
	// 生成 6 位验证码
	smsCode := generateSmsCode(6)
	templateParams := make(map[string]string)
	templateParams["code"] = smsCode
	// 如果创建失败则 code 为商户 id
	if req.GetTemplateAction() == pb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL {
		templateParams = map[string]string{}
	}
	// proto 内容转 domain
	templateAction := req.GetTemplateAction().Number()

	isSucceed, errSms := s.aliYunSms.SendSms(req.GetPhone(), s.signName, client.TemplateAction(templateAction), client.SimpleChinese, templateParams)
	// 序列化模版参数
	templateParam, errTemplateParam := json.Marshal(templateParams)
	if errTemplateParam != nil {
		return "", fmt.Errorf("failed to get marshal format %v: %s", templateParams, errTemplateParam.Error())
	}
	// 构建短信信息
	sms := &domain.Sms{
		SmsID:          xid.New().String(),
		Phone:          req.GetPhone(),
		TemplateAction: req.GetTemplateAction().String(),
		PlatformType:   AliyunPlatform,
		Language:       "SimpleChinese",
		TemplateParam:  string(templateParam),
		TxID:           txId,
	}
	// 添加短信状态
	if isSucceed {
		sms.SmsStatus = domain.SmsStatusSentSuccessful
	} else {
		sms.SmsStatus = domain.SmsStatusSendFailed
	}
	if errSms != nil {
		sms.SmsErrorLog = errSms.Error()
	}
	err := s.notificationStore.CreateSmsRecord(ctx, sms)
	if err != nil {
		return "", errors.Error(codes.DataAccessFailed, err.Error())
	}
	return txId, errSms
}

// 生成短信验证码
func generateSmsCode(width int) string {
	smsCode, _ := rand.RandomStringWithMask(rand.MaskDigits, width)
	return smsCode
}
