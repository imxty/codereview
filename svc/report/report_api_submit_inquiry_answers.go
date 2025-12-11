package report

import (
	"context"
	"encoding/json"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// SubmitInquiryAnswers 提交问诊回答
func (s *ReportAPIHandler) SubmitInquiryAnswers(ctx context.Context, req *pb.SubmitInquiryAnswersRequest, rsp *pb.SubmitInquiryAnswersResponse) error {
	err := validateSubmitInquiryAnswersRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 将提交的答案转换为json字符串
	answers := make([]domain.InquiryAnswer, len(req.GetAnswers()))
	for k, v := range req.GetAnswers() {
		answers[k] = domain.InquiryAnswer{
			InquiryID: v.GetInquiryId(),
			Answers:   v.GetAnswers(),
		}
	}
	data, err := json.Marshal(answers)
	if err != nil {
		return errors.Errorf(codes.Internal, "marshal answer json error: %v", err)
	}

	// 存储问诊信息
	err = s.reportStore.ModifyReportInquiryDiagnosis(ctx, req.GetReportId(), string(data))
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

func validateSubmitInquiryAnswersRequest(req *pb.SubmitInquiryAnswersRequest) error {
	if req.GetReportId() == "" {
		return gerr.New("report_id should not be empty")
	}
	if req.GetAnswers() == nil {
		return gerr.New("answers should not be nil")
	}
	return nil
}
