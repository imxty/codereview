package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 通过商户名和组织 ID 模糊查询商户信息
func (u *UserAPIHandler) SearchTenantByNameAndOrganizationID(ctx context.Context, req *pb.SearchTenantByNameAndOrganizationIDRequest, rsp *pb.SearchTenantByNameAndOrganizationIDResponse) error {
	err := validateSearchTenantByNameAndOrganizationIDRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	tenants, err := u.userStore.SearchTenantByNameAndOrganizationID(ctx, req.GetOrganizationId(), req.GetTenantName())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.TenantIds = tenants

	return nil
}

// 验证 request
func validateSearchTenantByNameAndOrganizationIDRequest(req *pb.SearchTenantByNameAndOrganizationIDRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
