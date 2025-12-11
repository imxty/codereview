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

// 创建商户
func (s *OrganizationAPIHandler) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest, rsp *pb.CreateTenantResponse) error {
	err := validateCreateTenantRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 创建商户
	createRsp, err := s.userAPI.CreateTenant(ctx, &userv1.CreateTenantRequest{
		OrganizationId: req.GetOrganizationId(),
		Tenant:         toSvcTenant(req.GetTenant()),
		PlainPassword:  req.GetPlainPassword(),
		SmsCode:        req.GetSmsCode(),
		TxId:           req.GetTxId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回响应
	rsp.Tenant = toApiTenant(createRsp.GetTenant(), s.s3Domain)

	return nil
}

// 验证request
func validateCreateTenantRequest(req *pb.CreateTenantRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetTenant() == nil {
		return gerr.New("tenant should not be empty")
	}
	tenant := req.GetTenant()
	if tenant.GetContactPhone() == "" {
		return gerr.New("contact phone should not be empty")
	}
	if tenant.GetContactName() == "" {
		return gerr.New("contact_name should not be empty")
	}
	if tenant.GetSafePhone() == "" {
		return gerr.New("safe phone should not be empty")
	}
	if tenant.GetLogoUrl() == "" {
		return gerr.New("logo should not be empty")
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
	if req.GetPlainPassword() == "" {
		return gerr.New("plain_password should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms_code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx_id should not be empty")
	}
	return nil
}
