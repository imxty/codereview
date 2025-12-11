package report

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// ListReports 查看常客/散客的报告
func (s *ReportAPIHandler) ListReports(ctx context.Context, req *pb.ListReportsRequest, rsp *pb.ListReportsResponse) error {

	// 1.验证request
	err := validateListReportsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取报告共享开关状态
	reportSharingRsp, err := s.userAPI.GetSharingReportSwitchStatus(ctx, &userpb.GetSharingReportSwitchStatusRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 如果开关打开但是没有传员工ID则返回错误
	if !reportSharingRsp.GetSharingReportSwitchStatus() && req.GetStaffId() == "" {
		return errors.Error(codes.InvalidRequest, "The report is in a shared state but there is no staff ID")
	}

	// 数据格式转化
	startTime := req.GetStartTime().AsTime()

	endTime := req.GetEndTime().AsTime()
	// 如果是常客就带着customer_id查询，如果不是常客，那么customer_id就是空
	reports, count, err := s.reportStore.ListReportsWithStaff(ctx, req.GetTenantId(), req.GetCustomerId(), req.GetStaffId(),
		reportSharingRsp.GetSharingReportSwitchStatus(), startTime, endTime, req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 如果没有报告直接返回空即可
	if len(reports) == 0 {
		return nil
	}

	// 类型转化
	// 获取所有员工的信息
	bookStaff := make(map[string]bool)
	staffIds := []string{}
	for _, v := range reports {
		if !bookStaff[v.GetStaffID()] {
			staffIds = append(staffIds, v.GetStaffID())
			bookStaff[v.GetStaffID()] = true
		}
	}
	getStaffsRsp, err := s.userAPI.BatchGetStaffs(ctx, &userpb.BatchGetStaffsRequest{
		StaffIds: staffIds,
	})
	if err != nil {
		return err
	}
	// 构建员工的信息map
	staffInfoMap := getStaffsRsp.GetStaffs()
	appReports := make([]*pb.SimpleHealthReport, len(reports))
	for k, v := range reports {
		appReports[k], err = toProtoSimpleReport(v)
		if err != nil {
			return errors.Errorf(codes.InvalidOperation, "failed to get simple report[%s]", err.Error())
		}
		// 查到对应的员工名称
		appReports[k].StaffName = staffInfoMap[v.GetStaffID()].GetName()
	}
	rsp.Reports = appReports
	rsp.TotalCount = int32(count)

	return nil
}

// 验证request
func validateListReportsRequest(req *pb.ListReportsRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be empty")
	}
	if req.GetEndTime() == nil {
		return gerr.New("end_time should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}
