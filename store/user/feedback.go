package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"gorm.io/gorm"
)

// Feedback 反馈意见
type Feedback struct {
	FeedbackID string    `gorm:"primary_key;column:feedback_id"`
	TenantID   string    `gorm:"column:tenant_id"`
	Phone      string    `gorm:"column:phone"`
	Content    string    `gorm:"column:content"`
	CreatedAt  time.Time // 创建时间
	UpdatedAt  time.Time // 更新时间
	DeletedAt  *gorm.DeletedAt
}

func (f Feedback) TableName() string {
	return "feedback"
}

// domain->db
func (f *Feedback) FromDomainFeedback(d domain.FeedbackIntf) {
	if f == nil || d == nil {
		return
	}

	f.FeedbackID = d.GetFeedbackID()
	f.TenantID = d.GetTenantID()
	f.Phone = d.GetPhone()
	f.Content = d.GetContent()
}

// db->domain
func (f *Feedback) ToDomainFeedback() domain.FeedbackIntf {
	if f == nil {
		return nil
	}

	p := domain.Feedback{
		FeedbackID: f.FeedbackID,
		TenantID:   f.TenantID,
		Phone:      f.Phone,
		Content:    f.Content,
	}
	return &p
}
