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

func (s *OrganizationAPIHandler) SearchReports(ctx context.Context, req *pb.SearchReportsRequest, rsp *pb.SearchReportsResponse) error {
	err := validateSearchReportsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 发送请求
	searchRsp, err := s.reportAPI.SearchReports(ctx, &reportv1.SearchReportsRequest{
		// 客户类型
		CustomerType: toReportCustomerType(req.GetCustomerType()),
		// 客户名称
		CustomerName: req.GetCustomerName(),
		// 客户手机号
		CustomerPhone: req.GetCustomerPhone(),
		// 开始时间
		StartTime: req.GetStartTime(),
		// 结束时间
		EndTime: req.GetEndTime(),
		// 脏腑辨证
		DirtyDialectics: req.GetDirtyDialectics(),
		// 商户名称
		TenantName: req.GetTenantName(),
		// 组织 ID
		OrganizationId: req.GetOrganizationId(),
		// 分页信息
		Pagination: toReportPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	rsp.TotalCount = searchRsp.GetTotalCount()
	sps := make([]*pb.SummaryReport, len(searchRsp.GetReports()))
	for k, v := range searchRsp.GetReports() {
		sps[k] = toAppSummaryReport(v)
	}
	rsp.Reports = sps

	return nil
}

// 验证request
func validateSearchReportsRequest(req *pb.SearchReportsRequest) error {
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	if req.GetStartTime().AsTime().After(req.GetEndTime().AsTime()) {
		return gerr.New("invalid time")
	}
	return nil
}
