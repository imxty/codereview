package report

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// SearchStaffReportsCount 查询报告测量总次数
func (s *ReportAPIHandler) SearchStaffReportsCount(ctx context.Context, req *pb.SearchStaffReportsCountRequest, rsp *pb.SearchStaffReportsCountResponse) error {
	err := validateSearchStaffReportsCountRequest(req)
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

	// 查询
	ccs, tcs, err := s.reportStore.SearchStaffReportsCount(ctx, req.GetTenantId(), staffIds, req.GetStartTime().AsTime().UTC(), req.GetEndTime().AsTime().UTC())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.CustomerCount = int32(ccs)
	rsp.TempCustomerCount = int32(tcs)
	return nil
}

// 验证request
func validateSearchStaffReportsCountRequest(req *pb.SearchStaffReportsCountRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	st := req.GetStartTime().AsTime()
	et := req.GetEndTime().AsTime()
	if st.After(et) {
		return gerr.New("invalid  time")
	}
	return nil
}
