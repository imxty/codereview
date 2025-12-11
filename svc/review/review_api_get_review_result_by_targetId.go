package review

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ReviewAPIHandler) GetReviewResultByTargetId(ctx context.Context, req *pb.GetReviewResultByTargetIdRequest, rsp *pb.GetReviewResultByTargetIdResponse) error {
	err := validateGetReviewResultByTargetIdRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 通过target id获取result
	issue, err := s.reviewStore.GetLatestReviewIssueByTargetId(ctx, req.GetTargetId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if issue == nil {
		rsp.Result = &pb.ReviewResult{}
		return nil
	}
	// 通过issue获取result
	result, err := s.reviewStore.GetReviewResultByIssueId(ctx, issue.GetReviewIssueID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 说明还没有审核完成
	if result == nil {
		rsp.Result = &pb.ReviewResult{
			// 审核类型
			ReviewType: toProtoReviewType(issue.GetTargetType()),
			// 审核状态
			ReiewStatus: toProtoReviewStatus(issue.GetReviewStatus()),
		}
		return nil
	}
	rsp.Result = &pb.ReviewResult{
		// 审核类型
		ReviewType: toProtoReviewType(issue.GetTargetType()),
		// 审核状态
		ReiewStatus: toProtoReviewStatus(issue.GetReviewStatus()),
		// 如果审核失败,失败原因
		FailReason: result.GetComment(),
	}
	return nil
}

func validateGetReviewResultByTargetIdRequest(req *pb.GetReviewResultByTargetIdRequest) error {
	if req.GetTargetId() == "" {
		return gerr.New("target_id should not be empty")
	}
	return nil
}
