package domain

import "time"

func (s *SubscriptionPeriod) GetSubscriptionPeriodID() string {
	if s == nil {
		return ""
	}
	return s.SubscriptionPeriodID
}

func (s *SubscriptionPeriod) GetTenantID() string {
	if s == nil {
		return ""
	}
	return s.TenantID
}

func (s *SubscriptionPeriod) GetOrganizationID() string {
	if s == nil {
		return ""
	}
	return s.OrganizationID
}

func (s *SubscriptionPeriod) GetOrganizationName() string {
	if s == nil {
		return ""
	}
	return s.OrganizationName
}

func (s *SubscriptionPeriod) GetTenantName() string {
	if s == nil {
		return ""
	}
	return s.TenantName
}

func (s *SubscriptionPeriod) GetExpiredTime() time.Time {
	if s == nil {
		return time.Time{}
	}
	return s.ExpiredTime
}

func (s *SubscriptionPeriod) GetUserID() string {
	if s == nil {
		return ""
	}
	return s.UserID
}

func (s *SubscriptionPeriod) GetContactName() string {
	if s == nil {
		return ""
	}
	return s.ContactName
}

func (s *SubscriptionPeriod) GetContactPhone() string {
	if s == nil {
		return ""
	}
	return s.ContactPhone
}

func (s *SubscriptionPeriod) GetYears() int32 {
	if s == nil {
		return 0
	}
	return s.Years
}

func (s *SubscriptionPeriod) GetRev() int32 {
	if s == nil {
		return 0
	}
	return s.Rev
}

func (s *SubscriptionPeriod) GetCreatedAt() time.Time {
	if s == nil {
		return time.Time{}
	}
	return s.CreatedAt
}

type SubscriptionPeriod struct {
	SubscriptionPeriodID string
	TenantID             string
	OrganizationID       string
	OrganizationName     string
	TenantName           string
	ExpiredTime          time.Time
	UserID               string
	ContactName          string
	ContactPhone         string
	Years                int32
	Rev                  int32
	CreatedAt            time.Time
}

type SubscriptionPeriodMapper interface {
	ToDomainSubscriptionPeriodMapper
	FromDomainSubscriptionPeriodMapper
}

type ToDomainSubscriptionPeriodMapper interface {
	ToDomainSubscriptionPeriod() SubscriptionPeriodIntf
}

type FromDomainSubscriptionPeriodMapper interface {
	FromDomainSubscriptionPeriod(SubscriptionPeriodIntf)
}

type SubscriptionPeriodIntf interface {
	GetSubscriptionPeriodID() string
	GetTenantID() string
	GetOrganizationID() string
	GetOrganizationName() string
	GetTenantName() string
	GetExpiredTime() time.Time
	GetUserID() string
	GetContactName() string
	GetContactPhone() string
	GetYears() int32
	GetRev() int32
	GetCreatedAt() time.Time
}

var _ SubscriptionPeriodIntf = (*SubscriptionPeriod)(nil)
