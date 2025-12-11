package user

import (
	"context"
	gerr "errors"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrStaffNotBelongTenant
	ErrStaffNotBelongTenant = 5018
)

// 通过验证码登录 App
func (u *UserAPIHandler) SignInAppBySmsCode(ctx context.Context, req *pb.SignInAppBySmsCodeRequest, rsp *pb.SignInAppBySmsCodeResponse) error {
	err := validateSignInAppBySmsCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 验证验证码
	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   req.GetPhone(),
		SmsCode: req.GetSmsCode(),
		TxId:    req.GetTxId(),
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_APP,
	})
	if err != nil {
		return err
	}

	// 	获取员工信息
	staff, err := u.userStore.GetStaffByPhone(ctx, "+86", req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if staff == nil {
		return errors.Error(ErrStaffNotExist, "staff not found")
	}
	// 员工不属于任何租户
	if !staff.GetIsActivated() {
		return errors.Error(ErrStaffNotBelongTenant, "staff not belong to any tenants")
	}
	// 获取租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, staff.GetTenantID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Error(ErrTenantNotFound, "tenant not found")
	}
	// 返回数据
	sf := &pb.Staff{
		// 员工 ID
		StaffId: staff.GetUserID(),
		// 员工姓名
		Name: staff.GetNickname(),
		// 员工手机号
		Phone: staff.GetPhone(),
		// 员工状态
		IsActivated: staff.GetIsActivated(),
		// 是否已删除
		IsDeleted: !staff.GetIsActivated(),
	}
	rsp.Staff = sf
	rsp.Staff.Tenant = toProtoTenant(tenantEntity)
	return nil
}

// 验证 request
func validateSignInAppBySmsCodeRequest(req *pb.SignInAppBySmsCodeRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms_code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx_id should not be empty")
	}
	return nil
}
