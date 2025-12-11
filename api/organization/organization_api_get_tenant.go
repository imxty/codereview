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

// 获取商户
func (s *OrganizationAPIHandler) GetTenant(ctx context.Context, req *pb.GetTenantRequest, rsp *pb.GetTenantResponse) error {
	err := validateGetTenantRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetTenant(ctx, &userv1.GetTenantRequest{
		OrganizationId: req.GetOrganizationId(),
		TenantId:       req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	rsp.Tenant = toApiTenant(getRsp.GetTenant(), s.s3Domain)

	return nil
}

// 验证request
func validateGetTenantRequest(req *pb.GetTenantRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
