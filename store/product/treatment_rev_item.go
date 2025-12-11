package product

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"gorm.io/gorm"
)

type TreatmentRevItem struct {
	TreatmentName        string          `gorm:"column:treatment_name"`
	TreatmentStatus      int32           `gorm:"column:treatment_status"`
	TreatmentItemID      string          `gorm:"column:treatment_item_id"`
	TreatmentRevID       string          `gorm:"column:treatment_rev_id"`
	TreatmentItemType    int32           `gorm:"column:treatment_item_type"`
	Symptom              string          `gorm:"column:symptom"`
	RecommendedProductID string          `gorm:"column:recommended_product_id"`
	CreatedAt            time.Time       // 创建时间
	UpdatedAt            time.Time       // 更新时间
	DeletedAt            *gorm.DeletedAt // 删除时间
}

// domain->db
func (t *TreatmentRevItem) FromDomainTreatmentRevItem(d domain.TreatmentRevItemIntf) {
	if t == nil || d == nil {
		return
	}

	t.TreatmentName = d.GetTreatmentName()
	t.TreatmentStatus = d.GetTreatmentStatus()
	t.TreatmentItemID = d.GetTreatmentItemID()
	t.TreatmentRevID = d.GetTreatmentRevID()
	t.TreatmentItemType = d.GetTreatmentItemType()
	t.Symptom = d.GetSymptom()
	t.RecommendedProductID = d.GetRecommendedProductID()
}

// db->domain
func (t *TreatmentRevItem) ToDomainTreatmentRevItem() domain.TreatmentRevItemIntf {
	if t == nil {
		return nil
	}

	p := domain.TreatmentRevItem{
		TreatmentName:        t.TreatmentName,
		TreatmentStatus:      t.TreatmentStatus,
		TreatmentItemID:      t.TreatmentItemID,
		TreatmentRevID:       t.TreatmentRevID,
		TreatmentItemType:    t.TreatmentItemType,
		Symptom:              t.Symptom,
		RecommendedProductID: t.RecommendedProductID,
	}
	return &p
}
