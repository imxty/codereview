package product

import (
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// proto->domain toDomainDrugClassification
var toDomainDrugClassification = map[bool]int32{
	false: 0,
	true:  1,
}

// proto->domain  toDomainProductStatus
func toDomainProductStatus(ps pb.ProductStatus) (int32, error) {
	switch ps {
	case pb.ProductStatus_PRODUCT_STATUS_UNSET:
		return domain.ProductStatusUnset, gerr.New("unset product status")
	case pb.ProductStatus_PRODUCT_STATUS_USING:
		return domain.ProductStatusUsing, nil
	case pb.ProductStatus_PRODUCT_STATUS_TO_BE_USED:
		return domain.ProductStatusToBeUsed, nil
	case pb.ProductStatus_PRODUCT_STATUS_UNUSED:
		return domain.ProductStatusUnused, nil
	case pb.ProductStatus_PRODUCT_STATUS_DELETED:
		return domain.ProductStatusDeleted, nil
	}
	return domain.ProductStatusInvalid, gerr.New("invalid product status")
}

// proto->domain toDomainProductType
func toDomainProductType(pt pb.ProductType) (int32, error) {
	switch pt {
	case pb.ProductType_PRODUCT_TYPE_UNSET:
		return domain.ProductTypeUnset, gerr.New("unnset product type")
	case pb.ProductType_PRODUCT_TYPE_SERVICE:
		return domain.ProductTypeService, nil
	case pb.ProductType_PRODUCT_TYPE_CPD:
		return domain.ProductTypeCpd, nil
	case pb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT:
		return domain.ProductTypeHealthCareProduct, nil
	case pb.ProductType_PRODUCT_TYPE_NUTRITION:
		return domain.ProductTypeNutrition, nil
	default:
		return domain.ProductTypeInvalid, gerr.New("invalid product type")
	}
}

// domain -> proto toProtoProductFromTreatmentProductAssociation
func toProtoProductFromTreatmentProductAssociation(t domain.TreatmentProductAssociationIntf) *pb.Product {
	if t == nil {
		return nil
	}
	drugValidityPeriod := timestamppb.New(t.GetDrugValidityPeriod())
	return &pb.Product{
		// 商品id
		ProductId: t.GetRecommendedProductID(),
		// 商品类型
		ProductType: toProtoProductType(t.GetProductType()),
		// 商品名称
		ProductName: t.GetProductName(),
		// 商品状态
		ProductStatus: toProtoProductStatus(t.GetProductStatus()),
		// 商品介绍
		Description: t.GetProductIntroduction(),
		// 商品图片
		ImageUrl: t.GetProductImageUrl(),
		// 商品备注
		Remarks: t.GetRemarks(),
		// 药品名称(商品类型为理疗服务方案不用传)
		DrugName: t.GetDrugName(),
		// 药品是非处方药 (商品类型为理疗服务方案不用传)
		IsOtc: t.GetDrugClassification() == domain.DrugClassificationIsOtc,
		// 准字号 (商品类型为理疗服务方案不用传)
		ApprovedNumber: t.GetApprovedNumber(),
		// 药品准效期 (商品类型为理疗服务方案不用传)
		DrugValidityPeriod: drugValidityPeriod,
	}
}

// domain -> proto toProtoProductFromTreatmentProductAssociationWithSymptom
func toProtoProductFromTreatmentProductAssociationWithSymptom(t domain.TreatmentProductAssociationIntf, symptomKey string) *pb.Product {
	if t == nil {
		return nil
	}
	drugValidityPeriod := timestamppb.New(t.GetDrugValidityPeriod())
	return &pb.Product{
		// 商品id
		ProductId: t.GetRecommendedProductID(),
		// 商品类型
		ProductType: toProtoProductType(t.GetProductType()),
		// 商品名称
		ProductName: t.GetProductName(),
		// 商品状态
		ProductStatus: toProtoProductStatus(t.GetProductStatus()),
		// 商品介绍
		Description: t.GetProductIntroduction(),
		// 商品图片
		ImageUrl: t.GetProductImageUrl(),
		// 商品备注
		Remarks: t.GetRemarks(),
		// 药品名称(商品类型为理疗服务方案不用传)
		DrugName: t.GetDrugName(),
		// 药品是非处方药 (商品类型为理疗服务方案不用传)
		IsOtc: t.GetDrugClassification() == domain.DrugClassificationIsOtc,
		// 准字号 (商品类型为理疗服务方案不用传)
		ApprovedNumber: t.GetApprovedNumber(),
		// 药品准效期 (商品类型为理疗服务方案不用传)
		DrugValidityPeriod: drugValidityPeriod,
		// 对应的症候Key
		SymptomKeys: []string{symptomKey},
	}
}

// domain->proto toProtoProductType
func toProtoProductType(productType int32) pb.ProductType {
	switch productType {
	case domain.ProductTypeUnset:
		return pb.ProductType_PRODUCT_TYPE_UNSET
	case domain.ProductTypeService:
		return pb.ProductType_PRODUCT_TYPE_SERVICE
	case domain.ProductTypeCpd:
		return pb.ProductType_PRODUCT_TYPE_CPD
	case domain.ProductTypeHealthCareProduct:
		return pb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT
	case domain.ProductTypeNutrition:
		return pb.ProductType_PRODUCT_TYPE_NUTRITION
	}
	return pb.ProductType_PRODUCT_TYPE_INVALID
}

// domain->proto toProtoProductStatus
func toProtoProductStatus(protoStatus int32) pb.ProductStatus {
	switch protoStatus {
	case domain.ProductStatusUnset:
		return pb.ProductStatus_PRODUCT_STATUS_UNSET
	case domain.ProductStatusUsing:
		return pb.ProductStatus_PRODUCT_STATUS_USING
	case domain.ProductStatusToBeUsed:
		return pb.ProductStatus_PRODUCT_STATUS_TO_BE_USED
	case domain.ProductStatusUnused:
		return pb.ProductStatus_PRODUCT_STATUS_UNUSED
	case domain.ProductStatusDeleted:
		return pb.ProductStatus_PRODUCT_STATUS_DELETED
	}
	return pb.ProductStatus_PRODUCT_STATUS_INVALID
}

// 转换得到proto treatment
func transferToProtoTreatment(dTreatmentRev domain.TreatmentRevIntf, rd, dd, pt, pd []*pb.TreatmentItem) *pb.Treatment {
	if dTreatmentRev == nil {
		return nil
	}
	dCreatedAt := timestamppb.New(dTreatmentRev.GetCreatedAt())
	return &pb.Treatment{
		TreatmentId:   dTreatmentRev.GetTreatmentID(),
		TreatmentName: dTreatmentRev.GetTreatmentName(),
		// OrganizationID:                dTreatmentRev.GetOrganizationID(),
		IsPublished:                      dTreatmentRev.GetIsPublished() == domain.TreatmentRevIsPublished,
		TreatmentStatus:                  toProtoTreatmentStatus(dTreatmentRev.GetTreatmentStatus()),
		TreatmentItemsRiskyDisease:       rd,
		TreatmentItemsDirtyDialectic:     dd,
		TreatmentItemsPhysicalTherapy:    pt,
		TreatmentItemsPhysicalDialectics: pd,
		CreatedTime:                      dCreatedAt,
	}
}

// domain->proto toProtoTreatmentStatus
func toProtoTreatmentStatus(status int32) pb.TreatmentStatus {
	switch status {
	case domain.TreatmentStatusUnset:
		return pb.TreatmentStatus_TREATMENT_STATUS_UNSET
	case domain.TreatmentStatusDraft:
		return pb.TreatmentStatus_TREATMENT_STATUS_DRAFT
	case domain.TreatmentStatusUnderReview:
		return pb.TreatmentStatus_TREATMENT_STATUS_UNDER_REVIEW
	case domain.TreatmentStatusApproved:
		return pb.TreatmentStatus_TREATMENT_STATUS_APPROVED
	case domain.TreatmentStatusUnapproved:
		return pb.TreatmentStatus_TREATMENT_STATUS_UNAPPROVED
	default:
		return pb.TreatmentStatus_TREATMENT_STATUS_INVALID
	}

}

// domain -> proto toProtoProductFromRecommendedProduct
func toProtoProductFromRecommendedProduct(t domain.RecommendedProductIntf) *pb.Product {
	if t == nil {
		return nil
	}
	drugValidityPeriod := timestamppb.New(t.GetDrugValidityPeriod())
	return &pb.Product{
		// 商品id
		ProductId: t.GetRecommendedProductID(),
		// 商品类型
		ProductType: toProtoProductType(t.GetProductType()),
		// 商品名称
		ProductName: t.GetProductName(),
		// 商品状态
		ProductStatus: toProtoProductStatus(t.GetProductStatus()),
		// 商品介绍
		Description: t.GetProductIntroduction(),
		// 商品图片
		ImageUrl: t.GetProductImageUrl(),
		// 商品备注
		Remarks: t.GetRemarks(),
		// 药品名称(商品类型为理疗服务方案不用传)
		DrugName: t.GetDrugName(),
		// 药品是非处方药 (商品类型为理疗服务方案不用传)
		IsOtc: t.GetDrugClassification() == domain.DrugClassificationIsOtc,
		// 准字号 (商品类型为理疗服务方案不用传)
		ApprovedNumber: t.GetApprovedNumber(),
		// 药品准效期 (商品类型为理疗服务方案不用传)
		DrugValidityPeriod: drugValidityPeriod,
		SymptomKeys:        toSymptomsKeys(t.GetSymptoms()),
	}
}

// proto->domain toDomainTreatmentStatus
func toDomainTreatmentStatus(status pb.TreatmentStatus) (int32, error) {
	switch status {
	case pb.TreatmentStatus_TREATMENT_STATUS_UNSET:
		return domain.TreatmentStatusUnset, gerr.New("unset treatment status")
	case pb.TreatmentStatus_TREATMENT_STATUS_DRAFT:
		return domain.TreatmentStatusDraft, nil
	case pb.TreatmentStatus_TREATMENT_STATUS_UNDER_REVIEW:
		return domain.TreatmentStatusUnderReview, nil
	case pb.TreatmentStatus_TREATMENT_STATUS_APPROVED:
		return domain.TreatmentStatusApproved, nil
	case pb.TreatmentStatus_TREATMENT_STATUS_UNAPPROVED:
		return domain.TreatmentStatusUnapproved, nil
	default:
		return domain.TreatmentStatusInvalid, gerr.New("invalid treatment status")
	}
}

// 转换为proto Treatment  不包括方案配置
func transferToProtoTreatmentInfo(t domain.TreatmentRevIntf, result bool, comment string) *pb.Treatment {
	dCreatedAt := timestamppb.New(t.GetCreatedAt())
	return &pb.Treatment{
		// 方案id
		TreatmentId: t.GetTreatmentID(),
		// 方案名称
		TreatmentName: t.GetTreatmentName(),
		// 审核是否通过
		ReviewPass: result,
		// 审核评论
		ReviewComment: comment,
		// 是否发布
		IsPublished: t.GetIsPublished() == domain.TreatmentRevIsPublished,
		// 方案状态
		TreatmentStatus: toProtoTreatmentStatus(t.GetTreatmentStatus()),
		// 创建时间
		CreatedTime: dCreatedAt,
	}
}
