package user

import (
	"context"
	gerr "errors"

	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
)

// 获取组织名称
func (u *UserAPIHandler) GetOrganizationsName(ctx context.Context, req *pb.GetOrganizationsNameRequest, rsp *pb.GetOrganizationsNameResponse) error {
	err := validateGetOrganizationsNameRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询组织信息
	us, err := u.userStore.ListOrganizations(ctx, req.GetOrganizationIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	oname := make(map[string]*pb.Organization)
	for _, v := range us {
		oname[v.GetOrganizationID()] = &pb.Organization{
			OrganizationId: v.GetOrganizationID(),
			Phone:          v.GetPhone(),
			Name:           v.GetUsername(),
		}
	}

	rsp.OrganizationName = oname
	return nil
}

func validateGetOrganizationsNameRequest(req *pb.GetOrganizationsNameRequest) error {
	if len(req.GetOrganizationIds()) == 0 {
		return gerr.New("organization ids should not be nil")
	}
	return nil
}
