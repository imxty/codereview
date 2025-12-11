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

// 查询商户
func (s *BossAPIHandler) SearchTenant(ctx context.Context, req *pb.SearchTenantRequest, rsp *pb.SearchTenantResponse) error {
	err := validateSearchTenantRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	if req.GetOrganizationId() == "" {
		rsp.Tenants = []*pb.SubscriptionPeriod{}
		rsp.TotalCount = 0
		return nil
	}

	searchRsp, err := s.userAPI.SearchTenantDetails(ctx, &userv1.SearchTenantDetailsRequest{
		OrganizationId: req.GetOrganizationId(),
		TenantName:     req.GetTenantName(),
		StartTime:      req.GetStartTime(),
		EndTime:        req.GetEndTime(),
		ContactName:    req.GetContactName(),
		ContactPhone:   req.GetContactPhone(),
		Status:         toSvcStatus(req.GetStatus()),
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

	// 返回数据
	rsp.Tenants = results
	rsp.TotalCount = searchRsp.GetTotal()

	return nil
}

// 验证request
func validateSearchTenantRequest(req *pb.SearchTenantRequest) error {
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
