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

func (s *OrganizationAPIHandler) CancelTreatmentToTenants(ctx context.Context, req *pb.CancelTreatmentToTenantsRequest, rsp *pb.CancelTreatmentToTenantsResponse) error {
	err := validateCancelTreatmentToTenantsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	if len(req.GetTenantIds()) < 1 {
		return nil
	}

	_, err = s.userAPI.CancelTreatmentToTenants(ctx, &userv1.CancelTreatmentToTenantsRequest{
		OrganizationId: req.GetOrganizationId(),
		TenantIds:      req.GetTenantIds(),
		TreatmentId:    req.GetTreatmentId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

func validateCancelTreatmentToTenantsRequest(req *pb.CancelTreatmentToTenantsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetTenantIds() == nil {
		return gerr.New("tenant_ids should not be nil")
	}
	return nil
}
