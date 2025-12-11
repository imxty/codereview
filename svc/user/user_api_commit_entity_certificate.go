package user

import (
	"context"
	gerr "errors"

	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

const (
	// ErrReviewPending
	ErrReviewPending = 5312
)

// 提交商户资质
func (u *UserAPIHandler) CommitEntityCertificate(ctx context.Context, req *pb.CommitEntityCertificateRequest, rsp *pb.CommitEntityCertificateResponse) error {
	err := validateCommitEntityCertificateRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 查找tenant
	tenant, err := u.userStore.GetTenantEntity(ctx, req.GetTenant().GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenant == nil {
		return errors.Errorf(ErrTenantNotFound, "tenant not found ]%s]", req.GetTenant().GetTenantId())
	}

	// 查询tenant的状态
	ot, err := u.userStore.GetOrganizationTenantByTenantId(ctx, tenant.GetTenantID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if ot == nil {
		return errors.Errorf(ErrOrganizationTenantNotExist, "organization tenant not found [%s]", tenant.GetTenantID())
	}

	// 如果是审核中则无法修改
	if ot.GetReviewStatus() == domain.TenantReviewStatusReviewing {
		return errors.Errorf(codes.InvalidRequest, "can not commit,tenant[%s] status is reviewing", tenant.GetTenantID())
	}
	// 查询tenant_revision
	revision, err := u.userStore.GetLatestTenantEntityRevision(ctx, req.GetTenant().GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 如果不存在说明商户id有问题
	if revision == nil {
		return errors.Errorf(codes.InvalidRequest, "revision not found,can not commit tenant entity[%s]", req.GetTenant().GetTenantId())
	}

	reqTenant := req.GetTenant()

	// 上传营业执照
	// 如果没有传就用原来的
	businessLicenseUrl := ""
	if reqTenant.GetBusinessLicenseUrl() == "" {
		// 如果有营业执照，上传获取URL
		if reqTenant.GetBusinessLicense() != nil {
			path, err := u.uploadImage(reqTenant.GetBusinessLicense())
			if err != nil {
				return errors.Error(codes.InvalidOperation, err.Error())
			}
			businessLicenseUrl = path
		} else {
			businessLicenseUrl = revision.GetBusinessLicenseUrl()
		}
	} else {
		businessLicenseUrl = checkLogoUrl(reqTenant.GetBusinessLicenseUrl())
		if businessLicenseUrl == "" {
			businessLicenseUrl = revision.GetBusinessLicenseUrl()
		}
	}

	// 上传logo
	// 如果没有传就用原来的
	logoUrl := checkLogoUrl(reqTenant.GetLogoUrl())
	if logoUrl == "" {
		// 如果有，则用最新提交的
		logoUrl = revision.GetLogoUrl()
	}
	// 创建新的revision
	nRevision := &domain.TenantEntityRevision{
		TenantEntityRevisionID: xid.New().String(),
		TenantID:               req.GetTenant().GetTenantId(),
		// 安全手机号无法提审，用最新的即可
		SafePhone:          tenant.GetSafePhone(),
		StoreName:          reqTenant.GetName(),
		Province:           reqTenant.GetAddress().GetProvince(),
		City:               reqTenant.GetAddress().GetCity(),
		District:           reqTenant.GetAddress().GetDistrict(),
		Street:             reqTenant.GetAddress().GetStreet(),
		ContactName:        reqTenant.GetContactName(),
		ContactPhone:       reqTenant.GetContactPhone(),
		SocialCreditCode:   reqTenant.GetSocialCreditCode(),
		BusinessLicenseUrl: businessLicenseUrl,
		LogoUrl:            logoUrl,
	}
	// 创建revision
	err = u.userStore.CreateTenantEntityRevision(ctx, nRevision)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	_, err = u.reviewAPI.CommitEntityCertificate(ctx, &reviewpb.CommitEntityCertificateRequest{
		// 组织ID
		OrganizationId: ot.GetOrganizationID(),
		// 商户ID
		TenantId: req.GetTenant().GetTenantId(),
		// 副本ID
		TenantEntityRevisionId: nRevision.GetTenantEntityRevisionID(),
		// 是否是组织提审
		IsOrganization: req.GetIsOrganization(),
	})
	if err != nil {
		return err
	}
	// 状态变审核中
	err = u.userStore.ChangeOrganizationTenantReviewStatus(ctx, tenant.GetTenantID(), domain.TenantReviewStatusReviewing, ot.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证request
func validateCommitEntityCertificateRequest(req *pb.CommitEntityCertificateRequest) error {
	if req.GetIsOrganization() && req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTenant() == nil {
		return gerr.New("tenant should not be nil")
	}
	if !req.GetIsOrganization() && req.GetTenant().GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
