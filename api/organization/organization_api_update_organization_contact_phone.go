package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *OrganizationAPIHandler) UpdateOrganizationContactPhone(ctx context.Context, req *pb.UpdateOrganizationContactPhoneRequest, rsp *pb.UpdateOrganizationContactPhoneResponse) error {
	err := validateUpdateOrganizationContactPhoneRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.UpdateOrganizationContactPhone(ctx, &userpb.UpdateOrganizationContactPhoneRequest{
		OrganizationId:           req.GetOrganizationId(),
		OrganizationContactPhone: req.GetOrganizationContactPhone(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

func validateUpdateOrganizationContactPhoneRequest(req *pb.UpdateOrganizationContactPhoneRequest) error {
	if req.GetOrganizationContactPhone() == "" {
		return gerr.New("organization_contact_phone should not be empty")
	}
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
