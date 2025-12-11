package user

import (
	"context"
	gerr "errors"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 通过手机号登录组织管理页面请求
func (u *UserAPIHandler) SignInOrganizationByPhone(ctx context.Context, req *pb.SignInOrganizationByPhoneRequest, rsp *pb.SignInOrganizationByPhoneResponse) error {
	// 验证 request
	err := validateSignInOrganizationByPhoneRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 验证手机短信
	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   req.GetPhone(),
		TxId:    req.GetTxId(),
		SmsCode: req.GetSmsCode(),
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_ORGANIZATION,
	})
	if err != nil {
		return err
	}

	// 查找组织
	organization, err := u.userStore.GetOrganizationByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if organization == nil {
		return errors.Errorf(ErrOrganizationNotExist, "failed to find organization by phone [%s]", req.GetPhone())
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
func validateSignInOrganizationByPhoneRequest(req *pb.SignInOrganizationByPhoneRequest) error {
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
