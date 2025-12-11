package customer

import (
	"context"
	gerr "errors"
	"math"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	InvalidPercent float64 = 2147483647
)

func (s *CustomerAPIHandler) GetEntityCustomerCompareSummary(ctx context.Context, req *pb.GetEntityCustomerCompareSummaryRequest, rsp *pb.GetEntityCustomerCompareSummaryResponse) error {
	err := validateGetEntityCustomerCompareSummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	start_time := req.GetDate().AsTime()
	// 当前月份的第一天
	currentMonthStart := time.Date(start_time.Year(), start_time.Month(), 1, 0, 0, 0, 0, start_time.Location())
	// 下个月第一天
	nextMonthStart := currentMonthStart.AddDate(0, 1, 0)
	// 去年同期的第一天
	lastYearStart := currentMonthStart.AddDate(-1, 0, 0)
	// 上个月的第一天
	lastMonthStart := currentMonthStart.AddDate(0, -1, 0)

	// 获取当月新增
	currentMonthCount, err := s.customerStore.GetTenantAddedCustomerCount(ctx, req.GetTenantId(), currentMonthStart, nextMonthStart)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获取上个月新增常
	lastMonthCount, err := s.customerStore.GetTenantAddedCustomerCount(ctx, req.GetTenantId(), lastMonthStart, currentMonthStart)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 查询去年新增常客数
	lastYearCount, err := s.customerStore.GetTenantAddedCustomerCount(ctx, req.GetTenantId(), lastYearStart, lastYearStart.AddDate(0, 1, 0))
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.MonthAddedCustomerCount = int32(currentMonthCount)
	if lastMonthCount == 0 {
		rsp.MonthOnMonth = InvalidPercent
	} else {
		rsp.MonthOnMonth = math.Floor((float64(currentMonthCount-lastMonthCount) / float64(lastMonthCount)) * 100)
	}
	if lastYearCount == 0 {
		rsp.YearOnYear = InvalidPercent
	} else {
		rsp.YearOnYear = math.Floor((float64(currentMonthCount-lastYearCount) / float64(lastYearCount)) * 100)
	}

	return nil
}

// 验证 request
func validateGetEntityCustomerCompareSummaryRequest(req *pb.GetEntityCustomerCompareSummaryRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetDate() == nil {
		return gerr.New("date should not be nil")
	}
	return nil
}
