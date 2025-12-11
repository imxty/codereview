package h5

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/h5/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *H5APIHandler) GetTongueFaceReport(ctx context.Context, req *pb.GetTongueFaceReportRequest, rsp *pb.GetTongueFaceReportResponse) error {
	err := validateGetTongueFaceReportRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.reportAPI.GetTongueFaceReport(ctx, &reportpb.GetTongueFaceReportRequest{
		TongueImageUrl: req.GetTongueImageUrl(),
		FaceImageUrl:   req.GetFaceImageUrl(),
		ReportId:       req.GetReportId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	report := getRsp.GetTongueFaceReport()
	rsp.TongueFaceReport = &pb.TongueFaceReport{
		Success:        report.GetSuccess(),
		FaceImageUrl:   s.getS3ImgUrl(report.GetFaceImageUrl()),
		TongueImageUrl: s.getS3ImgUrl(report.GetTongueImageUrl()),
		Data:           report.GetData(),
	}

	return nil
}

// 验证request
func validateGetTongueFaceReportRequest(req *pb.GetTongueFaceReportRequest) error {
	if req.GetReportId() == "" {
		return gerr.New("report id should not be empty")
	}
	if req.GetFaceImageUrl() == "" {
		return gerr.New("face image url should not be empty")
	}
	if req.GetTongueImageUrl() == "" {
		return gerr.New("tongue image url should not be empty")
	}
	return nil
}
