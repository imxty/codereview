package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) UpdateAdminPassword(ctx context.Context, req *pb.UpdateAdminPasswordRequest, rsp *pb.UpdateAdminPasswordResponse) error {
	err := validateUpdateAdminPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 检测密码强度
	ok := checkPassword(req.GetNewPlainPassword())
	if !ok {
		return errors.Error(codes.InvalidRequest, api.ErrorChineseMsg(api.ErrInvalidTenantPassword))
	}

	// 发送更新后台商户的密码请求
	_, err = s.userAPI.UpdateAdminPassword(ctx, &userpb.UpdateAdminPasswordRequest{
		// 租户ID
		TenantId: req.GetTenantId(),
		// 旧明文密码
		OldPlainPassword: req.GetOldPlainPassword(),
		// 新明文密码
		NewPlainPassword: req.GetNewPlainPassword(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 更新成功强制退出
	tenantID := req.GetTenantId()
	err = s.tokenStore.KickOutUser(ctx, tenantID)
	if err != nil {
		return errors.Error(api.ErrCreateAccessToken, api.ErrorChineseMsg(api.ErrKickedOut))
	}
	return nil
}

// 验证request
func validateUpdateAdminPasswordRequest(req *pb.UpdateAdminPasswordRequest) error {
	if req.GetOldPlainPassword() == "" {
		return gerr.New("old_plain_password should not be empty")
	}
	if req.GetNewPlainPassword() == "" {
		return gerr.New("new_plain_password should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	return nil
}
