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

func (s *TenantAPIHandler) GetEntitySummary(ctx context.Context, req *pb.GetEntitySummaryRequest, rsp *pb.GetEntitySummaryResponse) error {
	err := validateGetEntitySummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取待复查人数
	getRsp, err := s.userAPI.GetEntityOverdueSummary(ctx, &userv1.GetEntityOverdueSummaryRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.CustomerOverdueCount = getRsp.GetOverdueCustomerTotalCount()

	return nil
}

// 验证request
func validateGetEntitySummaryRequest(req *pb.GetEntitySummaryRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
