package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *OrganizationAPIHandler) UpdateOrganizationPassword(ctx context.Context, req *pb.UpdateOrganizationPasswordRequest, rsp *pb.UpdateOrganizationPasswordResponse) error {
	err := validateUpdateOrganizationPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 检测密码强度
	ok := checkPassword(req.GetNewPlainPassword())
	if !ok {
		return errors.Error(codes.InvalidRequest, api.ErrorChineseMsg(api.ErrInvalidTenantPassword))
	}

	// 发送更新后台组织的密码请求
	_, err = s.userAPI.UpdateOrganizationPassword(ctx, &userpb.UpdateOrganizationPasswordRequest{
		// 组织ID
		OrganizationId: req.GetOrganizationId(),
		// 旧明文密码
		OldPlainPassword: req.GetOldPlainPassword(),
		// 新明文密码
		NewPlainPassword: req.GetNewPlainPassword(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 更新成功强制退出
	oid := req.GetOrganizationId()
	err = s.tokenStore.KickOutUser(ctx, oid)
	if err != nil {
		return errors.Error(api.ErrCreateAccessToken, api.ErrorChineseMsg(api.ErrKickedOut))
	}
	return nil
}

// 验证request
func validateUpdateOrganizationPasswordRequest(req *pb.UpdateOrganizationPasswordRequest) error {
	if req.GetOldPlainPassword() == "" {
		return gerr.New("old_plain_password should not be empty")
	}
	if req.GetNewPlainPassword() == "" {
		return gerr.New("new_plain_password should not be empty")
	}
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	return nil
}

// 检测密码是否满足强度要求
func checkPassword(pass string) bool {
	if len(pass) < 8 || len(pass) > 16 {
		return false
	}
	return true
}
