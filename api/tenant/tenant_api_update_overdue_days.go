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

// 修改复查天数
func (s *TenantAPIHandler) UpdateOverdueDays(ctx context.Context, req *pb.UpdateOverdueDaysRequest, rsp *pb.UpdateOverdueDaysResponse) error {
	err := validateUpdateOverdueDaysRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	updateRsp, err := s.userAPI.UpdateOverdueDays(ctx, &userv1.UpdateOverdueDaysRequest{
		TenantId: req.GetTenantId(),
		Days:     req.GetDays(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.Days = updateRsp.GetDays()

	return nil
}

// 验证request
func validateUpdateOverdueDaysRequest(req *pb.UpdateOverdueDaysRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
