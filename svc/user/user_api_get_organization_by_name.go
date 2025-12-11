package user

import (
	"context"
	gerr "errors"

	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
)

// 通过名称获取组织（模糊查询）
func (u *UserAPIHandler) GetOrganizationByName(ctx context.Context, req *pb.GetOrganizationByNameRequest, rsp *pb.GetOrganizationByNameResponse) error {
	err := validateGetOrganizationByNameRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	organizations, err := u.userStore.GetOrganizationByVagueUsername(ctx, req.GetOrganizationName())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if len(organizations) == 0 {
		return nil
	}

	// 改成数组
	os := make([]*pb.Organization, len(organizations))
	for k, v := range organizations {
		os[k] = &pb.Organization{
			// 组织 ID
			OrganizationId: v.GetOrganizationID(),
			// 名称
			Name: v.GetUsername(),
			// 手机号
			Phone: v.GetPhone(),
		}
	}

	rsp.Organizations = os

	return nil
}

func validateGetOrganizationByNameRequest(req *pb.GetOrganizationByNameRequest) error {
	if req.GetOrganizationName() == "" {
		return gerr.New("organization name should not be empty")
	}
	return nil
}
