package user

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
)

type TenantSubscriptionMonthly struct {
	Date        time.Time
	TenantCount int32
	YearCount   int32
}

// domain->db
func (t *TenantSubscriptionMonthly) FromDomainTenantSubscriptionMonthly(d domain.TenantSubscriptionMonthlyIntf) {
	if t == nil || d == nil {
		return
	}

	t.Date = d.GetDate()
	t.TenantCount = d.GetTenantCount()
	t.YearCount = d.GetYearCount()
}

// db->domain
func (t *TenantSubscriptionMonthly) ToDomainTenantSubscriptionMonthly() domain.TenantSubscriptionMonthlyIntf {
	if t == nil {
		return nil
	}

	p := domain.TenantSubscriptionMonthly{
		Date:        t.Date,
		TenantCount: t.TenantCount,
		YearCount:   t.YearCount,
	}
	return &p
}
