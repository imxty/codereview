package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	reportv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) SubmitInquiryAnswers(ctx context.Context, req *pb.SubmitInquiryAnswersRequest, rsp *pb.SubmitInquiryAnswersResponse) error {
	err := validateSubmitInquiryAnswersRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	answers := make([]*reportv1.InquiryAnswer, len(req.GetAnswers()))
	for k, v := range req.GetAnswers() {
		answers[k] = toSvcAnswer(v)
	}

	_, err = s.reportAPI.SubmitInquiryAnswers(ctx, &reportv1.SubmitInquiryAnswersRequest{
		ReportId: req.GetReportId(),
		Answers:  answers,
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
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
