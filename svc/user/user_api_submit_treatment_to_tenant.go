package user

import (
	"context"
	gerr "errors"
	"slices"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 分配方案
func (u *UserAPIHandler) SubmitTreatmentToTenant(ctx context.Context, req *pb.SubmitTreatmentToTenantRequest, rsp *pb.SubmitTreatmentToTenantResponse) error {
	err := validateSubmitTreatmentToTenant(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取已有方案的商户
	tts, err := u.userStore.GetTenantTreatments(ctx, req.GetTenantIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 所有已有方案的商户ID
	hasIds := make([]string, 0)
	for _, each := range tts {
		hasIds = append(hasIds, each.GetTenantID())
	}

	// 所有没有方案的商户ID
	noIds := make([]string, 0)
	for _, each := range req.GetTenantIds() {
		if !slices.Contains(hasIds, each) {
			noIds = append(noIds, each)
		}
	}

	// 分配方案
	if len(hasIds) > 0 {
		err = u.userStore.UpdateTreatmentToTenant(ctx, req.GetOrganizationId(), hasIds, req.GetTreatmentId())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	if len(noIds) > 0 {
		err = u.userStore.SubmitTreatmentToTenant(ctx, req.GetOrganizationId(), noIds, req.GetTreatmentId())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	return nil
}

// 验证 request
func validateSubmitTreatmentToTenant(req *pb.SubmitTreatmentToTenantRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetTenantIds() == nil {
		return gerr.New("tenant_ids should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment_id should not be empty")
	}
	return nil
}
