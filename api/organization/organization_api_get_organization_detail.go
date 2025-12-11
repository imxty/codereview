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

// 获取组织详情
func (s *OrganizationAPIHandler) GetOrganizationDetail(ctx context.Context, req *pb.GetOrganizationDetailRequest, rsp *pb.GetOrganizationDetailResponse) error {
	err := validateGetOrganizationDetailRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取组织详情
	getRsp, err := s.userAPI.GetOrganizationDetail(ctx, &userv1.GetOrganizationDetailRequest{
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回响应
	rsp.Organization = toApiOrganization(getRsp.GetOrganization())

	return nil
}

// 验证request
func validateGetOrganizationDetailRequest(req *pb.GetOrganizationDetailRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
