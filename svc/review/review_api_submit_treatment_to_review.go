package review

import (
	"context"
	gerr "errors"
	"time"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

var (
	// 初始化时间
	InitTime = time.Date(1970, 1, 1, 0, 0, 1, 0, time.UTC)
)

func (s *ReviewAPIHandler) SubmitTreatmentToReview(ctx context.Context, req *pb.SubmitTreatmentToReviewRequest, rsp *pb.SubmitTreatmentToReviewResponse) error {
	err := validateSubmitTreatmentToReviewRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	now := time.Now().UTC()
	// 创建新的issue
	issue := &domain.ReviewIssue{
		ReviewIssueID:           xid.New().String(),
		SubmitterOrganizationID: req.GetOrganizationId(),
		SubmitTime:              now,
		AcceptanceTime:          InitTime,
		ReviewStatus:            domain.ReviewStatusPending,
		TargetType:              domain.ReviewTypeProduct,
		IsCancel:                false,
		TargetRevID:             req.GetTreatmentRevId(),
	}
	err = s.reviewStore.CreateReviewIssue(ctx, issue)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 更新方案状态
	// 改变租户方案状态为审核中
	_, err = s.productAPI.UpdateTreatmentStatus(ctx, &productpb.UpdateTreatmentStatusRequest{
		OrganizationId:  req.GetOrganizationId(),
		TreatmentRevId:  req.GetTreatmentRevId(),
		TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_UNDER_REVIEW,
	})
	if err != nil {
		return err
	}
	return nil
}

// 验证request
func validateSubmitTreatmentToReviewRequest(req *pb.SubmitTreatmentToReviewRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentRevId() == "" {
		return gerr.New("treatment rev id should not be empty")
	}
	return nil
}
