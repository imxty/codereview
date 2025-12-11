package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) GetStaffNamesByID(ctx context.Context, req *pb.GetStaffNamesByIDRequest, rsp *pb.GetStaffNamesByIDResponse) error {
	err := validateGetStaffNamesByIDRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetStaffNamesByID(ctx, &userv1.GetStaffNamesByIDRequest{
		StaffIds: req.GetStaffIds(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.StaffNames = getRsp.GetStaffNames()

	return nil
}

func validateGetStaffNamesByIDRequest(req *pb.GetStaffNamesByIDRequest) error {
	if req.GetStaffIds() == nil {
		return gerr.New("staff_ids should not be nil")
	}
	return nil
}
