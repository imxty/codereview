package user

import (
	"context"
	gerr "errors"

	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
)

// GetSharingReportSwitchStatus 获取报告分享开关状态
func (u *UserAPIHandler) GetSharingReportSwitchStatus(ctx context.Context, req *pb.GetSharingReportSwitchStatusRequest, rsp *pb.GetSharingReportSwitchStatusResponse) error {

	// 1.验证request
	err := validateGetSharingReportSwitchStatusRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询租户
	tenant, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenant == nil {
		return errors.Error(ErrTenantNotFound, "tenant not found")
	}

	// 获取租户信息
	tenantEntity, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenantEntity == nil {
		return errors.Error(ErrTenantNotFound, "tenant_entity not found")
	}

	rsp.SharingReportSwitchStatus = tenantEntity.GetReportSharingStatus()

	return nil
}

// 验证request
func validateGetSharingReportSwitchStatusRequest(req *pb.GetSharingReportSwitchStatusRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	return nil
}
