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

// 搜索商户
func (s *OrganizationAPIHandler) SearchTenants(ctx context.Context, req *pb.SearchTenantsRequest, rsp *pb.SearchTenantsResponse) error {
	err := validateSearchTenantsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchRsp, err := s.userAPI.SearchTenants(ctx, &userv1.SearchTenantsRequest{
		OrganizationId: req.GetOrganizationId(),
		Status:         toSvcStatus(req.GetStatus()),
		Pagination:     toSvcPagination(req.GetPagination()),
		TreatmentId:    req.GetTreatmentId(),
		TenantName:     req.GetTenantName(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 换转
	results := make([]*pb.TenantEntity, len(searchRsp.GetTenants()))
	for k, v := range searchRsp.GetTenants() {
		results[k] = toApiTenant(v, s.s3Domain)
	}

	st := make(map[string]*pb.TenantSubscriptionTimeline)
	for _, each := range searchRsp.GetSubscriptions() {
		st[each.GetTenantId()] = &pb.TenantSubscriptionTimeline{
			TenantId:   each.GetTenantId(),
			StartTime:  each.GetStartTime(),
			EndTime:    each.GetEndTime(),
			TenantName: each.GetTenantName(),
			Status:     toTenantStatus(each.GetStatus()),
			Years:      each.GetYears(),
		}
	}

	rsp.Tenants = results
	rsp.Subscriptions = st
	rsp.Treatments = searchRsp.GetTreatments()
	rsp.TotalCount = searchRsp.GetTotalCount()

	return nil
}

// 验证request
func validateSearchTenantsRequest(req *pb.SearchTenantsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
