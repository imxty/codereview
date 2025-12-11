package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) ListReports(ctx context.Context, req *pb.ListReportsRequest, rsp *pb.ListReportsResponse) error {
	err := validateListReportsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 发送请求
	listRsp, err := s.reportAPI.ListReports(ctx, &reportpb.ListReportsRequest{
		TenantId:   req.GetTenantId(),
		CustomerId: req.GetCustomerId(),
		StartTime:  req.GetStartTime(),
		EndTime:    req.GetEndTime(),
		Pagination: toSvcPagination(req.GetPagination()),
		StaffId:    req.GetStaffId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	// 返回数据
	pbReports := make([]*pb.HealthReport, len(listRsp.GetReports()))
	for k, v := range listRsp.GetReports() {
		pbReports[k] = toAppHealthReport(v)
	}
	rsp.Reports = pbReports
	rsp.TotalCount = listRsp.GetTotalCount()
	return nil
}

// 验证request
func validateListReportsRequest(req *pb.ListReportsRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be empty")
	}
	if req.GetEndTime() == nil {
		return gerr.New("end_time should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}
