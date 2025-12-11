package report

import (
	"context"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 批量获取组织商户每月报告统计请求
func (s *ReportAPIHandler) ListOrganizationTenantMonthlyReportCount(ctx context.Context, req *pb.ListOrganizationTenantMonthlyReportCountRequest, rsp *pb.ListOrganizationTenantMonthlyReportCountResponse) error {

	// 1.验证request
	err := validateListOrganizationTenantMonthlyReportCountRequest(req)
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
	ts, count, err := s.reportStore.SearchTenantReports(ctx, req.GetOrganizationId(), tids, req.GetStartTime().AsTime().UTC(),
		req.GetEndTime().AsTime().UTC(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回数据
	tpcs := make([]*pb.TenantReportCount, len(ts))
	for k, v := range ts {
		// 解析字符串为 time.Time，默认将日期设为1号
		layout := "2006-01" // Go 的时间格式必须按照这个 layout (年-月)
		parsedTime, err := time.Parse(layout, v.GetDate())
		if err != nil {
			return errors.Errorf(codes.InvalidOperation, "parse time error[%s]", err.Error())
		}
		tpcs[k] = &pb.TenantReportCount{
			// 商户名称
			TenantName: tenantsName[v.GetTenantID()],
			// 当月测量次数
			MonthlyCustomerMeasurementCount: v.GetMonthlyCount(),
			// 同比
			YearOnYear: v.GetYearOnYear(),
			// 环比
			MonthOnMonth: v.GetMonthOnMonth(),
			// 时间
			Time: timestamppb.New(parsedTime),
		}
	}

	rsp.Counts = tpcs
	rsp.TotalCount = int32(count)
	return nil
}

// 验证request
func validateListOrganizationTenantMonthlyReportCountRequest(req *pb.ListOrganizationTenantMonthlyReportCountRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
