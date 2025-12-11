package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/huimaibao-service/api/token/tokenstore"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *BossAPIHandler) RefreshAccessToken(ctx context.Context, req *pb.RefreshAccessTokenRequest, rsp *pb.RefreshAccessTokenResponse) error {
	err := validateRefreshAccessTokenRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 验证token，并获取token的详细信息
	rtDetails, err := s.tokenStore.GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return errors.ErrorWithCause(api.ErrInvalidRefreshToken, err, api.ErrorMsg(err))
	}
	// 如果用户ID不同则返回错误
	if rtDetails.UserInfo[tokenstore.UserID] != req.GetUserId() {
		return errors.ErrorWithCause(api.ErrInvalidRefreshToken, err, api.ErrorMsg(err))
	}
	// 创建新的accessToken
	atDetails, err := s.tokenStore.CreateAccessToken(ctx, rtDetails.UserInfo)
	if err != nil {
		return errors.ErrorWithCause(api.ErrInvalidRefreshToken, err, api.ErrorMsg(err))
	}
	// 刷新refreshToken时间
	err = s.tokenStore.RenewRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return errors.ErrorWithCause(api.ErrInvalidRefreshToken, err, api.ErrorMsg(err))
	}
	// 时间的转化
	expire := timestamppb.New(atDetails.ExpiredAt)
	accessToken := &pb.TokenDetails{
		Token:       atDetails.AccessToken,
		ExpiredTime: expire,
	}
	rsp.AccessToken = accessToken

	return nil
}

// 验证request
func validateRefreshAccessTokenRequest(req *pb.RefreshAccessTokenRequest) error {
	if req.GetRefreshToken() == "" {
		return gerr.New("refresh token should not be empty")
	}
	if req.GetUserId() == "" {
		return gerr.New("user id should not be empty")
	}
	return nil
}
