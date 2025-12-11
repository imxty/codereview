package customer

import (
	"context"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	ptime "github.com/jinmukeji/huimaibao-service/pkg/time"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *CustomerAPIHandler) GetCustomerCount(ctx context.Context, req *pb.GetCustomerCountRequest, rsp *pb.GetCustomerCountResponse) error {
	// 验证request
	err := validateGetCustomerCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	now := time.Now().UTC()
	startTime := ptime.DayBeginUTCTime(now, ptime.LocBeijing)
	// 查询统计数据
	totalCustomer, addCustomer, err := s.customerStore.GetCustomerStatistics(ctx, req.GetTenantId(), startTime, now.UTC())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 员工排行map
	staffRankMap := make(map[string]*pb.StaffCustomerCount)

	if len(req.GetStaffIds()) > 0 {
		staffTodayRank, err := s.customerStore.GetTodayStaffCustomerCount(ctx, req.GetTenantId(), req.GetStaffIds(), startTime, now.UTC())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		staffTotalRank, err := s.customerStore.GetTotalStaffCustomerCount(ctx, req.GetTenantId(), req.GetStaffIds())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}

		for k, v := range staffTotalRank {
			staffRankMap[k] = &pb.StaffCustomerCount{
				TotalCustomer: v,
				TodayCustomer: staffTodayRank[k],
			}
		}
	}

	rsp.CustomerTotalCount = int32(totalCustomer)
	rsp.TodayAddedCustomerCount = int32(addCustomer)
	rsp.StaffCustomerCount = staffRankMap

	return nil
}

// 验证request
func validateGetCustomerCountRequest(req *pb.GetCustomerCountRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
