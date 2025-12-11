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

func (s *CustomerAPIHandler) SearchAppCustomerLastStatus(ctx context.Context, req *pb.SearchAppCustomerLastStatusRequest, rsp *pb.SearchAppCustomerLastStatusResponse) error {
	err := validateSearchAppCustomerLastStatus(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商户信息
	getRsp, err := s.userAPI.GetTenantEntity(ctx, &userv1.GetTenantEntityRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return err
	}

	// 获取待复查常客信息
	lastCS, count, err := s.customerStore.BatchGetAppCustomersLastStatus(ctx, req.GetTenantId(), req.GetKey(), getRsp.GetOverdueDays(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	cIDs := make([]string, len(lastCS))
	for k, v := range lastCS {
		cIDs[k] = v.GetCustomerID()
	}

	cs, err := s.customerStore.BatchGetCustomers(ctx, cIDs, false)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	cMap := make(map[string]domain.CustomerIntf)
	for _, each := range cs {
		cMap[each.GetCustomerID()] = each
	}

	results := make([]*pb.CustomerLastStatus, len(lastCS))
	for k, v := range lastCS {
		updateTime := v.GetUpdatedAt()
		updateTime = time.Date(updateTime.Year(), updateTime.Month(), updateTime.Day(), 0, 0, 0, 0, time.Local)
		nowTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
		cInfo := cMap[v.GetCustomerID()]
		results[k] = &pb.CustomerLastStatus{
			TenantName:    getRsp.GetEntity().GetEntityName(),
			CustomerName:  cInfo.GetNickname(),
			CustomerPhone: cInfo.GetPhone(),
			OverdueCount:  int32(nowTime.Sub(updateTime).Hours() / 24),
			Gender:        toProtoGender(cInfo.GetGender()),
			Initial:       cInfo.GetInitial(),
		}
	}

	rsp.Customers = results
	rsp.TotalCount = int32(count)

	return nil
}

func validateSearchAppCustomerLastStatus(req *pb.SearchAppCustomerLastStatusRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
