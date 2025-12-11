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

func (s *OrganizationAPIHandler) CommitEntityCertificate(ctx context.Context, req *pb.CommitEntityCertificateRequest, rsp *pb.CommitEntityCertificateResponse) error {
	err := validateCommitEntityCertificateRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.CommitEntityCertificate(ctx, &userv1.CommitEntityCertificateRequest{
		// 组织ID
		OrganizationId: req.GetOrganizationId(),
		// 商户信息
		Tenant: toUserTenantEntity(req.GetTenant()),
		// 是否是组织提审
		IsOrganization: true,
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

// 验证request
func validateCommitEntityCertificateRequest(req *pb.CommitEntityCertificateRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTenant() == nil {
		return gerr.New("tenant should not be nil")
	}
	return nil
}
