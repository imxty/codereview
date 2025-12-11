package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) ModifyReportRemark(ctx context.Context, req *pb.ModifyReportRemarkRequest, rsp *pb.ModifyReportRemarkResponse) error {
	err := validateModifyReportRemarkRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 发送请求
	_, err = s.reportAPI.ModifyReportRemark(ctx, &reportpb.ModifyReportRemarkRequest{
		ReportId: req.GetReportId(),
		Remark:   req.GetRemark(),
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateModifyReportRemarkRequest(req *pb.ModifyReportRemarkRequest) error {
	if req.GetReportId() == "" {
		return gerr.New("report_id should  not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	return nil
}
