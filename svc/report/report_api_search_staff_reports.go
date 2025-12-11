package report

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SearchStaffReports 查询报告测量次数
func (s *ReportAPIHandler) SearchStaffReports(ctx context.Context, req *pb.SearchStaffReportsRequest, rsp *pb.SearchStaffReportsResponse) error {
	err := validateSearchStaffReportsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询员工ID
	searchRsp, err := s.userAPI.SearchStaffsByNameAndStatus(ctx, &userv1.SearchStaffsByNameAndStatusRequest{
		// 商户 ID
		TenantId: req.GetTenantId(),
		// 员工名
		StaffName: req.GetStaffName(),
		// 员工状态
		IsActivated: req.GetIsActivated(),
	})
	if err != nil {
		return err
	}
	// 获取员工ID
	staffs := searchRsp.GetStaffs()
	staffIds := make([]string, len(staffs))
	for k, v := range staffs {
		staffIds[k] = v.GetStaffId()
	}
	// 没有直接返回空即可
	if len(staffIds) == 0 {
		return nil
	}

	// 查询员工信息
	getRsp, err := s.userAPI.BatchGetStaffs(ctx, &userv1.BatchGetStaffsRequest{
		StaffIds: staffIds,
	})
	if err != nil {
		return err
	}
	staffInfo := getRsp.GetStaffs()

	// 查询
	ccs, counts, err := s.reportStore.SearchStaffReports(ctx, req.GetTenantId(), staffIds, req.GetStartTime().AsTime().UTC(), req.GetEndTime().AsTime().UTC(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.TotalCount = int32(counts)
	// 转化数据
	scs := make([]*pb.StaffReport, len(ccs))
	for k, v := range ccs {
		// 返回员工信息
		sid := v.GetStaffID()
		staffInfo := staffInfo[sid]
		scs[k] = &pb.StaffReport{
			Nickname:                staffInfo.GetName(),
			Phone:                   staffInfo.GetPhone(),
			IsActivated:             staffInfo.GetIsActivated(),
			CustomerReportCount:     int32(v.GetCustomerCount()),
			TempCustomerReportCount: int32(v.GetTempCount()),
			StartTime:               timestamppb.New(v.GetDate()),
		}
	}
	rsp.ReportCount = scs
	return nil
}

// 验证request
func validateSearchStaffReportsRequest(req *pb.SearchStaffReportsRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	st := req.GetStartTime().AsTime()
	et := req.GetEndTime().AsTime()
	if st.After(et) {
		return gerr.New("invalid  time")
	}
	return nil
}
