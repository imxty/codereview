package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/huimaibao-service/api/token/tokenstore"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) SignOutTenantAdmin(ctx context.Context, req *pb.SignOutTenantAdminRequest, rsp *pb.SignOutTenantAdminResponse) error {
	// 验证request
	err := validateSignOutTenantAdminRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 验证token，并获取token的详细信息
	rtDetails, err := s.tokenStore.GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return errors.ErrorWithCause(api.ErrInvalidRefreshToken, err, api.ErrorMsg(err))
	}
	// 如果用户ID不同则返回错误
	if rtDetails.UserInfo[tokenstore.UserID] != req.GetTenantId() {
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
func validateSignOutTenantAdminRequest(req *pb.SignOutTenantAdminRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetAccessToken() == "" {
		return gerr.New("access_token should not be empty")
	}
	if req.GetRefreshToken() == "" {
		return gerr.New("refresh_token should not be empty")
	}
	return nil
}
