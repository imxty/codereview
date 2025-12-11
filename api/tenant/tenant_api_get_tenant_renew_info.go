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

func (s *TenantAPIHandler) GetTenantRenewInfo(ctx context.Context, req *pb.GetTenantRenewInfoRequest, rsp *pb.GetTenantRenewInfoResponse) error {
	err := validateGetTenantRenewInfoRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetTenantRenewInfo(ctx, &userv1.GetTenantRenewInfoRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.IsExpiredSoon = getRsp.GetIsExpiredSoon()
	rsp.EndTime = getRsp.GetEndTime()
	rsp.OrganizationContactPhone = getRsp.GetOrganizationContactPhone()

	return nil
}

func validateGetTenantRenewInfoRequest(req *pb.GetTenantRenewInfoRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
