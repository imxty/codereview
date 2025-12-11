package customer

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *CustomerAPIHandler) GetOverdueCount(ctx context.Context, req *pb.GetOverdueCountRequest, rsp *pb.GetOverdueCountResponse) error {
	err := validateGetOverdueCount(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetTenantEntity(ctx, &userv1.GetTenantEntityRequest{
		TenantId: req.TenantId,
	})
	if err != nil {
		return err
	}

	tenant := getRsp.GetEntity()

	count, err := s.customerStore.GetTenantOverdueCount(ctx, tenant.GetTenantId(), tenant.GetOverdue())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.TotalCount = int32(count)

	return nil
}

func validateGetOverdueCount(req *pb.GetOverdueCountRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
