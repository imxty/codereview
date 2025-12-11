package review

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"gorm.io/gorm"
)

// ReviewNotification
type ReviewNotification struct {
	ReviewNotificationID string          `gorm:"primary_key;column:review_notification_id"`
	OrganizationID       string          `gorm:"column:organization_id"`
	TenantID             string          `gorm:"column:tenant_id"`
	NotificationTitle    string          `gorm:"column:notification_title"`
	TargetType           int32           `gorm:"column:target_type"`
	Result               bool            `gorm:"column:result"`
	Comment              string          `gorm:"column:comment"`
	HasRead              bool            `gorm:"column:has_read"`
	Rev                  int32           `gorm:"column:rev"`
	CreatedAt            time.Time       // 创建时间
	UpdatedAt            time.Time       // 更新时间
	DeletedAt            *gorm.DeletedAt // 删除时间
}

func (r ReviewNotification) TableName() string {
	return "review_notification"
}

// domain->db
func (r *ReviewNotification) FromDomainReviewNotification(d domain.ReviewNotificationIntf) {
	if r == nil || d == nil {
		return
	}

	r.ReviewNotificationID = d.GetReviewNotificationID()
	r.OrganizationID = d.GetOrganizationID()
	r.TenantID = d.GetTenantID()
	r.NotificationTitle = d.GetNotificationTitle()
	r.TargetType = d.GetTargetType()
	r.Result = d.GetResult()
	r.Comment = d.GetComment()
	r.HasRead = d.GetHasRead()
	r.Rev = d.GetRev()
}

// db->domain
func (r *ReviewNotification) ToDomainReviewNotification() domain.ReviewNotificationIntf {
	if r == nil {
		return nil
	}

	p := domain.ReviewNotification{
		ReviewNotificationID: r.ReviewNotificationID,
		OrganizationID:       r.OrganizationID,
		TenantID:             r.TenantID,
		NotificationTitle:    r.NotificationTitle,
		TargetType:           r.TargetType,
		Result:               r.Result,
		Comment:              r.Comment,
		HasRead:              r.HasRead,
		Rev:                  r.Rev,
		CreatedAt:            r.CreatedAt,
	}
	return &p
}
