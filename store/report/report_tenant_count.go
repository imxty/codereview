package report

import (
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
)

type ReportTenantCount struct {
	OrganizationID string  `gorm:"column:organization_id"`
	TenantID       string  `gorm:"column:tenant_id"`
	MonthlyCount   int32   `gorm:"column:monthly_count"`
	YearOnYear     float64 `gorm:"column:year_on_year"`
	MonthOnMonth   float64 `gorm:"column:month_on_month"`
	Date           string  `gorm:"column:date"`
}

// domain->db
func (r *ReportTenantCount) FromDomainReportTenantCount(d domain.ReportTenantCountIntf) {
	if r == nil || d == nil {
		return
	}

	r.OrganizationID = d.GetOrganizationID()
	r.TenantID = d.GetTenantID()
	r.MonthlyCount = d.GetMonthlyCount()
	r.YearOnYear = d.GetYearOnYear()
	r.MonthOnMonth = d.GetMonthOnMonth()
	r.Date = d.GetDate()
}

// db->domain
func (r *ReportTenantCount) ToDomainReportTenantCount() domain.ReportTenantCountIntf {
	if r == nil {
		return nil
	}

	p := domain.ReportTenantCount{
		OrganizationID: r.OrganizationID,
		TenantID:       r.TenantID,
		MonthlyCount:   r.MonthlyCount,
		YearOnYear:     r.YearOnYear,
		MonthOnMonth:   r.MonthOnMonth,
		Date:           r.Date,
	}
	return &p
}
