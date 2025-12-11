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

func (s *AppAPIHandler) GetInquiryQuestions(ctx context.Context, req *pb.GetInquiryQuestionsRequest, rsp *pb.GetInquiryQuestionsResponse) error {
	err := validateGetInquiryQuestionsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.reportAPI.GetInquiryQuestions(ctx, &reportv1.GetInquiryQuestionsRequest{
		Count: req.GetCount(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	results := make([]*pb.InquiryQuestion, len(getRsp.GetQuestions()))
	for k, v := range getRsp.GetQuestions() {
		results[k] = toAppInquiryQuestion(v)
	}

	rsp.Questions = results
	return nil
}

func validateGetInquiryQuestionsRequest(req *pb.GetInquiryQuestionsRequest) error {
	if req.GetCount() < 0 {
		return gerr.New("count should not less than 0")
	}
	return nil
}
