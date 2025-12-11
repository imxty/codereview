package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

type OrganizationTenant struct {
	OrganizationTenantID string          `gorm:"primary_key;column:organization_tenant_id"`
	OrganizationID       string          `gorm:"column:organization_id"`
	TenantID             string          `gorm:"column:tenant_id"`
	IsActivated          bool            `gorm:"column:is_activated"`
	ReviewStatus         int32           `gorm:"column:review_status"`
	Rev                  int32           `gorm:"column:rev"`
	CreatedAt            time.Time       // 创建时间
	UpdatedAt            time.Time       // 更新时间
	DeletedAt            *gorm.DeletedAt // 删除时间
}

func (o OrganizationTenant) TableName() string {
	return "organization_tenant"
}

// domain->db
func (o *OrganizationTenant) FromDomainOrganizationTenant(d domain.OrganizationTenantIntf) {
	if o == nil || d == nil {
		return
	}

	o.OrganizationTenantID = d.GetOrganizationTenantID()
	o.OrganizationID = d.GetOrganizationID()
	o.TenantID = d.GetTenantID()
	o.IsActivated = d.GetIsActivated()
	o.ReviewStatus = d.GetReviewStatus()
	o.Rev = d.GetRev()
}

// db->domain
func (o *OrganizationTenant) ToDomainOrganizationTenant() domain.OrganizationTenantIntf {
	if o == nil {
		return nil
	}

	p := domain.OrganizationTenant{
		OrganizationTenantID: o.OrganizationTenantID,
		OrganizationID:       o.OrganizationID,
		TenantID:             o.TenantID,
		IsActivated:          o.IsActivated,
		ReviewStatus:         o.ReviewStatus,
		Rev:                  o.Rev,
	}
	return &p
}
