package user

import (
	"context"
	gerr "errors"

	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 查询组织商户数据对比请求
func (u *UserAPIHandler) SearchOrganizationTenantCompare(ctx context.Context, req *pb.SearchOrganizationTenantCompareRequest, rsp *pb.SearchOrganizationTenantCompareResponse) error {
	err := validateSearchOrganizationTenantCompareRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchRsp, err := u.customerAPI.SearchOrganizationTenantCompare(ctx, &customerv1.SearchOrganizationTenantCompareRequest{
		OrganizationId: req.GetOrganizationId(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
		TenantName:     req.GetTenantName(),
		Pagination:     toCustomerPagination(req.GetPagination()),
	})
	if err != nil {
		return err
	}

	results := make([]*pb.OrganizationTenantCompare, len(searchRsp.GetTenants()))
	for k, v := range searchRsp.GetTenants() {
		results[k] = &pb.OrganizationTenantCompare{
			Date:         v.GetDate(),
			TenantName:   v.GetTenantName(),
			MonthCount:   v.GetMonthCount(),
			YearOnYear:   v.GetYearOnYear(),
			MonthOnMonth: v.GetMonthOnMonth(),
		}
	}

	rsp.Tenants = results
	rsp.TotalCount = searchRsp.TotalCount

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
		return gerr.New("pagination should not be empty")
	}
	return nil
}
