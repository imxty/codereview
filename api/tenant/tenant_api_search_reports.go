package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	reportv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) SearchReports(ctx context.Context, req *pb.SearchReportsRequest, rsp *pb.SearchReportsResponse) error {
	err := validateSearchReportsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchRsp, err := s.reportAPI.SearchReports(ctx, &reportv1.SearchReportsRequest{
		// 报告ID
		ReportId: req.GetReportId(),
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
		// 分页信息
		Pagination: toReportPagination(req.GetPagination()),
		// TODO: FIX NAME
		TenantId: req.GetOperatorId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 数据转换
	results := make([]*pb.SummaryReport, len(searchRsp.GetReports()))
	for k, v := range searchRsp.GetReports() {
		results[k] = toApiReport(v)
	}

	// 返回数据
	rsp.Reports = results
	rsp.TotalCount = searchRsp.GetTotalCount()

	return nil
}

// 验证request
func validateSearchReportsRequest(req *pb.SearchReportsRequest) error {
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be empty")
	}
	if req.GetEndTime() == nil {
		return gerr.New("end_time should not be empty")
	}
	if req.GetOperatorId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
