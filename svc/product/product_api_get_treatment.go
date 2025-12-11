package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ProductAPIHandler) GetTreatment(ctx context.Context, req *pb.GetTreatmentRequest, rsp *pb.GetTreatmentResponse) error {
	err := validateGetTreatmentRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 查询treatment
	treatment, err := s.productStore.GetTreatment(ctx, req.GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if treatment == nil {
		return errors.Errorf(ErrGetTreatmentFailed, "treatment[%s] not found", req.GetTreatmentId())
	}
	if treatment.GetOrganizationID() != req.GetOrganizationId() {
		return errors.Errorf(codes.InvalidRequest, "treatment[%s] not belong to organizationId[%s]", req.GetTreatmentId(), req.GetOrganizationId())
	}

	// 查询treatment相关的信息，审核信息，方案信息商品信息
	// 获取方案信息
	dTreatmentRev, err := s.productStore.GetTreatmentRev(ctx, req.GetOrganizationId(), req.GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if dTreatmentRev == nil {
		return errors.Errorf(ErrGetTreatmentFailed, "get treatment_rev by treatmentID[%s] failed", req.GetTreatmentId())
	}
	// 获取方案配置
	pItems, err := s.productStore.ListTreatmentItems(ctx, req.GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获得proto treatment_item
	itemsRiskyDisease, itemsDirtyDialectic, itemsPhysicalTherapy, itemsPhysical := transferToProtoTreatmentItems(pItems)

	// 通过revID 获取审核信息
	getRsp, err := s.reviewAPI.GetReviewResultByTargetId(ctx, &reviewpb.GetReviewResultByTargetIdRequest{
		TargetId: dTreatmentRev.GetTreatmentRevID(),
	})
	if err != nil {
		return err
	}
	reviewResult := getRsp.GetResult()
	// 只有审核通过才去返回pass
	rsp.Treatment = &pb.Treatment{
		// 方案id
		TreatmentId: req.GetTreatmentId(),
		// 方案名称
		TreatmentName: dTreatmentRev.GetTreatmentName(),
		// 审核是否通过
		ReviewPass: reviewResult.GetReiewStatus() == reviewpb.ReviewStatus_REVIEW_STATUS_SUCCESS,
		// 审核意见
		ReviewComment: reviewResult.GetFailReason(),
		// 是否发布
		IsPublished: dTreatmentRev.GetIsPublished(),
		// 方案状态
		TreatmentStatus: toProtoTreatmentStatus(dTreatmentRev.GetTreatmentStatus()),
		// 风险疾病方案配置
		TreatmentItemsRiskyDisease: itemsRiskyDisease,
		// 脏腑辩证方案配置
		TreatmentItemsDirtyDialectic: itemsDirtyDialectic,
		// 理疗方案配置
		TreatmentItemsPhysicalTherapy: itemsPhysicalTherapy,
		// 创建时间
		CreatedTime: timestamppb.New(dTreatmentRev.GetCreatedAt()),
		// 体质方案
		TreatmentItemsPhysicalDialectics: itemsPhysical,
	}
	return nil
}

// 验证request
func validateGetTreatmentRequest(req *pb.GetTreatmentRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organizationId should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment id should not be empty")
	}
	return nil
}
