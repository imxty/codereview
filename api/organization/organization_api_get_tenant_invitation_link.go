package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取商户邀请链接
func (s *OrganizationAPIHandler) GetTenantInvitationLink(ctx context.Context, req *pb.GetTenantInvitationLinkRequest, rsp *pb.GetTenantInvitationLinkResponse) error {
	err := validateGetTenantInvitationLinkRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	return nil
}

// 验证request
func validateGetTenantInvitationLinkRequest(req *pb.GetTenantInvitationLinkRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
