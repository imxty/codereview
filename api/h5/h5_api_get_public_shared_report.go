package h5

import (
	"context"
	gerr "errors"
	"net/url"
	"path"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/h5/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *H5APIHandler) GetPublicSharedReport(ctx context.Context, req *pb.GetPublicSharedReportRequest, rsp *pb.GetPublicSharedReportResponse) error {
	err := validateGetPublicSharedReportRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取公开分享的报告请求
	getRsp, err := s.reportAPI.GetPublicSharedReport(ctx, &reportpb.GetPublicSharedReportRequest{
		// 租户id
		TenantId: req.GetTenantId(),
		// 公开分享报告的 token
		Token: req.GetToken(),
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

	// 返回舌面报告
	if getRsp.GetTongueFaceReport() != nil {
		tfReport := getRsp.GetTongueFaceReport()
		rsp.TongueFaceReport = &pb.TongueFaceReport{
			Success:        tfReport.GetSuccess(),
			FaceImageUrl:   s.getS3ImgUrl(tfReport.GetFaceImageUrl()),
			TongueImageUrl: s.getS3ImgUrl(tfReport.GetTongueImageUrl()),
			Data:           tfReport.GetData(),
		}
	}
	return nil
}

// 验证request
func validateGetPublicSharedReportRequest(req *pb.GetPublicSharedReportRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetToken() == "" {
		return gerr.New("token should not be empty")
	}
	return nil
}

func (s *H5APIHandler) getS3ImgUrl(imgUrl string) string {
	if imgUrl == "" {
		return ""
	}
	link, _ := url.Parse(s.s3Domain)
	link.Path = path.Join(link.Path, imgUrl)
	return link.String()
}
