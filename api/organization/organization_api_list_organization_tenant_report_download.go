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

// 获取组织测量导出数据
func (s *OrganizationAPIHandler) ListOrganizationTenantReportDownload(ctx context.Context, req *pb.ListOrganizationTenantReportDownloadRequest, rsp *pb.ListOrganizationTenantReportDownloadResponse) error {
	err := validateListOrganizationTenantReportDownloadRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.reportAPI.ListOrganizationTenantReportDownload(ctx, &reportv1.ListOrganizationTenantReportDownloadRequest{
		OrganizationId: req.GetOrganizationId(),
		TenantName:     req.GetTenantName(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	tcs := make([]*pb.TenantReportCount, len(getRsp.GetCounts()))
	for k, v := range getRsp.GetCounts() {
		tcs[k] = toAppTenantReportCount(v)
	}

	rsp.Counts = tcs

	return nil
}

// 验证request
func validateListOrganizationTenantReportDownloadRequest(req *pb.ListOrganizationTenantReportDownloadRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be empty")
	}
	if req.GetEndTime() == nil {
		return gerr.New("start_time should not be empty")
	}
	return nil
}
