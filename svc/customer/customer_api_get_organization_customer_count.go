package customer

import (
	"context"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	ptime "github.com/jinmukeji/huimaibao-service/pkg/time"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *CustomerAPIHandler) GetOrganizationCustomerCount(ctx context.Context, req *pb.GetOrganizationCustomerCountRequest, rsp *pb.GetOrganizationCustomerCountResponse) error {
	err := validateGetOrganizationCustomerCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取组织下的所有商户
	getRsp, err := s.userAPI.GetOrganizationTenants(ctx, &userpb.GetOrganizationTenantsRequest{
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return err
	}
	tenants := getRsp.GetTenants()
	tenantIds := make([]string, len(tenants))

	for k, v := range tenants {
		tenantIds[k] = v.GetTenantId()
	}

	// 查询商户下常客数量
	customerCount, err := s.customerStore.ListCustomerTotalCount(ctx, tenantIds)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 查询当天常数量
	// 当天的0点到现在
	now := time.Now().UTC()
	startTime := ptime.DayBeginUTCTime(now, ptime.LocBeijing)
	todayCount, err := s.customerStore.ListTodayCustomerCount(ctx, tenantIds, startTime.UTC(), now)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.TodayAddedCustomerCount = todayCount
	rsp.CustomerTotalCount = customerCount
	return nil
}

// 验证request
func validateGetOrganizationCustomerCountRequest(req *pb.GetOrganizationCustomerCountRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	return nil
}
