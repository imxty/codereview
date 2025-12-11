package organization

import (
	"context"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 查询商户续费明细
func (s *OrganizationAPIHandler) SearchTenantSubscriptions(ctx context.Context, req *pb.SearchTenantSubscriptionsRequest, rsp *pb.SearchTenantSubscriptionsResponse) error {
	err := validateSearchTenantSubscriptionsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchRsp, err := s.userAPI.SearchTenantSubscriptions(ctx, &userv1.SearchTenantSubscriptionsRequest{
		OrganizationId: req.GetOrganizationId(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	results := make([]*pb.TenantSubscriptionTimeline, len(searchRsp.GetSubcriptions()))
	for k, v := range searchRsp.GetSubcriptions() {
		var status pb.TenantStatus
		if time.Now().Before(v.GetExpiredAt().AsTime()) {
			status = pb.TenantStatus_TENANT_STATUS_USING
		} else {
			status = pb.TenantStatus_TENANT_STATUS_PENDING
		}
		results[k] = &pb.TenantSubscriptionTimeline{
			TenantId:   v.GetTenantId(),
			StartTime:  v.GetCreatedAt(),
			EndTime:    v.GetExpiredAt(),
			TenantName: v.GetTenantName(),
			Years:      v.GetYears(),
			Status:     status,
		}
	}

	rsp.Subcriptions = results

	return nil
}

// 验证request
func validateSearchTenantSubscriptionsRequest(req *pb.SearchTenantSubscriptionsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be nil")
	}
	if req.GetEndTime() == nil {
		return gerr.New("end_time should not be nil")
	}
	return nil
}
