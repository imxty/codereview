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

// 获取员工排行榜
func (s *TenantAPIHandler) GetStaffRank(ctx context.Context, req *pb.GetStaffRankRequest, rsp *pb.GetStaffRankResponse) error {
	err := validateGetStaffRankRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取员工
	getRsp, err := s.userAPI.GetEntityRankSummary(ctx, &userv1.GetEntityRankSummaryRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	results := make([]*pb.StaffRank, len(getRsp.GetStaffRank()))
	for k, v := range getRsp.GetStaffRank() {
		results[k] = toApiStaffRank(v)
	}
	rsp.StaffRank = results

	return nil
}

// 验证request
func validateGetStaffRankRequest(req *pb.GetStaffRankRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
