package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/meta"
)

func (s *AppAPIHandler) SignOut(ctx context.Context, req *pb.SignOutRequest, rsp *pb.SignOutResponse) error {

	// 1.验证request
	err := validateSignOutRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}
	// 从context获取access_token
	at, ok := meta.Get(ctx, ContextAccessToken)
	if !ok {
		return errors.Error(api.ErrUserUnauthorized, api.ErrorChineseMsg(api.ErrUserUnauthorized))
	}
	// 获取refresh_token
	rt := req.GetRefreshToken()
	userID, ok := meta.Get(ctx, ContextUserID)
	if !ok {
		return errors.Error(api.ErrUserUnauthorized, api.ErrorChineseMsg(api.ErrUserUnauthorized))
	}
	// 判断用户是否相同
	if userID != req.GetStaffId() {
		return errors.Error(api.ErrUserUnauthorized, api.ErrorChineseMsg(api.ErrUserUnauthorized))
	}
	// 删除token
	err = s.tokenStore.DeleteAccessToken(ctx, at)
	if err != nil {
		return errors.Error(api.ErrInvalidAccessToken, err.Error())
	}
	err = s.tokenStore.DeleteRefreshToken(ctx, rt)
	if err != nil {
		return errors.Error(api.ErrInvalidAccessToken, err.Error())
	}

	return nil
}

// 验证request
func validateSignOutRequest(req *pb.SignOutRequest) error {
	if req.GetStaffId() == "" {
		return gerr.New("staaff_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetRefreshToken() == "" {
		return gerr.New("refresh_token should not be empty")
	}
	return nil
}
