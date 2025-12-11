package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
)

func (s *AppAPIHandler) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest, rsp *pb.UpdatePasswordResponse) error {

	// 1.验证request
	err := validateUpdatePasswordRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 密码不能少于6位
	if len(req.GetNewPlainPassword()) < 6 {
		return errors.Error(api.ErrStaffPasswordLength, api.ErrorChineseMsg(api.ErrStaffPasswordLength))
	}

	// 发送更新密码请求
	_, err = s.userAPI.UpdatePassword(ctx, &userpb.UpdatePasswordRequest{
		StaffId:          req.GetStaffId(),
		TenantId:         req.GetTenantId(),
		OldPlainPassword: req.GetOldPlainPassword(),
		NewPlainPassword: req.GetNewPlainPassword(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 更新成功强制退出
	err = s.tokenStore.KickOutUser(ctx, req.GetStaffId())
	if err != nil {
		return errors.Error(api.ErrCreateAccessToken, api.ErrorChineseMsg(api.ErrKickedOut))
	}

	return nil
}

// 验证request
func validateUpdatePasswordRequest(req *pb.UpdatePasswordRequest) error {
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetOldPlainPassword() == "" {
		return gerr.New("old_plain_password should not be empty")
	}
	if req.GetNewPlainPassword() == "" {
		return gerr.New("new_plain_password should not be empty")
	}
	return nil
}
