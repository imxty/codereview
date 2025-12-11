package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取组织商户概况
func (u *UserAPIHandler) GetOrganizationTenantSummary(ctx context.Context, req *pb.GetOrganizationTenantSummaryRequest, rsp *pb.GetOrganizationTenantSummaryResponse) error {
	err := validGetOrganizationTenantSummaryRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	tenantIds, err := u.userStore.ListAllOrganizationTenantIds(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	usingTenantIds, err := u.userStore.GetUsingTenantIds(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	todayAddedTenantIds, err := u.userStore.GetTodayAddedTenantIds(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	expiringTenantIds, err := u.userStore.GetExpiringTenantIds(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回响应
	rsp.TenantExpiringSoonCount = int32(len(expiringTenantIds))
	rsp.TenantTotalCount = int32(len(tenantIds))
	rsp.TenantUsingCount = int32(len(usingTenantIds))
	rsp.TodayAddedTenantCount = int32(len(todayAddedTenantIds))

	return nil
}

// 验证 request
func validGetOrganizationTenantSummaryRequest(req *pb.GetOrganizationTenantSummaryRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should")
	}
	return nil
}
