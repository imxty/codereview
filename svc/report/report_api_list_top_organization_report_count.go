package report

import (
	"context"
	gerr "errors"
	"sort"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	tt "github.com/jinmukeji/huimaibao-service/pkg/time"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

type reportsRank []domain.ReportRankIntf

// 获取组织商户测量top10
func (s *ReportAPIHandler) ListTopOrganizationReportCount(ctx context.Context, req *pb.ListTopOrganizationReportCountRequest, rsp *pb.ListTopOrganizationReportCountResponse) error {
	err := validateListTopOrganizationReportCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取组织下的所有商户
	getRsp, err := s.userAPI.GetOrganizationTenants(ctx, &userpb.GetOrganizationTenantsRequest{
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return err
	}
	tenants := getRsp.GetTenants()
	tenantIds := make([]string, len(tenants))
	tenantName := make(map[string]string)

	for k, v := range tenants {
		tenantIds[k] = v.GetTenantId()
		tenantName[v.GetTenantId()] = v.GetName()
	}

	// 如果没有商户直接返回即可
	if len(tenantIds) == 0 {
		return nil
	}

	// 获取所有商户的报告数量统计
	// 只获取本月的商品统计
	now := time.Now().In(tt.LocBeijing)
	// 获取本月开始时间和结束时间
	// [本月的开始，和下个月的第一天的0点)
	startTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, tt.LocBeijing)
	endTime := tt.GetLastDateOfMonth(now).AddDate(0, 0, 1)
	reportRanks, err := s.reportStore.ListTenantsReportCount(ctx, tenantIds, startTime, endTime)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	tenantReportRank := make(map[string]int64)
	for _, v := range reportRanks {
		tenantReportRank[v.GetTenantID()] = v.GetReportCount()
	}

	// 对商户报告统计排序
	sort.Sort(reportsRank(reportRanks))

	// 获取上个月的报告统计
	// [上月的开始，本月的第一天的0点)
	lastStartTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, tt.LocBeijing).AddDate(0, -1, 0)
	lastEndTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, tt.LocBeijing)
	lastReportRanks, err := s.reportStore.ListTenantsReportCount(ctx, tenantIds, lastStartTime, lastEndTime)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	lastTenantReportRank := make(map[string]int64)
	for _, v := range lastReportRanks {
		lastTenantReportRank[v.GetTenantID()] = v.GetReportCount()
	}

	length := 10
	if len(reportRanks) < 10 {
		length = len(reportRanks)
	}
	// 取较多的一个排名，不能超过10
	if len(lastReportRanks) > length {
		length = len(lastReportRanks)
		if length > 10 {
			length = 10
		}
	}
	protoCustomerRank := make([]*pb.TenantRank, length)
	// 返回数据
	if len(reportRanks) >= len(lastReportRanks) {
		// 返回前10的商户数据
		for i := 0; i < length; i++ {
			protoCustomerRank[i] = &pb.TenantRank{
				// 商户ID
				TenantId: reportRanks[i].GetTenantID(),
				// 商户名称
				StoreName: tenantName[reportRanks[i].GetTenantID()],
				// 上月次数
				LastMonthAmounts: int32(lastTenantReportRank[reportRanks[i].GetTenantID()]),
				// 本月次数
				Amounts: int32(reportRanks[i].GetReportCount()),
			}
		}
	} else {
		// 返回前10的商户数据
		for i := 0; i < length; i++ {
			protoCustomerRank[i] = &pb.TenantRank{
				// 商户ID
				TenantId: lastReportRanks[i].GetTenantID(),
				// 商户名称
				StoreName: tenantName[lastReportRanks[i].GetTenantID()],
				// 上月次数
				LastMonthAmounts: int32(lastReportRanks[i].GetReportCount()),
				// 本月次数
				Amounts: int32(tenantReportRank[lastReportRanks[i].GetTenantID()]),
			}
		}
	}

	rsp.MeasurementRank = protoCustomerRank
	return nil
}

// 验证request
func validateListTopOrganizationReportCountRequest(req *pb.ListTopOrganizationReportCountRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	return nil
}

// Len 获取数据集合元素个数
func (w reportsRank) Len() int {
	return len(w)
}

// Swap 交换 i 和 j 索引的两个元素的位置
func (w reportsRank) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

// Less 如果 i 索引的数据小于 j 索引的数据，返回 true，且不会调用下面的 Swap()，即数据升序排序。
func (w reportsRank) Less(i, j int) bool {
	return w[i].GetReportCount() > w[j].GetReportCount()
}
