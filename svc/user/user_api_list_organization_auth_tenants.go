package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取组织下所有认证商户
func (s *UserAPIHandler) ListOrganizationAuthTenants(ctx context.Context, req *pb.ListOrganizationAuthTenantsRequest, rsp *pb.ListOrganizationAuthTenantsResponse) error {
	err := validateListOrganizationAuthTenantsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	ots, err := s.userStore.ListOrganizationAuthTenants(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	tenantIds := make([]string, len(ots))
	// 统计数量
	var authNumber int32 = 0
	for k, v := range ots {
		tenantIds[k] = v.GetTenantID()
		if v.GetReviewStatus() == domain.TenantReviewStatusReviewSuccess {
			authNumber++
		}
	}

	timelines, err := s.userStore.ListSubscriptionTimeline(ctx, tenantIds)
	if err != nil {
		return err
	}

	pbTimelines := make([]*pb.TenantSubscriptionTimeline, len(timelines))
	for i := range timelines {
		tenant, err := s.userStore.GetTenantEntity(ctx, timelines[i].GetTenantID())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		if tenant == nil {
			return errors.Errorf(ErrTenantNotFound, "tenant[%s] not found", timelines[i].GetTenantID())
		}
		pbTimelines[i] = toProtoSubscriptionTimeline(timelines[i])
		pbTimelines[i].TenantName = tenant.GetStoreName()
	}

	rsp.Timelines = pbTimelines
	rsp.AuthNumber = authNumber

	return nil
}

// 验证request
func validateListOrganizationAuthTenantsRequest(req *pb.ListOrganizationAuthTenantsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	return nil
}
