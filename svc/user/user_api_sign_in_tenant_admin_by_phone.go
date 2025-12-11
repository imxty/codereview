package user

import (
	"context"
	gerr "errors"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 通过手机号登录商户管理页面请求
func (u *UserAPIHandler) SignInTenantAdminByPhone(ctx context.Context, req *pb.SignInTenantAdminByPhoneRequest, rsp *pb.SignInTenantAdminByPhoneResponse) error {
	// 验证 request
	err := validateSignInTenantAdminByPhoneRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 验证手机短信
	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   req.GetPhone(),
		TxId:    req.GetTxId(),
		SmsCode: req.GetSmsCode(),
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_TENANT,
	})
	if err != nil {
		return err
	}

	// 获取商户 user
	t, err := u.userStore.GetTenantUserByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if t == nil {
		return errors.Errorf(ErrTenantNotFound, "tenant phone[%s] not found", req.GetPhone())
	}

	// 获取商户状态
	ot, err := u.userStore.GetOrganizationTenantByTenantId(ctx, t.GetTenantID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if ot == nil {
		return errors.Errorf(ErrOrganizationTenantNotExist, "organization tenant not found [%s]", t.GetTenantID())
	}
	// 如果还是未认证则无法登录
	if !ot.GetIsActivated() {
		return errors.Errorf(ErrTenantUnAuth, "unauth tenant[%s]", t.GetTenantID())
	}

	// 获取租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, t.GetTenantID())
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
func validateSignInTenantAdminByPhoneRequest(req *pb.SignInTenantAdminByPhoneRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx id should not be empty")
	}
	return nil
}
