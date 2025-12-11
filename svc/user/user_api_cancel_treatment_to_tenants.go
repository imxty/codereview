package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) CancelTreatmentToTenants(ctx context.Context, req *pb.CancelTreatmentToTenantsRequest, rsp *pb.CancelTreatmentToTenantsResponse) error {
	err := validateCancelTreatmentToTenantsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	if len(req.GetTenantIds()) < 1 {
		return nil
	}

	err = u.userStore.CancelTreatmentToTenants(ctx, req.GetTenantIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
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
