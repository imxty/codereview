package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

type SubscriptionPeriod struct {
	SubscriptionPeriodID string          `gorm:"primary_key;column:subscription_period_id"`
	TenantID             string          `gorm:"column:tenant_id"`
	OrganizationID       string          `gorm:"column:organization_id"`
	OrganizationName     string          `gorm:"column:organization_name"`
	TenantName           string          `gorm:"column:tenant_name"`
	ExpiredTime          time.Time       `gorm:"column:expired_time"`
	UserID               string          `gorm:"column:user_id"`
	ContactName          string          `gorm:"column:contact_name"`
	ContactPhone         string          `gorm:"column:contact_phone"`
	Years                int32           `gorm:"column:years"`
	Rev                  int32           `gorm:"column:rev"`
	CreatedAt            time.Time       // 创建时间
	UpdatedAt            time.Time       // 更新时间
	DeletedAt            *gorm.DeletedAt // 删除时间
}

func (s SubscriptionPeriod) TableName() string {
	return "subscription_period"
}

// domain->db
func (s *SubscriptionPeriod) FromDomainSubscriptionPeriod(d domain.SubscriptionPeriodIntf) {
	if s == nil || d == nil {
		return
	}

	s.SubscriptionPeriodID = d.GetSubscriptionPeriodID()
	s.TenantID = d.GetTenantID()
	s.OrganizationID = d.GetOrganizationID()
	s.OrganizationName = d.GetOrganizationName()
	s.TenantName = d.GetTenantName()
	s.ExpiredTime = d.GetExpiredTime()
	s.UserID = d.GetUserID()
	s.ContactName = d.GetContactName()
	s.ContactPhone = d.GetContactPhone()
	s.Years = d.GetYears()
	s.Rev = d.GetRev()
	s.CreatedAt = d.GetCreatedAt()
}

// db->domain
func (s *SubscriptionPeriod) ToDomainSubscriptionPeriod() domain.SubscriptionPeriodIntf {
	if s == nil {
		return nil
	}

	p := domain.SubscriptionPeriod{
		SubscriptionPeriodID: s.SubscriptionPeriodID,
		TenantID:             s.TenantID,
		OrganizationID:       s.OrganizationID,
		OrganizationName:     s.OrganizationName,
		TenantName:           s.TenantName,
		ExpiredTime:          s.ExpiredTime,
		UserID:               s.UserID,
		ContactName:          s.ContactName,
		ContactPhone:         s.ContactPhone,
		Years:                s.Years,
		Rev:                  s.Rev,
		CreatedAt:            s.CreatedAt,
	}
	return &p
}
