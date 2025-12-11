package notification

import (
	"context"
	"encoding/json"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/svc/notification/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *NotificationAPIHandler) GetLatestPhoneVerificationCode(ctx context.Context, req *pb.GetLatestPhoneVerificationCodeRequest, rsp *pb.GetLatestPhoneVerificationCodeResponse) error {
	err := validateGetLatestPhoneVerificationCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	sms, err := s.notificationStore.GetLatestSms(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 如果没有短信直接返回空即可
	if sms == nil {
		return nil
	}

	// 解析 TemplateParam，获得纯验证码
	var st domain.SmsTemplateParam
	err = json.Unmarshal([]byte(sms.GetTemplateParam()), &st)
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}

	rsp.SmsCode = st.Code
	rsp.TxId = sms.GetTxID()
	return nil
}

// 验证 request
func validateGetLatestPhoneVerificationCodeRequest(req *pb.GetLatestPhoneVerificationCodeRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	return nil
}
