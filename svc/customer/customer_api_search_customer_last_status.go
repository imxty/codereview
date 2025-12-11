package customer

import (
	"context"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/customer/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *CustomerAPIHandler) SearchCustomerLastStatus(ctx context.Context, req *pb.SearchCustomerLastStatusRequest, rsp *pb.SearchCustomerLastStatusResponse) error {
	err := validateSearchCustomerLastStatusRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询商户 ID
	searchRsp, err := s.userAPI.SearchTenantByNameAndOrganizationID(ctx, &userv1.SearchTenantByNameAndOrganizationIDRequest{
		OrganizationId: req.GetOrganizationId(),
		TenantName:     req.GetTenantName(),
	})
	if err != nil {
		return err
	}

	// 查询商户名
	getRsp, err := s.userAPI.BatchGetTenantNamesByIDs(ctx, &userv1.BatchGetTenantNamesByIDsRequest{
		TenantIds: searchRsp.GetTenantIds(),
	})
	if err != nil {
		return err
	}

	tenant_names := getRsp.GetTenantNames()

	// 获取最新测量情况
	last_status, count, err := s.customerStore.BatchGetCustomersLastStatus(ctx, searchRsp.GetTenantIds(), req.GetCustomerName(), req.GetCustomerPhone(), req.GetOverdueMask(), req.GetPagination().GetOffset(), req.Pagination.GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	customer_ids := make([]string, len(last_status))
	for k, v := range last_status {
		customer_ids[k] = v.GetCustomerID()
	}

	// 获取客户信息
	infos, err := s.customerStore.BatchGetCustomers(ctx, customer_ids, false)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	info_map := make(map[string]domain.CustomerIntf)
	for _, each := range infos {
		info_map[each.GetCustomerID()] = each
	}

	// 数据转换
	results := make([]*pb.CustomerLastStatus, len(last_status))
	for k, v := range last_status {
		updateTime := v.GetUpdatedAt()
		updateTime = time.Date(updateTime.Year(), updateTime.Month(), updateTime.Day(), 0, 0, 0, 0, time.Local)
		nowTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
		results[k] = &pb.CustomerLastStatus{
			TenantName:    tenant_names[v.GetTenantID()],
			CustomerName:  info_map[v.GetCustomerID()].GetNickname(),
			CustomerPhone: info_map[v.GetCustomerID()].GetPhone(),
			OverdueCount:  int32(nowTime.Sub(updateTime).Hours() / 24),
			Gender:        toProtoGender(info_map[v.GetCustomerID()].GetGender()),
		}
	}

	rsp.Customers = results
	rsp.TotalCount = int32(count)

	return nil
}

func validateSearchCustomerLastStatusRequest(req *pb.SearchCustomerLastStatusRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
