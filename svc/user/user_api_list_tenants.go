package user

import (
	"context"
	gerr "errors"

	productv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取商户列表
func (u *UserAPIHandler) ListTenants(ctx context.Context, req *pb.ListTenantsRequest, rsp *pb.ListTenantsResponse) error {
	// 验证 request
	err := validateListTenantsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取商户组织列表
	ots, count, err := u.userStore.ListOrganizationTenantsWithoutTenantId(ctx, req.GetOrganizationId(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获取商户订阅信息
	tids := make([]string, 0, len(ots))
	for _, each := range ots {
		tids = append(tids, each.GetTenantID())
	}

	ts, err := u.userStore.ListSubscriptionTimeline(ctx, tids)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	ts_map := make(map[string]domain.SubscriptionTimelineIntf)
	pts_map := make(map[string]*pb.TenantSubscriptionTimeline)
	for _, each := range ts {
		ts_map[each.GetTenantID()] = each
		pts_map[each.GetTenantID()] = toProtoSubscriptionTimeline(each)
	}

	// 获取商户方案
	trs, err := u.userStore.ListTreatmentsByTenantIDs(ctx, tids)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	tr_ids := make([]string, len(trs))
	tr_tenant_map := make(map[string]string)
	for k, v := range trs {
		tr_ids[k] = v.GetTreatmentID()
		tr_tenant_map[v.GetTenantID()] = v.GetTreatmentID()
	}

	// 获取方案名称
	treatMap := make(map[string]string)
	if len(tr_ids) > 0 {
		listRsp, err := u.productAPI.ListTreatmentNameByIds(ctx, &productv1.ListTreatmentNameByIdsRequest{
			TreatmentIds: tr_ids,
		})
		if err != nil {
			return err
		}
		for k, v := range tr_tenant_map {
			treatMap[k] = listRsp.GetTreatmentNames()[v]
		}
	}

	// 审核失败 ID 集合
	reviewFailTenants := []string{}
	tenantEntity := make([]*pb.TenantEntity, len(ots))
	for k, v := range ots {
		entity, err := u.userStore.GetTenantEntity(ctx, v.GetTenantID())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		if entity == nil {
			return errors.Errorf(ErrTenantNotFound, "tenant[%s] not found", v.GetTenantID())
		}
		var tenant_status pb.TenantStatus
		if !v.GetIsActivated() {
			// 未认证
			tenant_status = pb.TenantStatus_TENANT_STATUS_UNAUTH
		} else {
			// 判断商户订阅状态
			if t, ok := ts_map[entity.GetTenantID()]; !ok {
				// 没有订阅，待续期状态
				tenant_status = pb.TenantStatus_TENANT_STATUS_PENDING
			} else {
				// 有订阅，判断有没有过期
				tenant_status = checkTenantSubscriptionTimeline(t)
			}
		}
		// 判断当前 tenant 的状态
		// 如果是审核状态则去获取(revision)，如果失败还要去获取失败原因
		// 其他情况获取 entity 即可
		if v.GetReviewStatus() != domain.TenantReviewStatusReviewSuccess {
			revision, err := u.userStore.GetLatestTenantEntityRevision(ctx, v.GetTenantID())
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
			// 待审核但没有 revision，不合法
			if revision == nil {
				return errors.Errorf(ErrTenantNotFound, "tenant[%s] has no revision", v.GetTenantID())
			}
			tenantEntity[k] = toProtoTenantEntityFromRevision(revision)
			tenantEntity[k].TenantStatus = tenant_status
			tenantEntity[k].TenantReviewStatus = toProtoTenantReviewStatus(v.GetReviewStatus())
			if v.GetReviewStatus() == domain.TenantReviewStatusReviewFailed {
				// 获取失败原因
				reviewFailTenants = append(reviewFailTenants, v.GetTenantID())
			}
		} else {
			// 其他情况获取 tenantEntity 即可
			tenantEntity[k] = toProtoTenantEntity(entity)
			tenantEntity[k].TenantStatus = tenant_status
			tenantEntity[k].TenantReviewStatus = toProtoTenantReviewStatus(v.GetReviewStatus())
		}
	}

	if len(reviewFailTenants) != 0 {
		// 获取失败原因
		listRsp, err := u.reviewAPI.ListCertificateReviewStatus(ctx, &reviewpb.ListCertificateReviewStatusRequest{
			TenantIds: reviewFailTenants,
		})
		if err != nil {
			return err
		}
		result := listRsp.GetResults()
		for k, v := range tenantEntity {
			if res, ok := result[v.GetTenantId()]; ok {
				tenantEntity[k].FailReason = res.GetFailReason()
			}
		}
	}

	rsp.Tenants = tenantEntity
	rsp.Subscriptions = pts_map
	rsp.Treatments = treatMap
	rsp.TotalCount = int32(count)

	return nil
}

// 验证 request
func validateListTenantsRequest(req *pb.ListTenantsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
