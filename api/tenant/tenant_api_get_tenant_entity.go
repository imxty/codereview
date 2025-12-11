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

func (s *TenantAPIHandler) GetTenantEntity(ctx context.Context, req *pb.GetTenantEntityRequest, rsp *pb.GetTenantEntityResponse) error {
	err := validateGetTenantEntityRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetTenantEntity(ctx, &userv1.GetTenantEntityRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回商户信息
	rsp.Entity = toAppTenantEntity(getRsp.GetEntity(), s.s3Domain)
	rsp.ReportSharingSwitchStatus = getRsp.GetReportSharingSwitchStatus()
	rsp.ConstitutionSwitchStatus = getRsp.GetConstitutionSwitchStatus()
	rsp.OverdueDays = getRsp.GetOverdueDays()
	return nil
}

// 验证request
func validateGetTenantEntityRequest(req *pb.GetTenantEntityRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	return nil
}
