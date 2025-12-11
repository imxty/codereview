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

func (s *ProductAPIHandler) GetTreatmentRevByTreatmentRevId(ctx context.Context, req *pb.GetTreatmentRevByTreatmentRevIdRequest, rsp *pb.GetTreatmentRevByTreatmentRevIdResponse) error {
	// 验证request
	err := validateGetTreatmentRevByTreatmentRevIdRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 根据租户id获取方案
	dTreatmentRev, err := s.productStore.GetTreatmentRevByTreatmentRevId(ctx, req.GetTreatmentRevId())
	if err != nil {
		return err
	}
	if dTreatmentRev == nil {
		return errors.Error(ErrGetTreatmentFailed, "get treatment_rev failed")
	}

	// 获取方案配置
	pItems, err := s.productStore.ListTreatmentRevItems(ctx, req.GetTreatmentRevId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获得proto treatment_item
	itemsRiskyDisease, itemsDirtyDialectic, itemsPhysicalTherapy, itemsPhysical := transferToProtoTreatmentItems(pItems)

	// 通过revID 获取审核信息
	getRsp, err := s.reviewAPI.GetReviewResultByTargetId(ctx, &reviewpb.GetReviewResultByTargetIdRequest{
		TargetId: req.GetTreatmentRevId(),
	})
	if err != nil {
		return err
	}
	reviewResult := getRsp.GetResult()
	// 只有审核通过才去返回pass
	rsp.Treatment = &pb.Treatment{
		// 方案id
		TreatmentId: dTreatmentRev.GetTreatmentID(),
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
func validateGetTreatmentRevByTreatmentRevIdRequest(req *pb.GetTreatmentRevByTreatmentRevIdRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentRevId() == "" {
		return gerr.New("treatment_id should not be empty")
	}
	return nil
}
