package user

import (
	"context"
	gerr "errors"

	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	reportv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取组织商户数据对比概况
func (u *UserAPIHandler) GetOrganizationRankSummary(ctx context.Context, req *pb.GetOrganizationRankSummaryRequest, rsp *pb.GetOrganizationRankSummaryResponse) error {
	err := validateGetOrganizationRankSummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取组织商户测量 top10
	listRsp, err := u.reportAPI.ListTopOrganizationReportCount(ctx, &reportv1.ListTopOrganizationReportCountRequest{
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return err
	}

	if len(listRsp.GetMeasurementRank()) < 1 {
		rsp.TenantRank = make([]*pb.TenantRank, 0)
		return nil
	}

	// 获取客户数据
	tenant_ids := make([]string, len(listRsp.GetMeasurementRank()))
	for k, v := range listRsp.GetMeasurementRank() {
		tenant_ids[k] = v.GetTenantId()
	}
	getRsp, err := u.customerAPI.GetTenantsCustomerRank(ctx, &customerv1.GetTenantsCustomerRankRequest{
		TenantIds: tenant_ids,
	})
	if err != nil {
		return err
	}

	customers := getRsp.GetTenants()

	// 整合数据
	results := make([]*pb.TenantRank, len(listRsp.GetMeasurementRank()))
	for i, each := range listRsp.GetMeasurementRank() {
		results[i] = &pb.TenantRank{
			TenantId:                    each.GetTenantId(),
			StoreName:                   each.GetStoreName(),
			LastMonthCustomerAmounts:    customers[each.GetTenantId()].GetLastMonthAmounts(),
			CustomerAmounts:             customers[each.GetTenantId()].GetAmounts(),
			LastMonthMeasurementAmounts: each.GetLastMonthAmounts(),
			MeasurementAmounts:          each.GetAmounts(),
		}
	}

	rsp.TenantRank = results

	return nil
}

// 验证 request
func validateGetOrganizationRankSummaryRequest(req *pb.GetOrganizationRankSummaryRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
