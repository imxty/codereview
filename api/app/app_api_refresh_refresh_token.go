package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/huimaibao-service/api/token/tokenstore"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AppAPIHandler) RefreshRefreshToken(ctx context.Context, req *pb.RefreshRefreshTokenRequest, rsp *pb.RefreshRefreshTokenResponse) error {

	// 1.验证request
	err := validateRefreshRefreshTokenRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 验证token，并获取token的详细信息
	rtDetails, err := s.tokenStore.GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return errors.ErrorWithCause(api.ErrInvalidRefreshToken, err, api.ErrorMsg(err))
	}
	// 如果用户ID不同则返回错误
	if rtDetails.UserInfo[tokenstore.UserID] != req.GetStaffId() {
		return errors.ErrorWithCause(api.ErrInvalidRefreshToken, err, api.ErrorMsg(err))
	}
	// 刷新token
	rtNewDetails, err := s.tokenStore.CreateRefreshToken(ctx, rtDetails.UserInfo)
	if err != nil {
		return errors.ErrorWithCause(api.ErrInvalidRefreshToken, err, api.ErrorMsg(err))
	}
	// 时间的转化
	expire := timestamppb.New(rtNewDetails.ExpiredAt)
	rt := &pb.TokenDetails{
		Token:       rtNewDetails.RefreshToken,
		ExpiredTime: expire,
	}
	rsp.RefreshToken = rt

	return nil
}

// 验证request
func validateRefreshRefreshTokenRequest(req *pb.RefreshRefreshTokenRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetRefreshToken() == "" {
		return gerr.New("refresh_token should not be empty")
	}
	return nil
}
