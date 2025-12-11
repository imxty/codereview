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

// 获取商户信息
func (u *UserAPIHandler) GetTenantEntity(ctx context.Context, req *pb.GetTenantEntityRequest, rsp *pb.GetTenantEntityResponse) error {
	// 验证 request
	err := validateGetTenantEntityRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 查询租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Error(ErrTenantNotFound, "tenant not found")
	}
	// 获取订阅时间线
	ts, err := u.userStore.GetSubscriptionTimeline(ctx, tenantEntity.GetTenantID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 判断商户状态
	var tenant_status pb.TenantStatus
	if ts == nil {
		// 没有订阅
		tenant_status = pb.TenantStatus_TENANT_STATUS_PENDING
	} else {
		// 检查订阅是否过期
		tenant_status = checkTenantSubscriptionTimeline(ts)
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
	// 未认证
	if !ot.GetIsActivated() {
		tenant_status = pb.TenantStatus_TENANT_STATUS_UNAUTH
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
	if revision != nil {
		// 返回数据
		rsp.Entity = toProtoEntityFromRevision(revision)
		rsp.Entity.TenantReviewStatus = toProtoTenantReviewStatus(ot.GetReviewStatus())
		rsp.Entity.TenantStatus = tenant_status
		rsp.Entity.OrganizationId = ot.GetOrganizationID()
		rsp.Entity.SafePhone = tenantEntity.GetSafePhone()
		rsp.Entity.FailReason = failReason
		rsp.Entity.Overdue = tenantEntity.GetOverdue()
		rsp.ReportSharingSwitchStatus = tenantEntity.GetReportSharingStatus()
		rsp.ConstitutionSwitchStatus = tenantEntity.GetConstitutionStatus()
		rsp.OverdueDays = tenantEntity.GetOverdue()
		return nil
	}
	// 返回数据
	rsp.Entity = toProtoEntity(tenantEntity)
	rsp.Entity.TenantReviewStatus = toProtoTenantReviewStatus(ot.GetReviewStatus())
	rsp.Entity.TenantStatus = tenant_status
	rsp.Entity.OrganizationId = ot.GetOrganizationID()
	rsp.Entity.FailReason = failReason
	rsp.Entity.Overdue = tenantEntity.GetOverdue()
	rsp.ReportSharingSwitchStatus = tenantEntity.GetReportSharingStatus()
	rsp.ConstitutionSwitchStatus = tenantEntity.GetConstitutionStatus()
	rsp.OverdueDays = tenantEntity.GetOverdue()
	return nil
}

// 验证 request
func validateGetTenantEntityRequest(req *pb.GetTenantEntityRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
