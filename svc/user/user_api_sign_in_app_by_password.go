package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 使用密码登录 App
func (u *UserAPIHandler) SignInAppByPassword(ctx context.Context, req *pb.SignInAppByPasswordRequest, rsp *pb.SignInAppByPasswordResponse) error {
	err := validateSignInAppByPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 生成 hash 密码
	hashedPassword := generateHashSHA256(req.GetPlainPassword())
	// 查询用户
	user, err := u.userStore.GetStaffByPhone(ctx, "+86", req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 员工不存在
	if user == nil {
		return errors.Error(ErrStaffNotExist, "user doesn't exist")
	}
	// 判断密码是否相同
	if user.GetHashedPassword() != hashedPassword {
		return errors.Error(ErrWrongPassword, "password wrong")
	}
	// 员工不属于任何租户
	if !user.GetIsActivated() {
		return errors.Error(ErrStaffNotBelongTenant, "staff not belong to any tenants")
	}
	// 获取租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, user.GetTenantID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Error(ErrTenantNotFound, "tenant not found")
	}

	// 返回数据
	rsp.Staff = toProtoStaff(user)
	rsp.Staff.Tenant = toProtoTenant(tenantEntity)

	return nil
}

// 验证 request
func validateSignInAppByPasswordRequest(req *pb.SignInAppByPasswordRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetPlainPassword() == "" {
		return gerr.New("plain_password should not be empty")
	}
	return nil
}
