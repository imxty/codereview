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

// 检查租户提审资格 当前租户的review情况
// 无法提交审核的情况:
func (s *ProductAPIHandler) SubmitTreatmentToReview(ctx context.Context, req *pb.SubmitTreatmentToReviewRequest, rsp *pb.SubmitTreatmentToReviewResponse) error {
	err := validateSubmitTreatmentToReviewRequest(req)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 提审
	// 获取方案信息
	treatment, err := s.productStore.GetTreatmentRev(ctx, req.GetOrganizationId(), req.GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if treatment == nil {
		return errors.Error(ErrGetTreatmentFailed, "get treatment failed")
	}
	// 检验是否处于草稿状态
	if treatment.GetTreatmentStatus() != int32(pb.TreatmentStatus_TREATMENT_STATUS_DRAFT) {
		return errors.Errorf(codes.InvalidRequest, "treatment[%s] status should be draft when submitting treatment review", req.GetTreatmentId())
	}
	// 提交方案id给review模块进行操作
	// 保存的是treatmentRevId
	_, err = s.reviewAPI.SubmitTreatmentToReview(ctx, &reviewpb.SubmitTreatmentToReviewRequest{
		OrganizationId: req.GetOrganizationId(),
		TreatmentRevId: treatment.GetTreatmentRevID(),
	})
	if err != nil {
		return err
	}

	// 提审即变为审核中
	err = s.productStore.ChangeTreatmentRevStatus(ctx, treatment.GetTreatmentRevID(), treatment.GetRev(), domain.TreatmentStatusUnderReview)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证request
func validateSubmitTreatmentToReviewRequest(req *pb.SubmitTreatmentToReviewRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment id should not be empty")
	}
	return nil
}
