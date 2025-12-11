package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 删除商户
func (s *BossAPIHandler) DeleteTenant(ctx context.Context, req *pb.DeleteTenantRequest, rsp *pb.DeleteTenantResponse) error {
	err := validateDeleteTenantRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.DeleteTenant(ctx, &userv1.DeleteTenantRequest{
		UserId:         req.GetUserId(),
		TenantId:       req.GetTenantId(),
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

// 验证request
func validateDeleteTenantRequest(req *pb.DeleteTenantRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
