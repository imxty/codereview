package domain

import "time"

func (t *TenantSubscriptionMonthly) GetDate() time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.Date
}

func (t *TenantSubscriptionMonthly) GetTenantCount() int32 {
	if t == nil {
		return 0
	}
	return t.TenantCount
}

func (t *TenantSubscriptionMonthly) GetYearCount() int32 {
	if t == nil {
		return 0
	}
	return t.YearCount
}

type TenantSubscriptionMonthly struct {
	Date        time.Time
	TenantCount int32
	YearCount   int32
}

type TenantSubscriptionMonthlyMapper interface {
	ToDomainTenantSubscriptionMonthlyMapper
	FromDomainTenantSubscriptionMonthlyMapper
}

type ToDomainTenantSubscriptionMonthlyMapper interface {
	ToDomainTenantSubscriptionMonthly() TenantSubscriptionMonthlyIntf
}

type FromDomainTenantSubscriptionMonthlyMapper interface {
	FromDomainTenantSubscriptionMonthly(TenantSubscriptionMonthlyIntf)
}

type TenantSubscriptionMonthlyIntf interface {
	GetDate() time.Time
	GetTenantCount() int32
	GetYearCount() int32
}

var _ TenantSubscriptionMonthlyIntf = (*TenantSubscriptionMonthly)(nil)
