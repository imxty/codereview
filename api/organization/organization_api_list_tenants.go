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

// 获取商户列表
func (s *OrganizationAPIHandler) ListTenants(ctx context.Context, req *pb.ListTenantsRequest, rsp *pb.ListTenantsResponse) error {
	err := validateListTenantsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	listRsp, err := s.userAPI.ListTenants(ctx, &userv1.ListTenantsRequest{
		OrganizationId: req.GetOrganizationId(),
		Pagination:     toApiPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	pTenants := listRsp.GetTenants()
	t := make([]*pb.TenantEntity, len(pTenants))
	for k, v := range pTenants {
		t[k] = toApiTenant(v, s.s3Domain)
	}

	st := make(map[string]*pb.TenantSubscriptionTimeline)
	for _, each := range listRsp.GetSubscriptions() {
		st[each.GetTenantId()] = &pb.TenantSubscriptionTimeline{
			TenantId:   each.GetTenantId(),
			StartTime:  each.GetStartTime(),
			EndTime:    each.GetEndTime(),
			TenantName: each.GetTenantName(),
			Status:     toTenantStatus(each.GetStatus()),
			Years:      each.GetYears(),
		}
	}

	rsp.Tenants = t
	rsp.Treatments = listRsp.GetTreatments()
	rsp.Subscriptions = st
	rsp.TotalCount = listRsp.GetTotalCount()

	return nil
}

// 验证request
func validateListTenantsRequest(req *pb.ListTenantsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
