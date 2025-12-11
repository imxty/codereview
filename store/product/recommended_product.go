package product

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"gorm.io/gorm"
)

type RecommendedProduct struct {
	RecommendedProductID string          `gorm:"primary_key;column:recommended_product_id"`
	OrganizationID       string          `gorm:"column:organization_id"`
	ProductImageUrl      string          `gorm:"column:product_image_url"`
	ProductName          string          `gorm:"column:product_name"`
	ProductType          int32           `gorm:"column:product_type"`
	ProductStatus        int32           `gorm:"column:product_status"`
	ProductIntroduction  string          `gorm:"column:product_introduction"`
	Remarks              string          `gorm:"column:remarks"`
	DrugName             string          `gorm:"column:drug_name"`
	DrugClassification   int32           `gorm:"column:drug_classification"`
	ApprovedNumber       string          `gorm:"column:approved_number"`
	DrugValidityPeriod   time.Time       `gorm:"column:drug_validity_period"`
	Symptoms             string          `gorm:"column:symptoms"`
	Rev                  int32           `gorm:"column:rev"`
	CreatedAt            time.Time       // 创建时间
	UpdatedAt            time.Time       // 更新时间
	DeletedAt            *gorm.DeletedAt // 删除时间
}

func (r RecommendedProduct) TableName() string {
	return "recommended_product"
}

// domain->db
func (r *RecommendedProduct) FromDomainRecommendedProduct(d domain.RecommendedProductIntf) {
	if r == nil || d == nil {
		return
	}

	r.RecommendedProductID = d.GetRecommendedProductID()
	r.OrganizationID = d.GetOrganizationID()
	r.ProductImageUrl = d.GetProductImageUrl()
	r.ProductName = d.GetProductName()
	r.ProductType = d.GetProductType()
	r.ProductStatus = d.GetProductStatus()
	r.ProductIntroduction = d.GetProductIntroduction()
	r.Remarks = d.GetRemarks()
	r.DrugName = d.GetDrugName()
	r.DrugClassification = d.GetDrugClassification()
	r.ApprovedNumber = d.GetApprovedNumber()
	r.DrugValidityPeriod = d.GetDrugValidityPeriod()
	r.Symptoms = d.GetSymptoms()
	r.Rev = d.GetRev()
}

// db->domain
func (r *RecommendedProduct) ToDomainRecommendedProduct() domain.RecommendedProductIntf {
	if r == nil {
		return nil
	}

	p := domain.RecommendedProduct{
		RecommendedProductID: r.RecommendedProductID,
		OrganizationID:       r.OrganizationID,
		ProductImageUrl:      r.ProductImageUrl,
		ProductName:          r.ProductName,
		ProductType:          r.ProductType,
		ProductStatus:        r.ProductStatus,
		ProductIntroduction:  r.ProductIntroduction,
		Remarks:              r.Remarks,
		DrugName:             r.DrugName,
		DrugClassification:   r.DrugClassification,
		ApprovedNumber:       r.ApprovedNumber,
		DrugValidityPeriod:   r.DrugValidityPeriod,
		Symptoms:             r.Symptoms,
		Rev:                  r.Rev,
	}
	return &p
}
