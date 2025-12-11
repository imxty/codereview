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

func (s *H5APIHandler) GetReport(ctx context.Context, req *pb.GetReportRequest, rsp *pb.GetReportResponse) error {
	err := validateGetReportRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 构建report service的answer
	m := make(map[string]*reportpb.AnswerList)
	for k, v := range req.GetModuleAnswers() {
		m[k] = toReportAnswerList(v)
	}
	getReportRsp, err := s.reportAPI.GetReport(ctx, &reportpb.GetReportRequest{
		// 租户id
		TenantId: req.GetTenantId(),
		// 报告ID
		ReportId: req.GetReportId(),
		// 语言
		LanguageCode: req.GetLanguageCode(),
		// 回答的问题，是模块名到回答的问题的映射关系
		ModuleAnswers: m,
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 获取开关
	getRsp, err := s.userAPI.GetTenantEntity(ctx, &userv1.GetTenantEntityRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.HasQuestions = getReportRsp.GetHasQuestions()
	rsp.IsCompleteReport = getReportRsp.GetIsCompleteReport()
	// 如果有问题就返回问题
	if rsp.GetHasQuestions() {
		questions := make(map[string]*pb.QuestionList)
		for k, v := range getReportRsp.GetModuleQuestions() {
			questions[k] = toAppQuestionList(v)
		}
		rsp.ModuleQuestions = questions
	}
	report := toAppReport(getReportRsp.GetReport(), s.s3Domain)
	// 如果体质辨证开关关闭则不返回体质
	if !getRsp.GetConstitutionSwitchStatus() && report.TcmReport != nil {
		report.TcmReport.PhysiqueDialecticsModule = nil
	}
	rsp.Report = report

	// 返回舌面报告
	if getReportRsp.GetTongueFaceReport() != nil {
		tfReport := getReportRsp.GetTongueFaceReport()
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
func validateGetReportRequest(req *pb.GetReportRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetReportId() == "" {
		return gerr.New("report_id should not be empty")
	}
	return nil
}
