package user

import (
	"context"
	gerr "errors"

	reviewv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取商户
func (u *UserAPIHandler) GetTenant(ctx context.Context, req *pb.GetTenantRequest, rsp *pb.GetTenantResponse) error {
	err := validateGetTenantRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取商户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Errorf(ErrTenantNotFound, "tenant not found [%s]", req.GetTenantId())
	}

	var tenant_status pb.TenantStatus
	// 获取商户订阅时间线
	ts, err := u.userStore.GetSubscriptionTimeline(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if ts == nil {
		// 没有时间线
		tenant_status = pb.TenantStatus_TENANT_STATUS_PENDING
	} else {
		// 有时间线，检查是否过期
		tenant_status = checkTenantSubscriptionTimeline(ts)
	}

	// 获取组织商户信息
	ot, err := u.userStore.GetOrganizationTenant(ctx, req.GetOrganizationId(), req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 未认证状态
	if !ot.GetIsActivated() {
		tenant_status = pb.TenantStatus_TENANT_STATUS_UNAUTH
	}

	var fail_reason string
	// 如果审核失败，获取失败原因
	if ot.GetReviewStatus() == domain.ReviewStatusReviewFailed {
		re, err := u.reviewAPI.ListCertificateReviewStatus(ctx, &reviewv1.ListCertificateReviewStatusRequest{
			TenantIds: []string{req.GetTenantId()},
		})
		if err != nil {
			return err
		}
		fail_reason = re.GetResults()[req.GetTenantId()].GetFailReason()
	}

	// 返回商户数据
	rsp.Tenant = toProtoTenantEntity(tenantEntity)
	rsp.Tenant.TenantStatus = tenant_status
	rsp.Tenant.TenantReviewStatus = toProtoTenantReviewStatus(ot.GetReviewStatus())
	rsp.Tenant.FailReason = fail_reason

	return nil
}

// 验证 request
func validateGetTenantRequest(req *pb.GetTenantRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
