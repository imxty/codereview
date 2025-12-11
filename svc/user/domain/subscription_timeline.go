package domain

import "time"

func (s *SubscriptionTimeline) GetSubscriptionTimelineID() string {
	if s == nil {
		return ""
	}
	return s.SubscriptionTimelineID
}

func (s *SubscriptionTimeline) GetTenantID() string {
	if s == nil {
		return ""
	}
	return s.TenantID
}

func (s *SubscriptionTimeline) GetStartTime() time.Time {
	if s == nil {
		return time.Time{}
	}
	return s.StartTime
}

func (s *SubscriptionTimeline) GetEndTime() time.Time {
	if s == nil {
		return time.Time{}
	}
	return s.EndTime
}

func (s *SubscriptionTimeline) GetRev() int32 {
	if s == nil {
		return 0
	}
	return s.Rev
}

func (s *SubscriptionTimeline) GetYears() int32 {
	if s == nil {
		return 0
	}
	return s.Years
}

type SubscriptionTimeline struct {
	SubscriptionTimelineID string
	TenantID               string
	StartTime              time.Time
	EndTime                time.Time
	Rev                    int32
	Years                  int32
}

type SubscriptionTimelineMapper interface {
	ToDomainSubscriptionTimelineMapper
	FromDomainSubscriptionTimelineMapper
}

type ToDomainSubscriptionTimelineMapper interface {
	ToDomainSubscriptionTimeline() SubscriptionTimelineIntf
}

type FromDomainSubscriptionTimelineMapper interface {
	FromDomainSubscriptionTimeline(SubscriptionTimelineIntf)
}

type SubscriptionTimelineIntf interface {
	GetSubscriptionTimelineID() string
	GetTenantID() string
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetRev() int32
	GetYears() int32
}

var _ SubscriptionTimelineIntf = (*SubscriptionTimeline)(nil)
