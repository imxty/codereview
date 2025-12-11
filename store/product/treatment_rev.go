package product

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"gorm.io/gorm"
)

type TreatmentRev struct {
	TreatmentRevID  string          `gorm:"primary_key;column:treatment_rev_id"`
	TreatmentID     string          `gorm:"column:treatment_id"`
	TreatmentName   string          `gorm:"column:treatment_name"`
	Remarks         string          `gorm:"column:remarks"`
	OrganizationID  string          `gorm:"column:organization_id"`
	IsPublished     bool            `gorm:"column:is_published"`
	TreatmentStatus int32           `gorm:"column:treatment_status"`
	Rev             int32           `gorm:"column:rev"`
	CreatedAt       time.Time       // 创建时间
	UpdatedAt       time.Time       // 更新时间
	DeletedAt       *gorm.DeletedAt // 删除时间
}

func (t TreatmentRev) TableName() string {
	return "treatment_rev"
}

// domain->db
func (t *TreatmentRev) FromDomainTreatmentRev(d domain.TreatmentRevIntf) {
	if t == nil || d == nil {
		return
	}

	t.TreatmentRevID = d.GetTreatmentRevID()
	t.TreatmentID = d.GetTreatmentID()
	t.TreatmentName = d.GetTreatmentName()
	t.Remarks = d.GetRemarks()
	t.OrganizationID = d.GetOrganizationID()
	t.IsPublished = d.GetIsPublished()
	t.TreatmentStatus = d.GetTreatmentStatus()
	t.Rev = d.GetRev()
	t.CreatedAt = d.GetCreatedAt()
}

// db->domain
func (t *TreatmentRev) ToDomainTreatmentRev() domain.TreatmentRevIntf {
	if t == nil {
		return nil
	}

	p := domain.TreatmentRev{
		TreatmentRevID:  t.TreatmentRevID,
		TreatmentID:     t.TreatmentID,
		TreatmentName:   t.TreatmentName,
		Remarks:         t.Remarks,
		OrganizationID:  t.OrganizationID,
		IsPublished:     t.IsPublished,
		TreatmentStatus: t.TreatmentStatus,
		Rev:             t.Rev,
		CreatedAt:       t.CreatedAt,
	}
	return &p
}
