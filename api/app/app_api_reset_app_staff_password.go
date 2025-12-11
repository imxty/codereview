package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) ResetAppStaffPassword(ctx context.Context, req *pb.ResetAppStaffPasswordRequest, rsp *pb.ResetAppStaffPasswordResponse) error {
	err := validateResetAppStaffPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.ResetAppStaffPassword(ctx, &userv1.ResetAppStaffPasswordRequest{
		Phone:            req.GetPhone(),
		NewPlainPassword: req.GetNewPlainPassword(),
		SmsCode:          req.GetSmsCode(),
		TxId:             req.GetTxId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

func validateResetAppStaffPasswordRequest(req *pb.ResetAppStaffPasswordRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetNewPlainPassword() == "" {
		return gerr.New("new_plain_password should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms_code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx_id should not be empty")
	}
	return nil
}
