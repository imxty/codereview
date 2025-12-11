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

func (s *ReportAPIHandler) GetOrganizationReportCount(ctx context.Context, req *pb.GetOrganizationReportCountRequest, rsp *pb.GetOrganizationReportCountResponse) error {
	err := validateGetOrganizationReportCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 查询组织下的常客/散客累计测量次数
	// 传入当日的时间即可
	now := time.Now().UTC()
	startTime := ptime.DayBeginUTCTime(now, ptime.LocBeijing)
	addCM, addTm, totalCM, totalTM, err := s.reportStore.GetOrganizationReportCount(ctx, req.GetOrganizationId(), startTime.UTC(), now)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 当日新增常客测量次数
	rsp.TodayCustomerMeasurementCount = int32(addCM)
	// 当日新增散客测量次数
	rsp.TodayTempCustomerMeasurementCount = int32(addTm)
	// 常客今年累计测量次数
	rsp.CustomerMeasurementYearCount = int32(totalCM)
	// 散客今年累计测量次数
	rsp.TempCustomerMeasurementYearCount = int32(totalTM)
	return nil
}

// 验证request
func validateGetOrganizationReportCountRequest(req *pb.GetOrganizationReportCountRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	return nil
}
