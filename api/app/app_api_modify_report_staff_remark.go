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

func (s *AppAPIHandler) ModifyReportStaffRemark(ctx context.Context, req *pb.ModifyReportStaffRemarkRequest, rsp *pb.ModifyReportStaffRemarkResponse) error {
	err := validateModifyReportStaffRemarkRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 发送请求
	_, err = s.reportAPI.ModifyReportStaffRemark(ctx, &reportpb.ModifyReportStaffRemarkRequest{
		ReportId:    req.GetReportId(),
		StaffRemark: req.GetStaffRemark(),
		TenantId:    req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateModifyReportStaffRemarkRequest(req *pb.ModifyReportStaffRemarkRequest) error {
	if req.GetReportId() == "" {
		return gerr.New("report_id should  not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	return nil
}
