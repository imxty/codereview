package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *BossAPIHandler) SearchOrganizationReports(ctx context.Context, req *pb.SearchOrganizationReportsRequest, rsp *pb.SearchOrganizationReportsResponse) error {
	err := validateSearchOrganizationReportsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	if req.GetOrganizationId() == "" {
		rsp.ReportCount = make([]*pb.ReportCount, 0)
		rsp.TotalCount = 0
		return nil
	}

	// 发送请求
	searchRsp, err := s.reportAPI.SearchOrganizationReports(ctx, &reportpb.SearchOrganizationReportsRequest{
		// 组织ID
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
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.TotalCount = searchRsp.GetTotalCount()
	rts := make([]*pb.ReportCount, len(searchRsp.GetReportCount()))
	for k, v := range searchRsp.GetReportCount() {
		rts[k] = toAppReportCount(v)
	}
	rsp.ReportCount = rts

	return nil
}

// 验证request
func validateSearchOrganizationReportsRequest(req *pb.SearchOrganizationReportsRequest) error {
	st := req.GetStartTime().AsTime()
	et := req.GetEndTime().AsTime()
	if st.After(et) {
		return gerr.New("invalid  time")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}
