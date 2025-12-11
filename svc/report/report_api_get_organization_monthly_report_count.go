package report

import (
	"context"
	gerr "errors"
	"math"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	ptime "github.com/jinmukeji/huimaibao-service/pkg/time"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 添加常客报告请求
func (s *ReportAPIHandler) GetOrganizationMonthlyReportCount(ctx context.Context, req *pb.GetOrganizationMonthlyReportCountRequest, rsp *pb.GetOrganizationMonthlyReportCountResponse) error {

	// 1.验证request
	err := validateGetOrganizationMonthlyReportCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	mTime := req.GetMeasurementTime()
	// 计算3个时间段的测量数据
	// 1. MeasurementTime当月测量次数
	// 2. 去年MeasurementTime测量次数
	// 3. 上个月MeasurementTime测量次数
	var measurementStartTime, measurementEndTime time.Time
	var lastYearMeasurementStartTime, lastYearMeasurementEndTime time.Time
	var lastMonthMeasurementStartTime, lastMonthMeasurementEndTime time.Time
	t := time.Date(int(mTime.GetYear()), time.Month(mTime.GetMonth()), 5, 0, 0, 0, 0, ptime.LocBeijing)
	measurementStartTime = ptime.GetFirstDateOfMonth(t)
	measurementEndTime = measurementStartTime.AddDate(0, 1, 0)
	lastYearMeasurementStartTime = measurementStartTime.AddDate(-1, 0, 0)
	lastYearMeasurementEndTime = measurementEndTime.AddDate(-1, 0, 0)
	lastMonthMeasurementStartTime = measurementStartTime.AddDate(0, -1, 0)
	lastMonthMeasurementEndTime = measurementEndTime.AddDate(0, -1, 0)

	// 分别查询报告数量
	monthCount, err := s.reportStore.GetOrganizationReportsCount(ctx, req.GetOrganizationId(), measurementStartTime, measurementEndTime)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	lastMonthCount, err := s.reportStore.GetOrganizationReportsCount(ctx, req.GetOrganizationId(), lastMonthMeasurementStartTime, lastMonthMeasurementEndTime)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	lastYearCount, err := s.reportStore.GetOrganizationReportsCount(ctx, req.GetOrganizationId(), lastYearMeasurementStartTime, lastYearMeasurementEndTime)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回数据
	rsp.MonthlyCustomerMeasurementCount = int32(monthCount)
	// 环比
	if lastMonthCount == 0 {
		rsp.MonthOnMonth = InvalidPercent
	} else {
		rsp.MonthOnMonth = math.Floor((float64(monthCount-lastMonthCount) / float64(lastMonthCount)) * 100)
	}
	// 同比
	if lastYearCount == 0 {
		rsp.YearOnYear = InvalidPercent
	} else {
		rsp.YearOnYear = math.Floor((float64(monthCount-lastYearCount) / float64(lastYearCount)) * 100)
	}

	return nil
}

// 验证request
func validateGetOrganizationMonthlyReportCountRequest(req *pb.GetOrganizationMonthlyReportCountRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetMeasurementTime() == nil {
		return gerr.New("time should not be nil")
	}
	if err := validateDateTime(req.GetMeasurementTime()); err != nil {
		return err
	}
	return nil
}
