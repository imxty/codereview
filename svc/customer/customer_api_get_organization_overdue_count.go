package customer

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	OverdueDays = 30
)

func (s *CustomerAPIHandler) GetOrganizationOverdueCount(ctx context.Context, req *pb.GetOrganizationOverdueCountRequest, rsp *pb.GetOrganizationOverdueCountResponse) error {
	err := validateGetOrganizationOverdueCountRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	count, err := s.customerStore.GetOrganizationOverdueCount(ctx, req.GetOrganizationId(), OverdueDays)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.TotalCount = int32(count)

	return nil
}

func validateGetOrganizationOverdueCountRequest(req *pb.GetOrganizationOverdueCountRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
