package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

type PrivilegePolicy struct {
	PrivilegeID   string          `gorm:"primary_key;column:privilege_id"`
	PrivilegeName string          `gorm:"column:privilege_name"`
	Remark        string          `gorm:"column:remark"`
	Rev           int32           `gorm:"column:rev"`
	CreatedAt     time.Time       // 创建时间
	UpdatedAt     time.Time       // 更新时间
	DeletedAt     *gorm.DeletedAt // 删除时间
}

func (p PrivilegePolicy) TableName() string {
	return "privilege_policy"
}

// domain->db
func (p *PrivilegePolicy) FromDomainPrivilegePolicy(d domain.PrivilegePolicyIntf) {
	if p == nil || d == nil {
		return
	}

	p.PrivilegeID = d.GetPrivilegeID()
	p.PrivilegeName = d.GetPrivilegeName()
	p.Remark = d.GetRemark()
	p.Rev = d.GetRev()
}

// db->domain
func (p *PrivilegePolicy) ToDomainPrivilegePolicy() domain.PrivilegePolicyIntf {
	if p == nil {
		return nil
	}

	pp := domain.PrivilegePolicy{
		PrivilegeID:   p.PrivilegeID,
		PrivilegeName: p.PrivilegeName,
		Remark:        p.Remark,
		Rev:           p.Rev,
	}
	return &pp
}
