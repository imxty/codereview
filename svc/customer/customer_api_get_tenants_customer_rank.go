package customer

import (
	"context"
	gerr "errors"
	"sort"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *CustomerAPIHandler) GetTenantsCustomerRank(ctx context.Context, req *pb.GetTenantsCustomerRankRequest, rsp *pb.GetTenantsCustomerRankResponse) error {
	err := validateGetTenantsCustomerRank(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商户名
	getRsp, err := s.userAPI.BatchGetTenantNamesByIDs(ctx, &userv1.BatchGetTenantNamesByIDsRequest{
		TenantIds: req.GetTenantIds(),
	})
	if err != nil {
		return err
	}
	tenant_names := getRsp.GetTenantNames()

	currentMonthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	currentMonthEnd := currentMonthStart.AddDate(0, 1, 0)
	lastMonthStart := currentMonthStart.AddDate(0, -1, 0)

	// 获取本月新增客户
	currentMonthCount, err := s.customerStore.BatchGetTenantsAddedCustomerCount(ctx, req.GetTenantIds(), currentMonthStart, currentMonthEnd)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 获取上月新增客户
	lastMonthCount, err := s.customerStore.BatchGetTenantsAddedCustomerCount(ctx, req.GetTenantIds(), lastMonthStart, currentMonthStart)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 排序
	temp := make([]*pb.TenantRank, 0)
	for k, v := range tenant_names {
		temp = append(temp, &pb.TenantRank{
			TenantId:         k,
			StoreName:        v,
			LastMonthAmounts: int32(lastMonthCount[k]),
			Amounts:          int32(currentMonthCount[k]),
		})
	}
	sort.Slice(temp, func(i, j int) bool {
		return (temp[i].Amounts + temp[i].LastMonthAmounts) > (temp[j].Amounts + temp[j].LastMonthAmounts)
	})

	// 返回数据
	results := make(map[string]*pb.TenantRank)
	for _, each := range temp {
		results[each.GetTenantId()] = each
	}

	rsp.Tenants = results

	return nil
}

// 验证 request
func validateGetTenantsCustomerRank(req *pb.GetTenantsCustomerRankRequest) error {
	if req.GetTenantIds() == nil {
		return gerr.New("tenant_ids should not be nil")
	}
	return nil
}
