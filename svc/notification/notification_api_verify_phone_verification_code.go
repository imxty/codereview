package notification

import (
	"context"
	"encoding/json"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/svc/notification/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrSmsCodeError
	ErrSmsCodeError = 5042
	// ErrSmsExpired
	ErrSmsExpired = 5308
)

func (s *NotificationAPIHandler) VerifyPhoneVerificationCode(ctx context.Context, req *pb.VerifyPhoneVerificationCodeRequest, rsp *pb.VerifyPhoneVerificationCodeResponse) error {
	err := validateVerifyPhoneVerificationCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	now := time.Now().UTC()

	// 通过 txid 获取短信
	sms, err := s.notificationStore.GetSmsByTxId(ctx, req.GetTxId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if sms == nil {
		return errors.Errorf(ErrSmsNotFound, "sms[%s] not found", req.GetTxId())
	}

	if req.GetAction().String() != sms.GetTemplateAction() || req.GetPhone() != sms.GetPhone() {
		return errors.Errorf(codes.InvalidRequest, "invalid sms[%v]", req)
	}
	// 验证短信验证码
	// 解析 TemplateParam，获得纯验证码
	var st domain.SmsTemplateParam
	err = json.Unmarshal([]byte(sms.GetTemplateParam()), &st)
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}
	if st.Code != req.GetSmsCode() {
		return errors.Errorf(ErrSmsCodeError, "sms code error[%s,%s]", req.GetSmsCode(), st.Code)
	}
	// 判断是否过期
	if sms.GetCreatedAt().Add(time.Minute * 10).UTC().After(now) {
		// 已经过期直接返回异常即可
		return errors.Errorf(ErrSmsExpired, "sms[%s] expired", req.GetTxId())
	}

	return nil
}

// 验证 request
func validateVerifyPhoneVerificationCodeRequest(req *pb.VerifyPhoneVerificationCodeRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx id should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms code should not be empty")
	}
	if req.GetAction() == pb.TemplateAction_TEMPLATE_ACTION_INVALID || req.GetAction() == pb.TemplateAction_TEMPLATE_ACTION_UNSET {
		return gerr.New("invalid template_action")
	}
	return nil
}
