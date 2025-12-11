package report

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 批量获取组织商户每月报告统计请求
func (s *ReportAPIHandler) ListOrganizationTenantReportDownload(ctx context.Context, req *pb.ListOrganizationTenantReportDownloadRequest, rsp *pb.ListOrganizationTenantReportDownloadResponse) error {
	// 1.验证request
	err := validateListOrganizationTenantReportDownloadRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询商户id
	searchRsp, err := s.userAPI.SearchTenantByNameAndOrganizationID(ctx, &userv1.SearchTenantByNameAndOrganizationIDRequest{
		OrganizationId: req.GetOrganizationId(),
		TenantName:     req.GetTenantName(),
	})
	if err != nil {
		return nil
	}
	// 如果没有商户直接返回空即可
	tids := searchRsp.GetTenantIds()
	if len(tids) == 0 {
		return nil
	}

	// 查询商户name
	getRsp, err := s.userAPI.BatchGetTenantNamesByIDs(ctx, &userv1.BatchGetTenantNamesByIDsRequest{
		TenantIds: tids,
	})
	if err != nil {
		return nil
	}
	tenantsName := getRsp.GetTenantNames()

	// 查询
	ts, _, err := s.reportStore.SearchTenantReportsDownload(ctx, req.GetOrganizationId(), tids, req.GetStartTime().AsTime().UTC(),
		req.GetEndTime().AsTime().UTC())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	tpcs := make([]*pb.TenantReportCount, len(ts))
	for k, v := range ts {
		tpcs[k] = &pb.TenantReportCount{
			TenantName:                      tenantsName[v.GetTenantID()],
			MonthlyCustomerMeasurementCount: v.GetMonthlyCount(),
		}
	}

	rsp.Counts = tpcs

	return nil
}

// 验证request
func validateListOrganizationTenantReportDownloadRequest(req *pb.ListOrganizationTenantReportDownloadRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be empty")
	}
	if req.GetEndTime() == nil {
		return gerr.New("end_time should not be empty")
	}
	return nil
}
