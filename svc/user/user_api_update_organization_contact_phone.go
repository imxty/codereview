package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) UpdateOrganizationContactPhone(ctx context.Context, req *pb.UpdateOrganizationContactPhoneRequest, rsp *pb.UpdateOrganizationContactPhoneResponse) error {
	err := validateUpdateOrganizationContactPhoneRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询组织信息
	og, err := u.userStore.GetOrganizationByOrganizationId(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if og == nil {
		return errors.Errorf(ErrOrganizationNotExist, "organization[%s] not found", req.GetOrganizationId())
	}

	// 更新组织联系人电话
	err = u.userStore.UpdateOrganizationContactPhone(ctx, req.GetOrganizationId(), req.GetOrganizationContactPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
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
