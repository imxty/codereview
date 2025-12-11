package h5

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/h5/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *H5APIHandler) GetWeeklyReport(ctx context.Context, req *pb.GetWeeklyReportRequest, rsp *pb.GetWeeklyReportResponse) error {
	err := validateGetWeeklyReportRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取公开分享的慧脉宝药报告请求
	getRsp, err := s.reportAPI.GetWeeklyReport(ctx, &reportpb.GetWeeklyReportRequest{
		// 租户id
		TenantId: req.GetTenantId(),
		// 常客ID
		CustomerId: req.GetCustomerId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 获取开关
	getTenantRsp, err := s.userAPI.GetTenantEntity(ctx, &userv1.GetTenantEntityRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回报告
	r := toAppReport(getRsp.GetReport(), s.s3Domain)
	// 如果体质辨证开关关闭则不返回体质
	if !getTenantRsp.GetConstitutionSwitchStatus() && r.TcmReport != nil {
		r.TcmReport.PhysiqueDialecticsModule = nil
	}
	rsp.Report = r
	rsp.ReportEnough = getRsp.GetReportEnough()

	return nil
}

// 验证request
func validateGetWeeklyReportRequest(req *pb.GetWeeklyReportRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomerId() == "" {
		return gerr.New("customer id should not be empty")
	}
	return nil
}
