package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	notificationv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) VerifyPhoneVerificationCode(ctx context.Context, req *pb.VerifyPhoneVerificationCodeRequest, rsp *pb.VerifyPhoneVerificationCodeResponse) error {
	err := validateVerifyPhoneVerificationCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	getRsp, err := s.notificationAPI.GetLatestPhoneVerificationCode(ctx, &notificationv1.GetLatestPhoneVerificationCodeRequest{
		Phone: req.GetPhone(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	if getRsp.GetTxId() != req.GetTxId() {
		return errors.Errorf(codes.InvalidRequest, "tx id not same[%s %s]", req.GetTxId(), getRsp.GetTxId())
	}
	if req.GetSmsCode() != getRsp.GetSmsCode() {
		return errors.Errorf(codes.InvalidRequest, "sms code not same[%s %s]", req.GetSmsCode(), getRsp.GetSmsCode())
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
	return nil
}
