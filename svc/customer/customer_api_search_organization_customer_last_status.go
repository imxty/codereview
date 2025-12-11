package customer

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/customer/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	MaxSize = 2147483647
)

// 获取组织客户最新测量状态
func (c *CustomerAPIHandler) SearchOrganizationCustomerLastStatus(ctx context.Context, req *pb.SearchOrganizationCustomerLastStatusRequest, rsp *pb.SearchOrganizationCustomerLastStatusResponse) error {
	err := validateOrganizationCustomerLastStatusRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 搜索商户
	searchRsp, err := c.userAPI.SearchTenantsByName(ctx, &userv1.SearchTenantsByNameRequest{
		TenantName: req.GetTenantName(),
	})
	if err != nil {
		return err
	}

	tids := searchRsp.GetTenantIds()

	// 获取商户名
	nameRsp, err := c.userAPI.BatchGetTenantNamesByIDs(ctx, &userv1.BatchGetTenantNamesByIDsRequest{
		TenantIds: tids,
	})
	if err != nil {
		return err
	}

	names_map := nameRsp.GetTenantNames()

	// 获取最新测量状态
	tcs, total_count, err := c.customerStore.BatchGetCustomersLastStatus(ctx, tids, req.GetCustomerName(), req.GetCustomerPhone(), 0, req.Pagination.GetOffset(), req.Pagination.GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获取客户信息
	cids := make([]string, len(tcs))
	for k, v := range tcs {
		cids[k] = v.GetCustomerID()
	}
	cs, err := c.customerStore.BatchGetCustomers(ctx, cids, false)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	cs_map := make(map[string]domain.CustomerIntf)
	for _, each := range cs {
		cs_map[each.GetCustomerID()] = each
	}

	results := make([]*pb.CustomerLastStatus, len(tcs))
	for k, v := range tcs {
		results[k] = toProtoCustomerLastStatus(v)
		results[k].CustomerName = cs_map[v.GetCustomerID()].GetNickname()
		results[k].CustomerPhone = cs_map[v.GetCustomerID()].GetPhone()
		results[k].TenantName = names_map[v.GetTenantID()]
	}

	rsp.Customers = results
	rsp.TotalCount = int32(total_count)

	return nil
}

func validateOrganizationCustomerLastStatusRequest(req *pb.SearchOrganizationCustomerLastStatusRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
