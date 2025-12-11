package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 首次认证失败再次创建商户
func (s *OrganizationAPIHandler) ReCreateTenant(ctx context.Context, req *pb.ReCreateTenantRequest, rsp *pb.ReCreateTenantResponse) error {
	err := validateReCreateTenantRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.ReCreateTenant(ctx, &userv1.ReCreateTenantRequest{
		// 商户信息
		Tenant: toSvcTenant(req.GetTenant()),
		// 密码
		PlainPassword: req.GetPlainPassword(),
		// 验证码
		SmsCode: req.GetSmsCode(),
		// 凭证 ID
		TxId: req.GetTxId(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

// 验证request
func validateReCreateTenantRequest(req *pb.ReCreateTenantRequest) error {
	if req.GetTenant().GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	tenant := req.GetTenant()
	if tenant.GetAddress() == nil {
		return gerr.New("address should not be empty")
	}
	if tenant.GetSocialCreditCode() == "" {
		return gerr.New("social credit code should not be empty")
	}
	if tenant.GetAddress().GetCity() == "" {
		return gerr.New("city should not be empty")
	}
	if tenant.GetAddress().GetProvince() == "" {
		return gerr.New("province should not be empty")
	}
	if tenant.GetAddress().GetStreet() == "" {
		return gerr.New("street should not be empty")
	}
	if tenant.GetName() == "" {
		return gerr.New("name should not be empty")
	}
	if tenant.GetContactPhone() == "" {
		return gerr.New("contact_phone should not be empty")
	}
	if tenant.GetContactName() == "" {
		return gerr.New("contact_name should not be empty")
	}
	if tenant.GetSafePhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx_id should not be empty")
	}
	if req.GetPlainPassword() == "" {
		return gerr.New("password should not be empty")
	}
	return nil
}
