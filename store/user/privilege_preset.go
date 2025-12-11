package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

type PrivilegePreset struct {
	PresetID   string          `gorm:"primary_key;column:preset_id"`
	PresetName string          `gorm:"column:preset_name"`
	Policies   string          `gorm:"column:policies"`
	Remark     string          `gorm:"column:remark"`
	Rev        int32           `gorm:"column:rev"`
	CreatedAt  time.Time       // 创建时间
	UpdatedAt  time.Time       // 更新时间
	DeletedAt  *gorm.DeletedAt // 删除时间
}

func (p PrivilegePreset) TableName() string {
	return "privilege_preset"
}

// domain->db
func (p *PrivilegePreset) FromDomainPrivilegePreset(d domain.PrivilegePresetIntf) {
	if p == nil || d == nil {
		return
	}

	p.PresetID = d.GetPresetID()
	p.PresetName = d.GetPresetName()
	p.Policies = d.GetPolicies()
	p.Remark = d.GetRemark()
	p.Rev = d.GetRev()
}

// db->domain
func (p *PrivilegePreset) ToDomainPrivilegePreset() domain.PrivilegePresetIntf {
	if p == nil {
		return nil
	}

	pp := domain.PrivilegePreset{
		PresetID:   p.PresetID,
		PresetName: p.PresetName,
		Policies:   p.Policies,
		Remark:     p.Remark,
		Rev:        p.Rev,
	}
	return &pp
}
