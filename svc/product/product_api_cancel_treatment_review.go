package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrGetTreatmentFailed
	ErrGetTreatmentFailed = 5028
)

func (s *ProductAPIHandler) CancelTreatmentReview(ctx context.Context, req *pb.CancelTreatmentReviewRequest, rsp *pb.CancelTreatmentReviewResponse) error {
	// 验证request
	err := validateCancelTreatmentReviewRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取方案版本信息
	treatmentRev, err := s.productStore.GetTreatmentRev(ctx, req.GetOrganizationId(), req.GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if treatmentRev == nil {
		return errors.Errorf(ErrGetTreatmentFailed, "get treatment_rev by treatment_id[%s] failed", req.GetTreatmentId())
	}
	// 验证租户方案状态 审核中才能取消审核
	if treatmentRev.GetTreatmentStatus() != domain.TreatmentStatusUnderReview {
		return errors.Errorf(codes.InvalidOperation, "cancel treatment review failed for wrong treatment status[%s]", req.GetTreatmentId())
	}
	// 发送取消方案审核请求给review模块
	getCancelTreatmentReview, err := s.reviewAPI.CancelTreatmentReview(ctx, &reviewpb.CancelTreatmentReviewRequest{
		OrganizationId: req.GetOrganizationId(),
		TreatmentId:    treatmentRev.GetTreatmentRevID(),
	})
	if err != nil {
		return err
	}

	// 取消后审核终结 租户方案状态改为草稿状态
	err = s.productStore.ChangeTreatmentRevStatus(ctx, treatmentRev.GetTreatmentRevID(), treatmentRev.GetRev(), domain.TreatmentStatusDraft)
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "change treatment status failed[%s]", err.Error())
	}

	rsp.ReviewIsTerminated = getCancelTreatmentReview.GetReviewIsTerminated()
	return nil
}

// 验证request
func validateCancelTreatmentReviewRequest(req *pb.CancelTreatmentReviewRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment_id should not be empty")
	}
	return nil
}
