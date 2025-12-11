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

func (s *AppAPIHandler) GetSharedReportLink(ctx context.Context, req *pb.GetSharedReportLinkRequest, rsp *pb.GetSharedReportLinkResponse) error {
	err := validateGetSharedReportLinkRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取分享报告的连接
	getRsp, err := s.reportAPI.GetSharedReportLink(ctx, &reportpb.GetSharedReportLinkRequest{
		ReportId: req.GetReportId(),
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	rsp.ReportLink = getRsp.GetReportLink()
	return nil
}

// 验证request
func validateGetSharedReportLinkRequest(req *pb.GetSharedReportLinkRequest) error {
	if req.GetReportId() == "" {
		return gerr.New("report_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
