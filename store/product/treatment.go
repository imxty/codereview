package product

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"gorm.io/gorm"
)

type Treatment struct {
	TreatmentID    string          `gorm:"primary_key;column:treatment_id"`
	OrganizationID string          `gorm:"column:organization_id"`
	LatestRev      string          `gorm:"column:latest_rev"`
	Rev            int32           `gorm:"column:rev"`
	CreatedAt      time.Time       // 创建时间
	UpdatedAt      time.Time       // 更新时间
	DeletedAt      *gorm.DeletedAt // 删除时间
}

func (t Treatment) TableName() string {
	return "treatment"
}

// domain->db
func (t *Treatment) FromDomainTreatment(d domain.TreatmentIntf) {
	if t == nil || d == nil {
		return
	}

	t.TreatmentID = d.GetTreatmentID()
	t.OrganizationID = d.GetOrganizationID()
	t.LatestRev = d.GetLatestRev()
	t.Rev = d.GetRev()

}

// db->domain
func (t *Treatment) ToDomainTreatment() domain.TreatmentIntf {
	if t == nil {
		return nil
	}

	p := domain.Treatment{
		TreatmentID:    t.TreatmentID,
		OrganizationID: t.OrganizationID,
		LatestRev:      t.LatestRev,
		Rev:            t.Rev,
	}
	return &p
}
