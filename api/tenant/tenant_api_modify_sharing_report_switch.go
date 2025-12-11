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

// 修改报告分享开关
func (s *TenantAPIHandler) ModifySharingReportSwitch(ctx context.Context, req *pb.ModifySharingReportSwitchRequest, rsp *pb.ModifySharingReportSwitchResponse) error {
	err := validateModifySharingReportSwitchRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	modifyRsp, err := s.userAPI.ModifySharingReportSwitch(ctx, &userv1.ModifySharingReportSwitchRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.ReportSharingSwitchStatus = modifyRsp.ReportSharingSwitchStatus

	return nil
}

// 验证request
func validateModifySharingReportSwitchRequest(req *pb.ModifySharingReportSwitchRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
