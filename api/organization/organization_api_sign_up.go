package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *OrganizationAPIHandler) SignUp(ctx context.Context, req *pb.SignUpRequest, rsp *pb.SignUpResponse) error {
	err := validateSignUpRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	signRsp, err := s.userAPI.SignUpOrganization(ctx, &userv1.SignUpOrganizationRequest{
		Phone:         req.GetPhone(),
		PlainPassword: req.GetPlainPassword(),
		SmsCode:       req.GetSmsCode(),
		TxId:          req.GetTxId(),
		Username:      req.GetUsername(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	rsp.Organization = toApiOrganization(signRsp.GetOrganization())
	return nil
}

// 验证request
func validateSignUpRequest(req *pb.SignUpRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetPlainPassword() == "" {
		return gerr.New("plain password should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx id should not be empty")
	}
	if req.GetUsername() == "" {
		return gerr.New("username should not be empty")
	}
	return nil
}
