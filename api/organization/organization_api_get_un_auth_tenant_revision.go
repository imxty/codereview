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

// 获取首次认证失败的商户信息
func (s *OrganizationAPIHandler) GetUnAuthTenantRevision(ctx context.Context, req *pb.GetUnAuthTenantRevisionRequest, rsp *pb.GetUnAuthTenantRevisionResponse) error {
	err := validateGetUnAuthTenantRevisionRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetUnAuthTenantRevision(ctx, &userv1.GetUnAuthTenantRevisionRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	rsp.HasUnauthRevision = getRsp.GetHasUnauthRevision()
	rsp.Tenant = toApiTenant(getRsp.GetTenant(), s.s3Domain)

	return nil
}

// 验证request
func validateGetUnAuthTenantRevisionRequest(req *pb.GetUnAuthTenantRevisionRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
