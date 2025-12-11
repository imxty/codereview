package user

import (
	"context"
	gerr "errors"

	reportv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取实体测量概况响应
func (u *UserAPIHandler) GetEntityMeasurementSummary(ctx context.Context, req *pb.GetEntityMeasurementSummaryRequest, rsp *pb.GetEntityMeasurementSummaryResponse) error {
	err := validateGetEntityMeasurementSummary(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询商户在职员工信息
	staffs, err := u.userStore.ListActivatedStaffs(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	staff_ids := make([]string, len(staffs))
	for i, each := range staffs {
		staff_ids[i] = each.GetUserID()
	}

	getRsp, err := u.reportAPI.GetTenantReportCount(ctx, &reportv1.GetTenantReportCountRequest{
		TenantId: req.GetTenantId(),
		StaffIds: staff_ids,
	})
	if err != nil {
		return err
	}

	// 返回响应
	rsp.CustomerMeasurementTotalCount = getRsp.GetCustomerMeasurementYearCount()
	rsp.TempCustomerMeasurementTotalCount = getRsp.GetTempCustomerMeasurementYearCount()
	rsp.TodayCustomerMeasurementCount = getRsp.GetTodayCustomerMeasurementCount()
	rsp.TodayTempCustomerMeasurementCount = getRsp.GetTodayTempCustomerMeasurementCount()

	return nil
}

// 验证 request
func validateGetEntityMeasurementSummary(req *pb.GetEntityMeasurementSummaryRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
