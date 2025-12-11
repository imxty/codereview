package report

import (
	"context"
	"encoding/json"
	gerr "errors"
	"math/rand"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ReportAPIHandler) GetInquiryQuestions(ctx context.Context, req *pb.GetInquiryQuestionsRequest, rsp *pb.GetInquiryQuestionsResponse) error {
	err := validateGetInquiryQuestionsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取所有问诊问题
	questions, err := s.reportStore.GetInquiryQuestions(ctx)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 如果问题数量不够，返回非法操作
	if questions == nil || len(questions) < int(req.GetCount()) {
		return errors.Errorf(codes.InvalidOperation, "question count not enough, want %d got %d", req.GetCount(), len(questions))
	}

	// 随机洗牌，获取指定数量问题
	// selections := shuffle(questions, req.GetCount())

	// 不随机抽取，按 ID 顺序返回问题
	selections := questions

	results := make([]*pb.InquiryQuestion, len(selections))
	for index, selection := range selections {
		// 解析选项json字符串
		var items []domain.InquiryAnswerItem
		if err := json.Unmarshal([]byte(selection.GetSelections()), &items); err != nil {
			return errors.Errorf(codes.Internal, "unmarshal question[%s] selections error: %v", selection.GetInquiryID(), err)
		}
		apiItems := make([]*pb.InquiryAnswerItem, len(items))
		for k, v := range items {
			apiItems[k] = toApiInquiryAnswerItem(v)
		}

		results[index] = &pb.InquiryQuestion{
			InquiryId:           selection.GetInquiryID(),
			Content:             selection.GetContent(),
			IsMultipleSelection: selection.GetIsMultipleSelection(),
			Items:               apiItems,
		}
	}

	// 返回问诊问题
	rsp.Questions = results
	return nil
}

func validateGetInquiryQuestionsRequest(req *pb.GetInquiryQuestionsRequest) error {
	if req.GetCount() < 0 {
		return gerr.New("count should not less than 0")
	}
	return nil
}

// shuffle 洗牌算法，随机选取指定个数问题
func shuffle(questions []domain.InquiryDiagnosisIntf, count int32) []domain.InquiryDiagnosisIntf {
	shuffled := make([]domain.InquiryDiagnosisIntf, len(questions))
	copy(shuffled, questions)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled[:count]
}
