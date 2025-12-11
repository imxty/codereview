package user

import (
	"context"
	gerr "errors"

	reviewv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取组织下的商户
func (s *UserAPIHandler) GetOrganizationTenants(ctx context.Context, req *pb.GetOrganizationTenantsRequest, rsp *pb.GetOrganizationTenantsResponse) error {
	err := validateGetOrganizationTenantsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取组织商户
	ots, err := s.userStore.ListOrganizationExistTenants(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	tids := make([]string, len(ots))
	tenants_map := make(map[string]domain.OrganizationTenantIntf)
	review_failed_tids := make([]string, 0)
	for k, v := range ots {
		tids[k] = v.GetTenantID()
		tenants_map[v.GetTenantID()] = v
		// 审核失败商户
		if v.GetReviewStatus() == domain.TenantReviewStatusReviewFailed {
			review_failed_tids = append(review_failed_tids, v.GetTenantID())
		}
	}
	// 查询商户的信息
	entities, err := s.userStore.ListTenantEntity(ctx, tids)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 查询商户订阅
	ts, err := s.userStore.ListSubscriptionTimeline(ctx, tids)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	ts_map := make(map[string]domain.SubscriptionTimelineIntf)
	for _, each := range ts {
		ts_map[each.GetTenantID()] = each
	}

	// 获取审核失败的商户的失败原因
	fail_results := make(map[string]*reviewv1.ReviewResult)
	if len(review_failed_tids) > 0 {
		getRsp, err := s.reviewAPI.ListCertificateReviewStatus(ctx, &reviewv1.ListCertificateReviewStatusRequest{
			TenantIds: review_failed_tids,
		})
		if err != nil {
			return err
		}
		fail_results = getRsp.GetResults()
	}

	// 返回所有商户信息
	protoEntity := make([]*pb.TenantEntity, len(entities))
	for k, v := range entities {
		protoEntity[k] = toProtoTenantEntity(v)
		// 获取商户审核状态
		protoEntity[k].TenantReviewStatus = toProtoTenantReviewStatus(tenants_map[v.GetTenantID()].GetReviewStatus())
		// 如果审核失败，返回失败原因
		if tenants_map[v.GetTenantID()].GetReviewStatus() == domain.TenantReviewStatusReviewFailed {
			protoEntity[k].FailReason = fail_results[v.GetTenantID()].GetFailReason()
		}
		// 判断商户状态
		if tenants_map[v.GetTenantID()].GetIsActivated() {
			// 已认证
			if v, ok := ts_map[v.GetTenantID()]; !ok {
				// 认证但没有时间线，待续期状态
				protoEntity[k].TenantStatus = pb.TenantStatus_TENANT_STATUS_PENDING
			} else {
				// 认证且有时间线，判断是否过期
				protoEntity[k].TenantStatus = checkTenantSubscriptionTimeline(v)
			}
		} else {
			// 未认证
			protoEntity[k].TenantStatus = pb.TenantStatus_TENANT_STATUS_UNAUTH
		}
	}

	rsp.Tenants = protoEntity
	return nil
}

func validateGetOrganizationTenantsRequest(req *pb.GetOrganizationTenantsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	return nil
}
