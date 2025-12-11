package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) CommitEntityCertificate(ctx context.Context, req *pb.CommitEntityCertificateRequest, rsp *pb.CommitEntityCertificateResponse) error {
	err := validateCommitEntityCertificateRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	commitRsp, err := s.userAPI.CommitEntityCertificate(ctx, &userv1.CommitEntityCertificateRequest{
		// 商户信息
		Tenant: &userv1.TenantEntity{
			// 商户ID
			TenantId: req.GetTenantId(),
			// 商户名称
			Name: req.GetEntityName(),
			// 地址
			Address: toUserAddress(req.GetAddress()),
			// 联系人姓名
			ContactName: req.GetContactName(),
			// 营业执照
			BusinessLicense: toUserImage(req.GetBusinessLicense()),
			// 社会信用代码
			SocialCreditCode: req.GetSocialCreditCode(),
			// 联系人手机号
			ContactPhone: req.GetContactPhone(),
			// logo_url
			LogoUrl: req.GetLogoUrl(),
		},
		// 是否是组织提审
		IsOrganization: false,
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.Entity = toAppEntity(commitRsp.GetEntity(), s.s3Domain)
	return nil
}

// 验证request
func validateCommitEntityCertificateRequest(req *pb.CommitEntityCertificateRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetEntityName() == "" {
		return gerr.New("entity_name should not be empty")
	}
	if req.GetAddress() == nil {
		return gerr.New("address should not be empty")
	}
	address := req.GetAddress()
	if address.GetCity() == "" {
		return gerr.New("city should not be empty")
	}
	if address.GetProvince() == "" {
		return gerr.New("province should not be empty")
	}
	if address.GetStreet() == "" {
		return gerr.New("street should not be empty")
	}
	if req.GetContactName() == "" {
		return gerr.New("contact_name should not be empty")
	}
	if req.GetContactPhone() == "" {
		return gerr.New("contact_phone should not be empty")
	}
	if req.GetEntityName() == "" {
		return gerr.New("entity_name should not be empty")
	}
	if req.GetSocialCreditCode() == "" {
		return gerr.New("social_credit_code should not be empty")
	}
	return nil
}
