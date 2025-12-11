package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取商户排行
func (s *OrganizationAPIHandler) ListTopTenantRank(ctx context.Context, req *pb.ListTopTenantRankRequest, rsp *pb.ListTopTenantRankResponse) error {
	err := validateListTopTenantRankRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetOrganizationRankSummary(ctx, &userv1.GetOrganizationRankSummaryRequest{
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	results := make([]*pb.TenantRank, len(getRsp.GetTenantRank()))
	for k, v := range getRsp.GetTenantRank() {
		results[k] = &pb.TenantRank{
			TenantId:                    v.GetTenantId(),
			StoreName:                   v.GetStoreName(),
			LastMonthCustomerAmounts:    v.GetLastMonthCustomerAmounts(),
			CustomerAmounts:             v.GetCustomerAmounts(),
			LastMonthMeasurementAmounts: v.GetLastMonthMeasurementAmounts(),
			MeasurementAmounts:          v.GetMeasurementAmounts(),
		}
	}

	rsp.RankCount = results

	return nil
}

// 验证request
func validateListTopTenantRankRequest(req *pb.ListTopTenantRankRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
