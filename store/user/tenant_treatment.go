package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

type TenantTreatment struct {
	TenantTreatmentID string          `gorm:"primary_key;column:tenant_treatment_id"`
	TreatmentID       string          `gorm:"column:treatment_id"`
	TenantID          string          `gorm:"column:tenant_id"`
	Rev               int32           `gorm:"column:rev"`
	CreatedAt         time.Time       // 创建时间
	UpdatedAt         time.Time       // 更新时间
	DeletedAt         *gorm.DeletedAt // 删除时间
}

func (t TenantTreatment) TableName() string {
	return "tenant_treatment"
}

// domain->db
func (t *TenantTreatment) FromDomainTenantTreatment(d domain.TenantTreatmentIntf) {
	if t == nil || d == nil {
		return
	}

	t.TenantTreatmentID = d.GetTenantTreatmentID()
	t.TreatmentID = d.GetTreatmentID()
	t.TenantID = d.GetTenantID()
	t.Rev = d.GetRev()
}

// db->domain
func (t *TenantTreatment) ToDomainTenantTreatment() domain.TenantTreatmentIntf {
	if t == nil {
		return nil
	}

	p := domain.TenantTreatment{
		TenantTreatmentID: t.TenantTreatmentID,
		TreatmentID:       t.TreatmentID,
		TenantID:          t.TenantID,
		Rev:               t.Rev,
	}
	return &p
}
