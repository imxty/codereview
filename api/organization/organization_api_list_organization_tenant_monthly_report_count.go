package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	reportv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取组织测量概况
func (s *OrganizationAPIHandler) ListOrganizationTenantMonthlyReportCount(ctx context.Context, req *pb.ListOrganizationTenantMonthlyReportCountRequest, rsp *pb.ListOrganizationTenantMonthlyReportCountResponse) error {
	err := validateListOrganizationTenantMonthlyReportCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取测量概况
	getRsp, err := s.reportAPI.ListOrganizationTenantMonthlyReportCount(ctx, &reportv1.ListOrganizationTenantMonthlyReportCountRequest{
		OrganizationId: req.GetOrganizationId(),
		// 商户名称
		TenantName: req.GetTenantName(),
		// 开始时间
		StartTime: req.GetStartTime(),
		// 结束时间
		EndTime: req.GetEndTime(),
		// 分页
		Pagination: toReportPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	tcs := make([]*pb.TenantReportCount, len(getRsp.GetCounts()))
	for k, v := range getRsp.GetCounts() {
		tcs[k] = toAppTenantReportCount(v)
	}
	rsp.Counts = tcs
	rsp.TotalCount = getRsp.GetTotalCount()

	return nil
}

// 验证request
func validateListOrganizationTenantMonthlyReportCountRequest(req *pb.ListOrganizationTenantMonthlyReportCountRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
