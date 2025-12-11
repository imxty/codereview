package customer

import (
	"context"
	gerr "errors"
	"time"

	"github.com/jinmukeji/go-pkg/v2/age"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *CustomerAPIHandler) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest, rsp *pb.GetCustomerResponse) error {
	// 验证request
	err := validateGetCustomerRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取domain customer
	dCustomer, err := s.customerStore.GetCustomer(ctx, req.GetTenantId(), req.GetCustomerId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if dCustomer == nil {
		return errors.Errorf(ErrCustomerNotFound, "customer[%s] not found", req.GetCustomerId())
	}

	var overdue_count int32

	// 获取常客上次测量
	status, err := s.customerStore.ListCustomerLastStatus(ctx, []string{dCustomer.GetCustomerID()})
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 计算距上次测量天数
	if len(status) < 1 {
		overdue_count = 0
	} else {
		updateTime := status[0].GetUpdatedAt()
		updateTime = time.Date(updateTime.Year(), updateTime.Month(), updateTime.Day(), 0, 0, 0, 0, time.Local)
		nowTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
		overdue_count = int32(nowTime.Sub(updateTime).Hours() / 24)
	}

	// 转换
	appCustomer := toProtoCustomer(dCustomer, getNicknameInitial(dCustomer.GetNickname()), int32(age.Age((dCustomer.GetBirthday()))))

	// 返回结果
	rsp.Customer = appCustomer
	rsp.OverdueCount = overdue_count
	return nil
}

// 验证request
func validateGetCustomerRequest(req *pb.GetCustomerRequest) error {
	if req.GetCustomerId() == "" {
		return gerr.New("customer id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
