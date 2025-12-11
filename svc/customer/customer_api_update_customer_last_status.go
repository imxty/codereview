package customer

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/svc/customer/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *CustomerAPIHandler) UpdateCustomerLastStatus(ctx context.Context, req *pb.UpdateCustomerLastStatusRequest, rsp *pb.UpdateCustomerLastStatusResponse) error {
	err := validateUpdateCustomerLastStatusRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取客户的最新测量数据
	status, err := s.customerStore.ListCustomerLastStatus(ctx, []string{req.GetCustomerId()})
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	customers, err := s.customerStore.BatchGetCustomers(ctx, []string{req.GetCustomerId()}, false)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 客户不存在
	if len(customers) < 1 {
		return errors.Errorf(ErrCustomerNotFound, "customer_id[%s] not found", req.GetCustomerId())
	}

	// 如果没有数据，说明是第一次测量，创建数据
	if len(status) < 1 {
		err := s.customerStore.CreateCustomerLastStatus(ctx, &domain.CustomerLastStatus{
			CustomerID: req.GetCustomerId(),
			ReportID:   req.GetReportId(),
			TenantID:   customers[0].GetTenantID(),
		})
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		return nil
	}

	// 有数据，更新客户状态
	err = s.customerStore.UpdateCustomerLastStatus(ctx, req.GetCustomerId(), req.GetReportId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

func validateUpdateCustomerLastStatusRequest(req *pb.UpdateCustomerLastStatusRequest) error {
	if req.GetCustomerId() == "" {
		return gerr.New("customer_id should not be empty")
	}
	if req.GetReportId() == "" {
		return gerr.New("report_id should not be empty")
	}
	return nil
}
