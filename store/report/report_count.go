package report

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
)

type ReportCount struct {
	OrganizationID string    `gorm:"column:organization_id"`
	TenantID       string    `gorm:"column:tenant_id"`
	CustomerCount  int64     `gorm:"column:customer_count"`
	TempCount      int64     `gorm:"column:temp_count"`
	Date           time.Time `gorm:"column:date"`
}

func (r ReportCount) TableName() string {
	return "report_count"
}

// domain->db
func (r *ReportCount) FromDomainReportCount(d domain.ReportCountIntf) {
	if r == nil || d == nil {
		return
	}

	r.OrganizationID = d.GetOrganizationID()
	r.TenantID = d.GetTenantID()
	r.CustomerCount = d.GetCustomerCount()
	r.TempCount = d.GetTempCount()
	r.Date = d.GetDate()
}

// db->domain
func (r *ReportCount) ToDomainReportCount() domain.ReportCountIntf {
	if r == nil {
		return nil
	}

	p := domain.ReportCount{
		OrganizationID: r.OrganizationID,
		TenantID:       r.TenantID,
		CustomerCount:  r.CustomerCount,
		TempCount:      r.TempCount,
		Date:           r.Date,
	}
	return &p
}
