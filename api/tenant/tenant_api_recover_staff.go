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

func (s *TenantAPIHandler) RecoverStaff(ctx context.Context, req *pb.RecoverStaffRequest, rsp *pb.RecoverStaffResponse) error {
	err := validateRecoverStaffRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.RecoverStaff(ctx, &userv1.RecoverStaffRequest{
		TenantId: req.GetTenantId(),
		StaffId:  req.GetStaffId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

func validateRecoverStaffRequest(req *pb.RecoverStaffRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetStaffId() == "" {
		return gerr.New("staff_id shoudl not be empty")
	}
	return nil
}
