package report

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
)

type ReportStaffCount struct {
	OrganizationID string    `gorm:"column:organization_id"`
	TenantID       string    `gorm:"column:tenant_id"`
	StaffID        string    `gorm:"column:staff_id"`
	CustomerCount  int64     `gorm:"column:customer_count"`
	TempCount      int64     `gorm:"column:temp_count"`
	Date           time.Time `gorm:"column:date"`
}

// domain->db
func (r *ReportStaffCount) FromDomainReportStaffCount(d domain.ReportStaffCountIntf) {
	if r == nil || d == nil {
		return
	}

	r.OrganizationID = d.GetOrganizationID()
	r.TenantID = d.GetTenantID()
	r.StaffID = d.GetStaffID()
	r.CustomerCount = d.GetCustomerCount()
	r.TempCount = d.GetTempCount()
	r.Date = d.GetDate()
}

// db->domain
func (r *ReportStaffCount) ToDomainReportStaffCount() domain.ReportStaffCountIntf {
	if r == nil {
		return nil
	}

	p := domain.ReportStaffCount{
		OrganizationID: r.OrganizationID,
		TenantID:       r.TenantID,
		StaffID:        r.StaffID,
		CustomerCount:  r.CustomerCount,
		TempCount:      r.TempCount,
		Date:           r.Date,
	}
	return &p
}
