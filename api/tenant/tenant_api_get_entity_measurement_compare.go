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

// 获取实体测量比较数据
func (s *TenantAPIHandler) GetEntityMeasurementCompare(ctx context.Context, req *pb.GetEntityMeasurementCompareRequest, rsp *pb.GetEntityMeasurementCompareResponse) error {
	err := validateGetEntityMeasurementCompareRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetEntityMeasurementCompareSummary(ctx, &userv1.GetEntityMeasurementCompareSummaryRequest{
		TenantId: req.GetTenantId(),
		Date:     req.GetDate(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.MeasurementAddedCount = getRsp.GetMonthMeasurementCount()
	rsp.MonthOnMonth = getRsp.GetMonthOnMonth()
	rsp.YearOnYear = getRsp.GetYearOnYear()

	return nil
}

// 验证request
func validateGetEntityMeasurementCompareRequest(req *pb.GetEntityMeasurementCompareRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetDate() == nil {
		return gerr.New("date should not be empty")
	}
	return nil
}
