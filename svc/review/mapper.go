package review

import (
	gerr "errors"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toProtoNotification
func toProtoNotification(n domain.ReviewNotificationIntf) *pb.ReviewNotification {
	if n == nil {
		return nil
	}
	createdAt := timestamppb.New(n.GetCreatedAt())
	return &pb.ReviewNotification{
		// 通知ID
		NotificationId: n.GetReviewNotificationID(),
		// 通知标题
		NotificationTitle: n.GetNotificationTitle(),
		// 是否通过
		Pass: n.GetResult(),
		// 审核类型
		ReviewType: toProtoReviewType(n.GetTargetType()),
		// 审核意见
		ReviewComment: n.GetComment(),
		// 创建时间
		CreatedTime: createdAt,
	}
}

// toProtoReviewType
func toProtoReviewType(t int32) pb.ReviewType {
	switch t {
	case domain.ReviewTypeCertificate:
		return pb.ReviewType_REVIEW_TYPE_CERTIFICATE
	case domain.ReviewTypeProduct:
		return pb.ReviewType_REVIEW_TYPE_PRODUCT
	}
	return pb.ReviewType_REVIEW_TYPE_UNSET
}

// toReviewProtoProductFromProduct
func toReviewProtoProductFromProduct(t *productpb.Product) (*pb.Product, error) {
	pt, err := toReviewProtoProductTypeFromProduct(t.GetProductType())
	if err != nil {
		return nil, err
	}
	ps, err := toReviewProtoProductStatusFromProduct(t.GetProductStatus())
	if err != nil {
		return nil, err
	}
	return &pb.Product{
		// 商品id
		ProductId: t.GetProductId(),
		// 商品类型
		ProductType: pt,
		// 商品名称
		ProductName: t.GetProductName(),
		// 商品状态
		ProductStatus: ps,
		// 商品介绍
		Description: t.GetDescription(),
		// 商品图片
		ImageUrl: t.GetImageUrl(),
		// 商品备注
		Remarks: t.GetRemarks(),
		// 药品名称(商品类型为理疗服务方案不用传)
		DrugName: t.GetDrugName(),
		// 药品是非处方药 (商品类型为理疗服务方案不用传)
		IsOtc: t.GetIsOtc(),
		// 准字号 (商品类型为理疗服务方案不用传)
		ApprovedNumber: t.GetApprovedNumber(),
		// 药品准效期 (商品类型为理疗服务方案不用传)
		DrugValidityPeriod: t.GetDrugValidityPeriod(),
	}, nil
}

// toReviewProtoProductTypeFromProduct
func toReviewProtoProductTypeFromProduct(pt productpb.ProductType) (pb.ProductType, error) {
	switch pt {
	case productpb.ProductType_PRODUCT_TYPE_UNSET:
		return pb.ProductType_PRODUCT_TYPE_UNSET, gerr.New("unnset product type")
	case productpb.ProductType_PRODUCT_TYPE_SERVICE:
		return pb.ProductType_PRODUCT_TYPE_SERVICE, nil
	case productpb.ProductType_PRODUCT_TYPE_CPD:
		return pb.ProductType_PRODUCT_TYPE_CPD, nil
	case productpb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT:
		return pb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT, nil
	case productpb.ProductType_PRODUCT_TYPE_NUTRITION:
		return pb.ProductType_PRODUCT_TYPE_NUTRITION, nil
	default:
		return pb.ProductType_PRODUCT_TYPE_INVALID, gerr.New("invalid product type")
	}
}

// toReviewProtoProductStatusFromProduct
func toReviewProtoProductStatusFromProduct(ps productpb.ProductStatus) (pb.ProductStatus, error) {
	switch ps {
	case productpb.ProductStatus_PRODUCT_STATUS_UNSET:
		return pb.ProductStatus_PRODUCT_STATUS_UNSET, gerr.New("unset product status")
	case productpb.ProductStatus_PRODUCT_STATUS_USING:
		return pb.ProductStatus_PRODUCT_STATUS_USING, nil
	case productpb.ProductStatus_PRODUCT_STATUS_TO_BE_USED:
		return pb.ProductStatus_PRODUCT_STATUS_TO_BE_USED, nil
	case productpb.ProductStatus_PRODUCT_STATUS_UNUSED:
		return pb.ProductStatus_PRODUCT_STATUS_UNUSED, nil
	case productpb.ProductStatus_PRODUCT_STATUS_DELETED:
		return pb.ProductStatus_PRODUCT_STATUS_DELETED, nil
	}
	return pb.ProductStatus_PRODUCT_STATUS_INVALID, gerr.New("invalid product status")
}

// toDomainReviewStatus
func toDomainReviewStatus(s pb.ReviewStatus) (int32, error) {
	switch s {
	case pb.ReviewStatus_REVIEW_STATUS_PENDING_REVIEW:
		return domain.ReviewStatusPending, nil
	case pb.ReviewStatus_REVIEW_STATUS_REVIEWING:
		return domain.ReviewStatusReviewing, nil
	case pb.ReviewStatus_REVIEW_STATUS_SUCCESS:
		return domain.ReviewStatusReviewSuccess, nil
	case pb.ReviewStatus_REVIEW_STATUS_FAIL:
		return domain.ReviewStatusReviewFailed, nil
	}
	return 0, gerr.New("invalid review_status")
}

func toProtoReviewIssue(r domain.ReviewIssueWithResultIntf, nickname, organizationName string) *pb.ReviewIssue {
	if r == nil {
		return nil
	}
	submitTime := timestamppb.New(r.GetSubmitTime())
	return &pb.ReviewIssue{
		ReviewIssueId: r.GetReviewIssueID(),
		// 租户ID
		SubmitterTenantId:       r.GetSubmitterTenantID(),
		SubmitterOrganizationId: r.GetSubmitterOrganizationID(),
		// 提审时间
		SubmitTime: submitTime,
		// 审核人名称
		ReviewNickName: nickname,
		// 工单状态
		Status: toProtoReviewStatus(r.GetReviewStatus()),
		// 审核类型
		ReviewType: toProtoReviewType(r.GetTargetType()),
		// 是否被用户取消
		IsCancled: r.GetIsCancel(),
		// 是否通过
		IsPassed: r.GetResult(),
		// 未通过原因
		FailReason: r.GetComment(),
		// 审核人ID
		ReviewerUserId:   r.GetReviewerUserID(),
		OrganizationName: organizationName,
	}
}

// toProtoReviewStatus
func toProtoReviewStatus(s int32) pb.ReviewStatus {
	switch s {
	case domain.ReviewStatusPending:
		return pb.ReviewStatus_REVIEW_STATUS_PENDING_REVIEW
	case domain.ReviewStatusReviewing:
		return pb.ReviewStatus_REVIEW_STATUS_REVIEWING
	case domain.ReviewStatusReviewSuccess:
		return pb.ReviewStatus_REVIEW_STATUS_SUCCESS
	case domain.ReviewStatusReviewFailed:
		return pb.ReviewStatus_REVIEW_STATUS_FAIL
	}
	return pb.ReviewStatus_REVIEW_STATUS_UNSET
}

// toProtoTenantEntity
func toProtoTenantEntity(e *userpb.TenantEntity) *pb.TenantEntity {
	if e == nil {
		return nil
	}
	return &pb.TenantEntity{
		// 租户ID
		TenantId: e.GetTenantId(),
		// 商铺名称
		EntityName: e.GetName(),
		// 地址
		Address: toProtoAddress(e.GetAddress()),
		// 联系人姓名
		ContactName: e.GetContactName(),
		// 联系人电话
		ContactPhone: e.GetContactPhone(),
		// 营业执照
		BusinessLicenseUrl: e.GetBusinessLicenseUrl(),
		// 社会信用代码
		SocialCreditCode: e.GetSocialCreditCode(),
		LogoUrl:          e.GetLogoUrl(),
	}
}

// toProtoAddress
func toProtoAddress(a *userpb.Address) *pb.Address {
	if a == nil {
		return nil
	}
	return &pb.Address{
		// 省
		Province: a.GetProvince(),
		// 市
		City: a.GetCity(),
		// 区
		District: a.GetDistrict(),
		// 街道
		Street: a.GetStreet(),
	}
}
