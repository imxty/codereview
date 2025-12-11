package domain

import (
	"time"
)

const (
	// 无效的商品类型
	ProductTypeInvalid int32 = 0
	// 未设置的商品类型
	ProductTypeUnset int32 = 1
	// 理疗服务
	ProductTypeService int32 = 2
	// 中成药
	ProductTypeCpd int32 = 3
	// 保健品
	ProductTypeHealthCareProduct int32 = 4
	// 营养食品
	ProductTypeNutrition int32 = 5
)
const (
	// 无效的商品状态
	ProductStatusInvalid int32 = 0
	// 未设置的商品状态
	ProductStatusUnset int32 = 1
	// 使用中
	ProductStatusUsing int32 = 2
	// 待使用
	ProductStatusToBeUsed int32 = 3
	// 未使用
	ProductStatusUnused int32 = 4
	// 已删除
	ProductStatusDeleted int32 = 5
)
const (
	// 处方药
	DrugClassificationIsNotOtc = 0
	// 非处方药
	DrugClassificationIsOtc = 1
)

type RecommendedProduct struct {
	RecommendedProductID string
	OrganizationID       string
	ProductImageUrl      string
	ProductName          string
	ProductType          int32
	ProductStatus        int32
	ProductIntroduction  string
	Remarks              string
	DrugName             string
	DrugClassification   int32
	ApprovedNumber       string
	DrugValidityPeriod   time.Time
	Symptoms             string
	Rev                  int32
}

type RecommendedProductMapper interface {
	ToDomainRecommendedProductMapper
	FromDomainRecommendedProductMapper
}

type ToDomainRecommendedProductMapper interface {
	ToDomainRecommendedProduct() RecommendedProductIntf
}

type FromDomainRecommendedProductMapper interface {
	FromDomainRecommendedProduct(RecommendedProductIntf)
}

type RecommendedProductIntf interface {
	GetRecommendedProductID() string
	GetOrganizationID() string
	GetProductImageUrl() string
	GetProductName() string
	GetProductType() int32
	GetProductStatus() int32
	GetProductIntroduction() string
	GetRemarks() string
	GetDrugName() string
	GetDrugClassification() int32
	GetApprovedNumber() string
	GetDrugValidityPeriod() time.Time
	GetSymptoms() string
	GetRev() int32
}

var _ RecommendedProductIntf = (*RecommendedProduct)(nil)

func (r *RecommendedProduct) GetRecommendedProductID() string {
	if r == nil {
		return ""
	}
	return r.RecommendedProductID
}

func (r *RecommendedProduct) GetOrganizationID() string {
	if r == nil {
		return ""
	}
	return r.OrganizationID
}

func (r *RecommendedProduct) GetProductImageUrl() string {
	if r == nil {
		return ""
	}
	return r.ProductImageUrl
}

func (r *RecommendedProduct) GetProductName() string {
	if r == nil {
		return ""
	}
	return r.ProductName
}

func (r *RecommendedProduct) GetProductType() int32 {
	if r == nil {
		return 0
	}
	return r.ProductType
}

func (r *RecommendedProduct) GetProductStatus() int32 {
	if r == nil {
		return 0
	}
	return r.ProductStatus
}

func (r *RecommendedProduct) GetProductIntroduction() string {
	if r == nil {
		return ""
	}
	return r.ProductIntroduction
}

func (r *RecommendedProduct) GetRemarks() string {
	if r == nil {
		return ""
	}
	return r.Remarks
}

func (r *RecommendedProduct) GetDrugName() string {
	if r == nil {
		return ""
	}
	return r.DrugName
}

func (r *RecommendedProduct) GetDrugClassification() int32 {
	if r == nil {
		return 0
	}
	return r.DrugClassification
}

func (r *RecommendedProduct) GetApprovedNumber() string {
	if r == nil {
		return ""
	}
	return r.ApprovedNumber
}

func (r *RecommendedProduct) GetDrugValidityPeriod() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.DrugValidityPeriod
}

func (r *RecommendedProduct) GetSymptoms() string {
	if r == nil {
		return ""
	}
	return r.Symptoms
}

func (r *RecommendedProduct) GetRev() int32 {
	if r == nil {
		return 0
	}
	return r.Rev
}
