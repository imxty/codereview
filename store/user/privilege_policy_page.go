package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

type PrivilegePolicyPage struct {
	PrivilegeID string          `gorm:"primary_key;column:privilege_id"`
	PresetID    string          `gorm:"column:preset_id"`
	CreatedAt   time.Time       // 创建时间
	UpdatedAt   time.Time       // 更新时间
	DeletedAt   *gorm.DeletedAt // 删除时间
}

func (p PrivilegePolicyPage) TableName() string {
	return "privilege_policy_page"
}

// domain->db
func (p *PrivilegePolicyPage) FromDomainPrivilegePolicyPage(d domain.PrivilegePolicyPageIntf) {
	if p == nil || d == nil {
		return
	}

	p.PrivilegeID = d.GetPrivilegeID()
	p.PresetID = d.GetPresetID()
}

// db->domain
func (p *PrivilegePolicyPage) ToDomainPrivilegePolicyPage() domain.PrivilegePolicyPageIntf {
	if p == nil {
		return nil
	}

	pp := domain.PrivilegePolicyPage{
		PrivilegeID: p.PrivilegeID,
		PresetID:    p.PresetID,
	}
	return &pp
}
