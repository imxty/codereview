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

// 修改体质辩证开关
func (s *TenantAPIHandler) UpdateConstitutionSwitchStatus(ctx context.Context, req *pb.UpdateConstitutionSwitchStatusRequest, rsp *pb.UpdateConstitutionSwitchStatusResponse) error {
	err := validateUpdateConstitutionSwitchStatusRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	tenant, err := s.userAPI.GetTenantEntity(ctx, &userv1.GetTenantEntityRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	if tenant.GetConstitutionSwitchStatus() == req.GetConstitutionSwitchStatus() {
		rsp.ConstitutionSwitchStatus = tenant.GetConstitutionSwitchStatus()
		return nil
	}

	updateRsp, err := s.userAPI.UpdateTenantConstitutionSwitch(ctx, &userv1.UpdateTenantConstitutionSwitchRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.ConstitutionSwitchStatus = updateRsp.GetConstitutionSwitchStatus()

	return nil
}

// 验证request
func validateUpdateConstitutionSwitchStatusRequest(req *pb.UpdateConstitutionSwitchStatusRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
