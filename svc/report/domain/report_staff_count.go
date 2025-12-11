package domain

import "time"

type ReportStaffCount struct {
	OrganizationID string
	TenantID       string
	StaffID        string
	CustomerCount  int64
	TempCount      int64
	Date           time.Time
}

type ReportStaffCountMapper interface {
	ToDomainReportStaffCountMapper
	FromDomainReportStaffCountMapper
}

type ToDomainReportStaffCountMapper interface {
	ToDomainReportStaffCount() ReportStaffCountIntf
}

type FromDomainReportStaffCountMapper interface {
	FromDomainReportStaffCount(ReportStaffCountIntf)
}

type ReportStaffCountIntf interface {
	GetOrganizationID() string
	GetTenantID() string
	GetStaffID() string
	GetCustomerCount() int64
	GetTempCount() int64
	GetDate() time.Time
}

var _ ReportStaffCountIntf = (*ReportStaffCount)(nil)

func (r *ReportStaffCount) GetOrganizationID() string {
	if r == nil {
		return ""
	}
	return r.OrganizationID
}

func (r *ReportStaffCount) GetTenantID() string {
	if r == nil {
		return ""
	}
	return r.TenantID
}

func (r *ReportStaffCount) GetStaffID() string {
	if r == nil {
		return ""
	}
	return r.StaffID
}

func (r *ReportStaffCount) GetCustomerCount() int64 {
	if r == nil {
		return 0
	}
	return r.CustomerCount
}

func (r *ReportStaffCount) GetTempCount() int64 {
	if r == nil {
		return 0
	}
	return r.TempCount
}

func (r *ReportStaffCount) GetDate() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.Date
}
