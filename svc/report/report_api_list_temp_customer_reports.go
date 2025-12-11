package report

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ReportAPIHandler) ListTempCustomerReports(ctx context.Context, req *pb.ListTempCustomerReportsRequest, rsp *pb.ListTempCustomerReportsResponse) error {

	// 1.验证request
	err := validateListTempCustomerReportsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取报告
	if req.GetReportId() != "" {
		// 查询报告
		report, err := s.reportStore.GetReportByID(ctx, req.GetTenantId(), req.GetReportId())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		if report == nil {
			return errors.Error(ErrReportNotFound, "report not found")
		}
		// 转化格式
		simpleReports, err := toProtoSimpleReport(report)
		if err != nil {
			return errors.Errorf(codes.InvalidOperation, "failed to get simple report[%s]", err.Error())
		}
		rsp.Reports = []*pb.SimpleHealthReport{simpleReports}
		rsp.TotalCount = 1
	} else {
		// 检测时间
		if req.GetTimeRange() == nil {
			// 如果没有时间就是查询全部
			reports, reportSize, err := s.reportStore.ListReportsWithoutTime(ctx, req.GetTenantId(), "", req.GetPagination().GetOffset(), req.GetPagination().GetSize())
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
			// 转化数据格式
			simpleReports := make([]*pb.SimpleHealthReport, len(reports))
			for k, report := range reports {
				simpleReports[k], err = toProtoSimpleReport(report)
				if err != nil {
					return errors.Errorf(codes.InvalidOperation, "failed to get simple report[%s]", err.Error())
				}
			}
			rsp.Reports = simpleReports
			rsp.TotalCount = int32(reportSize)
			return nil
		}
		timeRange := req.GetTimeRange()
		startTime := timeRange.GetStartTime().AsTime()
		endTime := timeRange.GetEndTime().AsTime()
		if startTime.After(endTime) {
			return errors.Error(codes.InvalidRequest, "end_time should not earlier than start_time")
		}
		// 查询报告
		reports, reportSize, err := s.reportStore.ListReports(ctx, req.GetTenantId(), "", startTime, endTime, req.GetPagination().GetOffset(), req.GetPagination().GetSize())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 转化数据格式
		simpleReports := make([]*pb.SimpleHealthReport, len(reports))
		for k, report := range reports {
			simpleReports[k], err = toProtoSimpleReport(report)
			if err != nil {
				return errors.Errorf(codes.InvalidOperation, "failed to get simple report[%s]", err.Error())
			}
		}
		rsp.Reports = simpleReports
		rsp.TotalCount = int32(reportSize)
	}

	return nil
}

// 验证request
func validateListTempCustomerReportsRequest(req *pb.ListTempCustomerReportsRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}
