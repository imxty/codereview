package review

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 取消商品推荐方案审核请求
func (s *ReviewAPIHandler) CancelTreatmentReview(ctx context.Context, req *pb.CancelTreatmentReviewRequest, rsp *pb.CancelTreatmentReviewResponse) error {
	err := validateCancelTreatmentReviewRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取review issue信息(treatment Rev id)
	reviewIssue, err := s.reviewStore.GetLatestReviewIssueByTargetId(ctx, req.GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if reviewIssue == nil {
		return errors.Errorf(ErrReviewIssueNotFound, "get review issue by treatment id failed[%s]", req.GetTreatmentId())
	}

	// 验证获取审核资质状态 不是审核中或待审核状态 报错
	if reviewIssue.GetReviewStatus() != domain.ReviewStatusPending {
		return errors.Errorf(codes.InvalidRequest, "cancel treatment review failed for wrong review status[%s]", reviewIssue.GetReviewIssueID())
	}

	// 修改review issue状态
	err = s.reviewStore.CancelTreatmentReviewIssue(ctx, reviewIssue.GetReviewIssueID(), reviewIssue.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证request
func validateCancelTreatmentReviewRequest(req *pb.CancelTreatmentReviewRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment id should not be empty")
	}
	return nil
}
