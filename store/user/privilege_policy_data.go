package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

type PrivilegePolicyData struct {
	DataPrivilegeID string          `gorm:"primary_key;column:data_privilege_id"`
	PrivilegeID     string          `gorm:"column:privilege_id"`
	OrganizationID  string          `gorm:"column:organization_id"`
	Rev             int32           `gorm:"column:rev"`
	CreatedAt       time.Time       // 创建时间
	UpdatedAt       time.Time       // 更新时间
	DeletedAt       *gorm.DeletedAt // 删除时间
}

func (p PrivilegePolicyData) TableName() string {
	return "privilege_policy_data"
}

// domain->db
func (p *PrivilegePolicyData) FromDomainPrivilegePolicyData(d domain.PrivilegePolicyDataIntf) {
	if p == nil || d == nil {
		return
	}

	p.DataPrivilegeID = d.GetDataPrivilegeID()
	p.PrivilegeID = d.GetPrivilegeID()
	p.OrganizationID = d.GetOrganizationID()
	p.Rev = d.GetRev()
}

// db->domain
func (p *PrivilegePolicyData) ToDomainPrivilegePolicyData() domain.PrivilegePolicyDataIntf {
	if p == nil {
		return nil
	}

	pp := domain.PrivilegePolicyData{
		DataPrivilegeID: p.DataPrivilegeID,
		PrivilegeID:     p.PrivilegeID,
		OrganizationID:  p.OrganizationID,
		Rev:             p.Rev,
	}
	return &pp
}
