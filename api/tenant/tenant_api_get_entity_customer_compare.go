package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) GetEntityCustomerCompare(ctx context.Context, req *pb.GetEntityCustomerCompareRequest, rsp *pb.GetEntityCustomerCompareResponse) error {
	err := validateGetEntityCustomerCompareRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetEntityCustomerCompareSummary(ctx, &userv1.GetEntityCustomerCompareSummaryRequest{
		TenantId: req.GetTenantId(),
		Date:     req.GetDate(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.CustomerAddedCount = getRsp.GetMonthAddedCustomerCount()
	rsp.YearOnYear = getRsp.GetYearOnYear()
	rsp.MonthOnMonth = getRsp.GetMonthOnMonth()

	return nil
}

// 验证request
func validateGetEntityCustomerCompareRequest(req *pb.GetEntityCustomerCompareRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetDate() == nil {
		return gerr.New("date should not be nil")
	}
	return nil
}
