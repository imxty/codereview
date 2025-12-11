package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 修改商户体质辩证开关请求
func (u *UserAPIHandler) UpdateTenantConstitutionSwitch(ctx context.Context, req *pb.UpdateTenantConstitutionSwitchRequest, rsp *pb.UpdateTenantConstitutionSwitchResponse) error {
	err := validateUpdateTenantConstitutionSwitchRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商户信息
	tenant, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if tenant == nil {
		return errors.Errorf(ErrTenantNotFound, "tenant[%s] not found", req.GetTenantId())
	}

	// 修改体质辩证开关
	err = u.userStore.UpdateTenantConstitutionSwitch(ctx, req.GetTenantId(), !tenant.GetConstitutionStatus())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.ConstitutionSwitchStatus = !tenant.GetConstitutionStatus()

	return nil
}

// 验证 request
func validateUpdateTenantConstitutionSwitchRequest(req *pb.UpdateTenantConstitutionSwitchRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
