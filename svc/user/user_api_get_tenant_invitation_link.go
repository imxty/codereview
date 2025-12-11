package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) GetTenantInvitationLink(ctx context.Context, req *pb.GetTenantInvitationLinkRequest, rsp *pb.GetTenantInvitationLinkResponse) error {
	err := validateGetTenantInvitationLinkRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	return nil
}

func validateGetTenantInvitationLinkRequest(req *pb.GetTenantInvitationLinkRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	return nil
}
