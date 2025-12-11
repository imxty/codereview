package product

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
)

type TreatmentProductAssociation struct {
	TreatmentItemID      string    `gorm:"column:treatment_item_id"`
	TreatmentRevID       string    `gorm:"column:treatment_rev_id"`
	TreatmentItemType    int32     `gorm:"column:treatment_item_type"`
	Symptom              string    `gorm:"column:symptom"`
	RecommendedProductID string    `gorm:"column:recommended_product_id"`
	OrganizationID       string    `gorm:"column:organization_id"`
	ProductImageUrl      string    `gorm:"column:product_image_url"`
	ProductName          string    `gorm:"column:product_name"`
	ProductType          int32     `gorm:"column:product_type"`
	ProductStatus        int32     `gorm:"column:product_status"`
	ProductIntroduction  string    `gorm:"column:product_introduction"`
	Remarks              string    `gorm:"column:remarks"`
	DrugName             string    `gorm:"column:drug_name"`
	DrugClassification   int32     `gorm:"column:drug_classification"`
	ApprovedNumber       string    `gorm:"column:approved_number"`
	DrugValidityPeriod   time.Time `gorm:"column:drug_validity_period"`
}

func (t TreatmentProductAssociation) TableName() string {
	return "treatment_product_association"
}

// domain->db
func (t *TreatmentProductAssociation) FromDomainTreatmentProductAssociation(d domain.TreatmentProductAssociationIntf) {
	if t == nil || d == nil {
		return
	}
	t.TreatmentItemID = d.GetTreatmentItemID()
	t.TreatmentRevID = d.GetTreatmentRevID()
	t.TreatmentItemType = d.GetTreatmentItemType()
	t.Symptom = d.GetSymptom()
	t.RecommendedProductID = d.GetRecommendedProductID()
	t.OrganizationID = d.GetOrganizationID()
	t.ProductImageUrl = d.GetProductImageUrl()
	t.ProductName = d.GetProductName()
	t.ProductType = d.GetProductType()
	t.ProductStatus = d.GetProductStatus()
	t.ProductIntroduction = d.GetProductIntroduction()
	t.Remarks = d.GetRemarks()
	t.DrugName = d.GetDrugName()
	t.DrugClassification = d.GetDrugClassification()
	t.ApprovedNumber = d.GetApprovedNumber()
	t.DrugValidityPeriod = d.GetDrugValidityPeriod()

}

// db->domain
func (t *TreatmentProductAssociation) ToDomainTreatmentProductAssociation() domain.TreatmentProductAssociationIntf {
	if t == nil {
		return nil
	}

	p := domain.TreatmentProductAssociation{
		TreatmentItemID:      t.TreatmentItemID,
		TreatmentRevID:       t.TreatmentRevID,
		TreatmentItemType:    t.TreatmentItemType,
		Symptom:              t.Symptom,
		RecommendedProductID: t.RecommendedProductID,
		OrganizationID:       t.OrganizationID,
		ProductImageUrl:      t.ProductImageUrl,
		ProductName:          t.ProductName,
		ProductType:          t.ProductType,
		ProductStatus:        t.ProductStatus,
		ProductIntroduction:  t.ProductIntroduction,
		Remarks:              t.Remarks,
		DrugName:             t.DrugName,
		DrugClassification:   t.DrugClassification,
		ApprovedNumber:       t.ApprovedNumber,
		DrugValidityPeriod:   t.DrugValidityPeriod,
	}
	return &p
}
