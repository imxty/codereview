package customer

import (
	"context"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *CustomerAPIHandler) SearchOrganizationTenantCompare(ctx context.Context, req *pb.SearchOrganizationTenantCompareRequest, rsp *pb.SearchOrganizationTenantCompareResponse) error {
	err := validateSearchOrganizationTenantCompareRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询商户信息
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

	tenant_names_map := getRsp.GetTenantNames()

	// 查询商户新增客户
	currentRanks, count, err := s.customerStore.BatchGetTenantMonthlyAddedCustomerCount(ctx, searchRsp.GetTenantIds(), req.GetStartTime().AsTime().UTC(), req.GetEndTime().AsTime().UTC(), req.GetPagination().GetOffset(), req.GetPagination().GetSize())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 整合数据
	results := make([]*pb.TenantCustomerCompare, len(currentRanks))
	for k, v := range currentRanks {
		parsedTime, _ := time.Parse("2006-01", v.Month)
		results[k] = &pb.TenantCustomerCompare{
			TenantName:   tenant_names_map[v.TenantID],
			MonthCount:   int32(v.AddedCount),
			MonthOnMonth: v.MonthOnMonth,
			YearOnYear:   v.YearOnYear,
			Date:         timestamppb.New(parsedTime),
		}
	}

	rsp.Tenants = results
	rsp.TotalCount = int32(count)

	return nil
}

func validateSearchOrganizationTenantCompareRequest(req *pb.SearchOrganizationTenantCompareRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be empty")
	}
	if req.GetEndTime() == nil {
		return gerr.New("end_time should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
