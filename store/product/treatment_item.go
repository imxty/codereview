package product

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"gorm.io/gorm"
)

type TreatmentItem struct {
	TreatmentItemID      string          `gorm:"primary_key;column:treatment_item_id"`
	TreatmentRevID       string          `gorm:"column:treatment_rev_id"`
	TreatmentItemType    int32           `gorm:"column:treatment_item_type"`
	Symptom              string          `gorm:"column:symptom"`
	RecommendedProductID string          `gorm:"column:recommended_product_id"`
	CreatedAt            time.Time       // 创建时间
	UpdatedAt            time.Time       // 更新时间
	DeletedAt            *gorm.DeletedAt // 删除时间
}

func (t TreatmentItem) TableName() string {
	return "treatment_item"
}

// domain->db
func (t *TreatmentItem) FromDomainTreatmentItem(d domain.TreatmentItemIntf) {
	if t == nil || d == nil {
		return
	}

	t.TreatmentItemID = d.GetTreatmentItemID()
	t.TreatmentRevID = d.GetTreatmentRevID()
	t.TreatmentItemType = d.GetTreatmentItemType()
	t.Symptom = d.GetSymptom()
	t.RecommendedProductID = d.GetRecommendedProductID()
}

// db->domain
func (t *TreatmentItem) ToDomainTreatmentItem() domain.TreatmentItemIntf {
	if t == nil {
		return nil
	}

	p := domain.TreatmentItem{
		TreatmentItemID:      t.TreatmentItemID,
		TreatmentRevID:       t.TreatmentRevID,
		TreatmentItemType:    t.TreatmentItemType,
		Symptom:              t.Symptom,
		RecommendedProductID: t.RecommendedProductID,
	}
	return &p
}
