package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/huimaibao-service/api/token/tokenstore"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *BossAPIHandler) SignOut(ctx context.Context, req *pb.SignOutRequest, rsp *pb.SignOutResponse) error {
	// 验证request
	err := validateSignOutRequest(req)
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

	// 删除token
	err = s.tokenStore.DeleteAccessToken(ctx, req.GetAccessToken())
	if err != nil {
		return errors.Error(api.ErrInvalidAccessToken, err.Error())
	}
	err = s.tokenStore.DeleteRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return errors.Error(api.ErrInvalidAccessToken, err.Error())
	}
	return nil
}

// 验证request
func validateSignOutRequest(req *pb.SignOutRequest) error {
	if req.GetAccessToken() == "" {
		return gerr.New("access token should not be empty")
	}
	if req.GetRefreshToken() == "" {
		return gerr.New("refresh token should not be empty")
	}
	if req.GetUserId() == "" {
		return gerr.New("user id should not be empty")
	}
	return nil
}
