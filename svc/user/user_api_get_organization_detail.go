package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取组织详情
func (u *UserAPIHandler) GetOrganizationDetail(ctx context.Context, req *pb.GetOrganizationDetailRequest, rsp *pb.GetOrganizationDetailResponse) error {
	err := validateGetOrganizationDetailRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取组织信息
	organization, err := u.userStore.GetOrganizationByOrganizationId(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if organization == nil {
		return errors.Errorf(ErrOrganizationNotExist, "organization not found [%s]", req.GetOrganizationId())
	}

	// 查询组织下的商户
	tf, _, err := u.userStore.ListOrganizationTenantsWithoutTenantId(ctx, organization.GetOrganizationID(), 0, 1)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	hasTenant := false
	if len(tf) != 0 {
		hasTenant = true
	}

	// 返回组织信息
	rsp.Organization = &pb.Organization{
		// 组织 ID
		OrganizationId: req.GetOrganizationId(),
		// 名称
		Name: organization.GetUsername(),
		// 手机号
		Phone:                    organization.GetPhone(),
		HasTenant:                hasTenant,
		OrganizationContactPhone: organization.GetOrganizationContactPhone(),
	}
	return nil
}

// 验证 request
func validateGetOrganizationDetailRequest(req *pb.GetOrganizationDetailRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization Id should not be empty")
	}
	return nil
}
