package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *OrganizationAPIHandler) SubmitTreatmentToTenant(ctx context.Context, req *pb.SubmitTreatmentToTenantRequest, rsp *pb.SubmitTreatmentToTenantResponse) error {
	err := validateSubmitTreatmentToTenantRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.SubmitTreatmentToTenant(ctx, &userv1.SubmitTreatmentToTenantRequest{
		OrganizationId: req.GetOrganizationId(),
		TenantIds:      req.GetTenantIds(),
		TreatmentId:    req.GetTreatmentId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

// 验证request
func validateSubmitTreatmentToTenantRequest(req *pb.SubmitTreatmentToTenantRequest) error {
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
