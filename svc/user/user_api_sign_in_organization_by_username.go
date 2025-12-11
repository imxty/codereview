package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrOrganizationNotExist
	ErrOrganizationNotExist = 5001
	// ErrWrongPassword
	ErrWrongPassword = 5014
)

// 通过账号登录组织管理页面请求
func (u *UserAPIHandler) SignInOrganizationByUsername(ctx context.Context, req *pb.SignInOrganizationByUsernameRequest, rsp *pb.SignInOrganizationByUsernameResponse) error {
	// 验证 request
	err := validateSignInOrganizationByUsernameRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 通过用户名查找组织
	organization, err := u.userStore.GetOrganizationByUsername(ctx, req.GetUsername())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if organization == nil {
		return errors.Errorf(ErrOrganizationNotExist, "failed to find organization by username [%s]", req.GetUsername())
	}

	// 生成 hash 密码
	hashedPassword := generateHashSHA256(req.GetPassword())
	// 判断密码是否相同
	if organization.GetHashedPassword() != hashedPassword {
		return errors.Error(ErrWrongPassword, "password wrong")
	}

	// 查询组织下的商户
	tf, _, err := u.userStore.ListOrganizationTenantsWithoutTenantId(ctx, organization.GetOrganizationID(), 0, 1)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	hasTenant := false
	if len(tf) != 0 {
		hasTenant = true
	}

	// 返回组织信息
	org := &pb.Organization{
		// 组织 ID
		OrganizationId: organization.GetOrganizationID(),
		// 名称
		Name: organization.GetUsername(),
		// 手机号
		Phone:     organization.GetPhone(),
		HasTenant: hasTenant,
	}
	rsp.Organization = org

	return nil
}

// 验证 request
func validateSignInOrganizationByUsernameRequest(req *pb.SignInOrganizationByUsernameRequest) error {
	if req.GetUsername() == "" {
		return gerr.New("username should not be empty")
	}
	if req.GetPassword() == "" {
		return gerr.New("password should not be empty")
	}
	return nil
}
