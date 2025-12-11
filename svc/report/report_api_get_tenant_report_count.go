package report

import (
	"context"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	ptime "github.com/jinmukeji/huimaibao-service/pkg/time"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取报告统计
func (s *ReportAPIHandler) GetTenantReportCount(ctx context.Context, req *pb.GetTenantReportCountRequest, rsp *pb.GetTenantReportCountResponse) error {
	err := validateGetTenantReportCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 统计报告数据
	now := time.Now()
	startTime := ptime.DayBeginUTCTime(now, ptime.LocBeijing)
	// 获取报告统计(当日新增常客测量次数，当日新增散客测量次数，常客累计测量次数,散客累计测量次数)
	// TM = tempCustomerMeasurementCount CM = customerMeasrementCount
	addCM, addTm, totalCM, totalTM, err := s.reportStore.GetTenantReportCount(ctx, req.GetTenantId(), startTime.UTC(), now.UTC())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 员工排行榜
	staffRankMap := make(map[string]*pb.StaffReportCount)
	if len(req.GetStaffIds()) != 0 {
		staffRank, err := s.reportStore.GetStaffMeasurementCount(ctx, req.GetTenantId(), req.GetStaffIds(), startTime, now.UTC())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 构建map返回
		for _, v := range staffRank {
			staffRankMap[v.GetStaffID().String()] = &pb.StaffReportCount{
				// 今日测量次数
				TodayMeasurement: int32(v.GetTodayMeasurement()),
				// 累计测量次数
				TotalMeasurement: int32(v.GetTotalMeasurement()),
			}
		}
	}

	// 当日新增常客测量次数
	rsp.TodayCustomerMeasurementCount = int32(addCM)
	// 当日新增散客测量次数
	rsp.TodayTempCustomerMeasurementCount = int32(addTm)
	// 常客今年累计测量次数
	rsp.CustomerMeasurementYearCount = int32(totalCM)
	// 散客今年累计测量次数
	rsp.TempCustomerMeasurementYearCount = int32(totalTM)
	rsp.StaffCustomerCount = staffRankMap

	return nil
}

// 验证request
func validateGetTenantReportCountRequest(req *pb.GetTenantReportCountRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
