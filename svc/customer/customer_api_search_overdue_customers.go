package customer

import (
	"context"
	gerr "errors"
	"time"

	"github.com/jinmukeji/go-pkg/v2/age"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/customer/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *CustomerAPIHandler) SearchOverdueCustomers(ctx context.Context, req *pb.SearchOverdueCustomersRequest, rsp *pb.SearchOverdueCustomersResponse) error {
	err := validateSearchOverdueCustomers(req)
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

	tenant := getRsp.GetEntity()

	// 搜索待复查常客信息
	cs, count, err := s.customerStore.SearchTenantOverdueCustomers(ctx, tenant.GetTenantId(), req.GetKey(), tenant.GetOverdue(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 查询员工信息
	cids := make([]string, len(cs))
	for k, v := range cs {
		cids[k] = v.GetStaffID()
	}
	if len(cids) < 1 {
		rsp.Customers = make([]*pb.OverdueCustomer, 0)
		rsp.TotalCount = 0
		return nil
	}
	batchRsp, err := s.userAPI.BatchGetStaffs(ctx, &userv1.BatchGetStaffsRequest{
		StaffIds: cids,
	})
	if err != nil {
		return err
	}

	staffs := batchRsp.GetStaffs()
	staff_names := make(map[string]string)
	for _, each := range staffs {
		staff_names[each.StaffId] = each.Name
	}

	csIDs := make([]string, len(cs))
	for k, v := range cs {
		csIDs[k] = v.GetCustomerID()
	}

	// 搜索最近的测量记录
	status, err := s.customerStore.ListCustomerLastStatus(ctx, csIDs)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	status_map := make(map[string]domain.CustomerLastStatusIntf)
	for _, each := range status {
		status_map[each.GetCustomerID()] = each
	}

	// 数据转换
	results := make([]*pb.OverdueCustomer, len(cs))
	for k, v := range cs {
		updateTime := status_map[v.GetCustomerID()].GetUpdatedAt()
		updateTime = time.Date(updateTime.Year(), updateTime.Month(), updateTime.Day(), 0, 0, 0, 0, time.Local)
		nowTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
		results[k] = &pb.OverdueCustomer{
			Customer:             toProtoCustomer(v, getNicknameInitial(v.GetNickname()), int32(age.Age(v.GetBirthday()))),
			LastMeasurementTime:  timestamppb.New(status_map[v.GetCustomerID()].GetUpdatedAt()),
			SinceLastMeasurement: int32(nowTime.Sub(updateTime).Hours() / 24),
			StaffName:            staff_names[v.GetStaffID()],
		}
	}

	rsp.Customers = results
	rsp.TotalCount = int32(count)

	return nil
}

func validateSearchOverdueCustomers(req *pb.SearchOverdueCustomersRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
