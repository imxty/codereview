package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/huimaibao-service/api/token/tokenstore"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 通过手机号登陆组织管理页面
func (s *OrganizationAPIHandler) SignInByPhone(ctx context.Context, req *pb.SignInByPhoneRequest, rsp *pb.SignInByPhoneResponse) error {
	err := validateSignInByPhoneRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	signRsp, err := s.userAPI.SignInOrganizationByPhone(ctx, &userv1.SignInOrganizationByPhoneRequest{
		Phone:   req.GetPhone(),
		SmsCode: req.GetSmsCode(),
		TxId:    req.GetTxId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 生成token
	organizationId := signRsp.GetOrganization().GetOrganizationId()
	userInfo := map[string]interface{}{
		tokenstore.UserID: organizationId,
	}
	at, err := s.tokenStore.CreateAccessToken(ctx, userInfo)
	if err != nil {
		return errors.Error(api.ErrCreateAccessToken, api.ErrorMsg(err))
	}
	rt, err := s.tokenStore.CreateRefreshToken(ctx, userInfo)
	if err != nil {
		return errors.ErrorWithCause(api.ErrCreateRefreshToken, err, api.ErrorMsg(err))
	}

	atExpireAt := timestamppb.New(at.ExpiredAt)
	rtExpireAt := timestamppb.New(rt.ExpiredAt)

	rsp.AccessToken = &pb.TokenDetails{
		Token:       at.AccessToken,
		ExpiredTime: atExpireAt,
	}
	rsp.RefreshToken = &pb.TokenDetails{
		Token:       rt.RefreshToken,
		ExpiredTime: rtExpireAt,
	}
	rsp.Organization = toApiOrganization(signRsp.GetOrganization())
	return nil
}

// 验证request
func validateSignInByPhoneRequest(req *pb.SignInByPhoneRequest) error {
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
