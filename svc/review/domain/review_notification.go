package domain

import (
	"time"
)

type ReviewNotification struct {
	ReviewNotificationID string
	OrganizationID       string
	TenantID             string
	NotificationTitle    string
	TargetType           int32
	Result               bool
	Comment              string
	HasRead              bool
	Rev                  int32
	CreatedAt            time.Time
}

type ReviewNotificationMapper interface {
	ToDomainReviewNotificationMapper
	FromDomainReviewNotificationMapper
}

type ToDomainReviewNotificationMapper interface {
	ToDomainReviewNotification() ReviewNotificationIntf
}

type FromDomainReviewNotificationMapper interface {
	FromDomainReviewNotification(ReviewNotificationIntf)
}

type ReviewNotificationIntf interface {
	GetReviewNotificationID() string
	GetOrganizationID() string
	GetTenantID() string
	GetNotificationTitle() string
	GetTargetType() int32
	GetResult() bool
	GetComment() string
	GetHasRead() bool
	GetRev() int32
	GetCreatedAt() time.Time
}

var _ ReviewNotificationIntf = (*ReviewNotification)(nil)

func (r *ReviewNotification) GetReviewNotificationID() string {
	if r == nil {
		return ""
	}
	return r.ReviewNotificationID
}

func (r *ReviewNotification) GetOrganizationID() string {
	if r == nil {
		return ""
	}
	return r.OrganizationID
}

func (r *ReviewNotification) GetTenantID() string {
	if r == nil {
		return ""
	}
	return r.TenantID
}

func (r *ReviewNotification) GetNotificationTitle() string {
	if r == nil {
		return ""
	}
	return r.NotificationTitle
}

func (r *ReviewNotification) GetTargetType() int32 {
	if r == nil {
		return 0
	}
	return r.TargetType
}

func (r *ReviewNotification) GetResult() bool {
	if r == nil {
		return false
	}
	return r.Result
}

func (r *ReviewNotification) GetComment() string {
	if r == nil {
		return ""
	}
	return r.Comment
}

func (r *ReviewNotification) GetHasRead() bool {
	if r == nil {
		return false
	}
	return r.HasRead
}

func (r *ReviewNotification) GetRev() int32 {
	if r == nil {
		return 0
	}
	return r.Rev
}

func (r *ReviewNotification) GetCreatedAt() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.CreatedAt
}
