package user

import (
	"context"
	gerr "errors"

	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
)

// ModifySharingReportSwitch 修改报告分享开关请求
func (u *UserAPIHandler) ModifySharingReportSwitch(ctx context.Context, req *pb.ModifySharingReportSwitchRequest, rsp *pb.ModifySharingReportSwitchResponse) error {

	// 1.验证request
	err := validateModifySharingReportSwitchRequest(req)
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

	// 修改报告分享开关请求
	err = u.userStore.ModifyReportSharingStatus(ctx, req.GetTenantId(), !tenantEntity.GetReportSharingStatus(), tenantEntity.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.ReportSharingSwitchStatus = !tenantEntity.GetReportSharingStatus()

	return nil
}

// 验证request
func validateModifySharingReportSwitchRequest(req *pb.ModifySharingReportSwitchRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	return nil
}
