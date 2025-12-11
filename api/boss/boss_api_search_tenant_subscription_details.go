package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 查询商户续期详情
func (s *BossAPIHandler) SearchTenantSubscriptionDetails(ctx context.Context, req *pb.SearchTenantSubscriptionDetailsRequest, rsp *pb.SearchTenantSubscriptionDetailsResponse) error {
	err := validateSearchTenantSubscriptionDetailsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询商户续期详情
	searchRsp, err := s.userAPI.SearchTenantSubscriptionDetails(ctx, &userv1.SearchTenantSubscriptionDetailsRequest{
		OrganizationId: req.GetOrganizationId(),
		TenantId:       req.GetTenantId(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
		ContactName:    req.GetContactName(),
		ContactPhone:   req.GetContactPhone(),
		Pagination:     toUserPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 数据转换
	results := make([]*pb.SubscriptionPeriod, len(searchRsp.GetSubscriptions()))
	for k, v := range searchRsp.GetSubscriptions() {
		results[k] = toApiSubscriptionPeriod(v)
	}

	rsp.Subscriptions = results
	rsp.Total = searchRsp.GetTotal()
	rsp.TotalYears = searchRsp.GetTotalYears()

	return nil
}

// 验证request
func validateSearchTenantSubscriptionDetailsRequest(req *pb.SearchTenantSubscriptionDetailsRequest) error {
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
