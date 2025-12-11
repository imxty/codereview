package user

import (
	"context"
	gerr "errors"

	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) GetOrganizationOverdueSummary(ctx context.Context, req *pb.GetOrganizationOverdueSummaryRequest, rsp *pb.GetOrganizationOverdueSummaryResponse) error {
	err := validateGetOrganizationOverdueSummary(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := u.customerAPI.GetOrganizationOverdueCount(ctx, &customerv1.GetOrganizationOverdueCountRequest{
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return err
	}

	rsp.OverdueCustomerTotalCount = getRsp.GetTotalCount()

	return nil
}

// 验证 request
func validateGetOrganizationOverdueSummary(req *pb.GetOrganizationOverdueSummaryRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
