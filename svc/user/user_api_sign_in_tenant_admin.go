package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrOrganizationTenantNotExist
	ErrOrganizationTenantNotExist = 5002
	// ErrTenantUnAuth
	ErrTenantUnAuth = 5307
)

// 登录商户后台
func (u *UserAPIHandler) SignInTenantAdmin(ctx context.Context, req *pb.SignInTenantAdminRequest, rsp *pb.SignInTenantAdminResponse) error {
	// 验证 request
	err := validateSignInTenantAdminRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询租户
	hashedPassword := generateHashSHA256(req.GetPassword())
	tenant, err := u.userStore.GetTenantUserByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenant == nil {
		return errors.Error(ErrTenantNotFound, "tenant not found")
	}

	// 获取商户状态
	ot, err := u.userStore.GetOrganizationTenantByTenantId(ctx, tenant.GetTenantID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if ot == nil {
		return errors.Errorf(ErrOrganizationTenantNotExist, "organization tenant not found [%s]", tenant.GetTenantID())
	}
	// 如果还是未认证则无法登录
	if !ot.GetIsActivated() {
		return errors.Errorf(ErrTenantUnAuth, "unauth tenant[%s]", tenant.GetTenantID())
	}

	// 对比密码
	if tenant.GetHashedPassword() != hashedPassword {
		return errors.Error(ErrWrongPassword, "wrong tenant password")
	}

	// 获取租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, tenant.GetTenantID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Error(ErrTenantNotFound, "tenant_entity not found")
	}

	rsp.Entity = toProtoEntity(tenantEntity)
	rsp.ReportSharingSwitchStatus = tenantEntity.GetReportSharingStatus()

	return nil
}

// 验证 request
func validateSignInTenantAdminRequest(req *pb.SignInTenantAdminRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("tenant id should not be empty")
	}
	if req.GetPassword() == "" {
		return gerr.New("password should not be empty")
	}
	return nil
}
