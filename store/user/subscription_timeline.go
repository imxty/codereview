package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

type SubscriptionTimeline struct {
	SubscriptionTimelineID string          `gorm:"primary_key;column:subscription_timeline_id"`
	TenantID               string          `gorm:"column:tenant_id"`
	StartTime              time.Time       `gorm:"column:start_time"`
	EndTime                time.Time       `gorm:"column:end_time"`
	Rev                    int32           `gorm:"column:rev"`
	Years                  int32           `gorm:"column:years"`
	CreatedAt              time.Time       // 创建时间
	UpdatedAt              time.Time       // 更新时间
	DeletedAt              *gorm.DeletedAt // 删除时间
}

func (s SubscriptionTimeline) TableName() string {
	return "subscription_timeline"
}

// domain->db
func (s *SubscriptionTimeline) FromDomainSubscriptionTimeline(d domain.SubscriptionTimelineIntf) {
	if s == nil || d == nil {
		return
	}

	s.SubscriptionTimelineID = d.GetSubscriptionTimelineID()
	s.TenantID = d.GetTenantID()
	s.StartTime = d.GetStartTime()
	s.EndTime = d.GetEndTime()
	s.Rev = d.GetRev()
	s.Years = d.GetYears()
}

// db->domain
func (s *SubscriptionTimeline) ToDomainSubscriptionTimeline() domain.SubscriptionTimelineIntf {
	if s == nil {
		return nil
	}

	p := domain.SubscriptionTimeline{
		SubscriptionTimelineID: s.SubscriptionTimelineID,
		TenantID:               s.TenantID,
		StartTime:              s.StartTime,
		EndTime:                s.EndTime,
		Rev:                    s.Rev,
		Years:                  s.Years,
	}
	return &p
}
