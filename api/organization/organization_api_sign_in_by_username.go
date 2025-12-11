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

func (s *OrganizationAPIHandler) SignInByUsername(ctx context.Context, req *pb.SignInByUsernameRequest, rsp *pb.SignInByUsernameResponse) error {
	err := validateSignInByUsernameRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	signRsp, err := s.userAPI.SignInOrganizationByUsername(ctx, &userv1.SignInOrganizationByUsernameRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
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
func validateSignInByUsernameRequest(req *pb.SignInByUsernameRequest) error {
	if req.GetPassword() == "" {
		return gerr.New("password should not be empty")
	}
	if req.GetUsername() == "" {
		return gerr.New("username should not be empty")
	}
	return nil
}
