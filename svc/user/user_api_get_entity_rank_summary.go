package user

import (
	"context"
	gerr "errors"
	"sort"

	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	reportv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) GetEntityRankSummary(ctx context.Context, req *pb.GetEntityRankSummaryRequest, rsp *pb.GetEntityRankSummaryResponse) error {
	err := validateGetEntityRankSummary(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询商户在职员工信息
	staffs, err := u.userStore.ListActivatedStaffs(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	staff_names := make(map[string]string)
	staff_ids := make([]string, len(staffs))
	for i, each := range staffs {
		staff_ids[i] = each.GetUserID()
		staff_names[each.GetUserID()] = each.GetNickname()
	}

	// 查询常客添加情况
	customerRsp, err := u.customerAPI.GetCustomerCount(ctx, &customerv1.GetCustomerCountRequest{
		TenantId: req.GetTenantId(),
		StaffIds: staff_ids,
	})
	if err != nil {
		return err
	}

	customer_summary := customerRsp.GetStaffCustomerCount()

	// 获取测量情况
	measureRsp, err := u.reportAPI.GetTenantReportCount(ctx, &reportv1.GetTenantReportCountRequest{
		TenantId: req.GetTenantId(),
		StaffIds: staff_ids,
	})
	if err != nil {
		return err
	}

	measuremant_summary := measureRsp.GetStaffCustomerCount()

	// 整合数据
	results := make([]*pb.StaffRank, len(staffs))
	for k, v := range staffs {
		var mToday, mTotal, cToday, cTotal int32
		if mInfo, ok := measuremant_summary[v.GetUserID()]; ok {
			mToday = mInfo.GetTodayMeasurement()
			mTotal = mInfo.GetTotalMeasurement()
		}
		if cInfo, ok := customer_summary[v.GetUserID()]; ok {
			cToday = cInfo.GetTodayCustomer()
			cTotal = cInfo.GetTotalCustomer()
		}
		results[k] = &pb.StaffRank{
			Name:             v.GetNickname(),
			TodayMeasurement: mToday,
			TotalMeasurement: mTotal,
			TodayCustomer:    cToday,
			TotalCustomer:    cTotal,
		}
	}

	// 排序
	sort.Slice(results, func(i, j int) bool {
		if results[i].TotalCustomer == results[j].TotalCustomer {
			return results[i].TodayCustomer > results[j].TodayCustomer
		}
		return results[i].TotalCustomer > results[j].TotalCustomer
	})

	rsp.StaffRank = results

	return nil
}

// 验证 request
func validateGetEntityRankSummary(req *pb.GetEntityRankSummaryRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
