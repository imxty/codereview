package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	reviewv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *BossAPIHandler) CommitTreatmentReviewResult(ctx context.Context, req *pb.CommitTreatmentReviewResultRequest, rsp *pb.CommitTreatmentReviewResultResponse) error {
	err := validateCommitTreatmentReviewResultRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 发送请求
	commitRsp, err := s.reviewAPI.CommitTreatmentReviewResult(ctx, &reviewv1.CommitTreatmentReviewResultRequest{
		ReviewIssueId:  req.GetReviewIssueId(),
		ReviewerUserId: req.GetReviewerUserId(),
		Pass:           req.GetPass(),
		Comment:        req.GetComment(),
	})

	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回结果
	rsp.CancelReview = commitRsp.GetCancelReview()
	return nil
}

// 验证request
func validateCommitTreatmentReviewResultRequest(req *pb.CommitTreatmentReviewResultRequest) error {
	if req.GetReviewIssueId() == "" {
		return gerr.New("issue_id should not be empty")
	}
	if req.GetReviewerUserId() == "" {
		return gerr.New("review_user_id should not be empty")
	}
	if !req.GetPass() && req.GetComment() == "" {
		return gerr.New("comment should not be empty")
	}
	return nil
}
