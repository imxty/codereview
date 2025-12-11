package organization

import (
	"context"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
)

// 查询组织商户数据对比
func (s *OrganizationAPIHandler) SearchOrganizationTenantCompare(ctx context.Context, req *pb.SearchOrganizationTenantCompareRequest, rsp *pb.SearchOrganizationTenantCompareResponse) error {
	searchRsp, err := s.userAPI.SearchOrganizationTenantCompare(ctx, &userv1.SearchOrganizationTenantCompareRequest{
		OrganizationId: req.GetOrganizationId(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
		TenantName:     req.GetTenantName(),
		Pagination:     toSvcPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回响应
	results := make([]*pb.OrganizationTenantCompare, len(searchRsp.GetTenants()))
	for k, v := range searchRsp.GetTenants() {
		results[k] = toApiOrganizationTenantCompare(v)
	}
	rsp.Tenants = results
	rsp.TotalCount = searchRsp.GetTotalCount()

	return nil
}
