package user

import (
	"context"
	gerr "errors"

	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取首次认证失败的商户信息
func (u *UserAPIHandler) GetUnAuthTenantRevision(ctx context.Context, req *pb.GetUnAuthTenantRevisionRequest, rsp *pb.GetUnAuthTenantRevisionResponse) error {
	// 验证request
	err := validateGetUnAuthTenantRevisionRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 查询租户信息
	tenantUser, err := u.userStore.GetTenantUser(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantUser == nil {
		return errors.Error(ErrTenantNotFound, "tenant not found")
	}
	// 获取最新的副本更新情况
	revision, err := u.userStore.GetLatestTenantEntityRevision(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 获取商户组织
	ot, err := u.userStore.GetOrganizationTenantByTenantId(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if ot == nil {
		return errors.Errorf(ErrOrganizationTenantNotExist, "organization tenant[%s] not found", req.GetTenantId())
	}

	// 如果不是审核失败且初次创建则不返回
	if ot.GetIsActivated() || ot.GetReviewStatus() != domain.TenantReviewStatusReviewFailed {
		rsp.HasUnauthRevision = false
		return nil
	}

	// 如果审核失败
	failReason := ""
	if ot.GetReviewStatus() == domain.TenantReviewStatusReviewFailed {
		// 获取失败原因
		listRsp, err := u.reviewAPI.ListCertificateReviewStatus(ctx, &reviewpb.ListCertificateReviewStatusRequest{
			TenantIds: []string{req.GetTenantId()},
		})
		if err != nil {
			return err
		}
		failReason = listRsp.GetResults()[req.GetTenantId()].GetFailReason()
	}

	// 返回数据
	rsp.HasUnauthRevision = true
	rsp.Tenant = toProtoTenantEntityFromRevision(revision)
	rsp.Tenant.TenantReviewStatus = toProtoTenantReviewStatus(ot.GetReviewStatus())
	rsp.Tenant.TenantStatus = pb.TenantStatus_TENANT_STATUS_UNAUTH
	rsp.Tenant.FailReason = failReason
	rsp.Tenant.OrganizationId = ot.GetOrganizationID()

	return nil
}

// 验证request
func validateGetUnAuthTenantRevisionRequest(req *pb.GetUnAuthTenantRevisionRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
