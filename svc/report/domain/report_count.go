package domain

import "time"

type ReportCount struct {
	OrganizationID string
	TenantID       string
	CustomerCount  int64
	TempCount      int64
	Date           time.Time
}

type ReportCountMapper interface {
	ToDomainReportCountMapper
	FromDomainReportCountMapper
}

type ToDomainReportCountMapper interface {
	ToDomainReportCount() ReportCountIntf
}

type FromDomainReportCountMapper interface {
	FromDomainReportCount(ReportCountIntf)
}

type ReportCountIntf interface {
	GetOrganizationID() string
	GetTenantID() string
	GetCustomerCount() int64
	GetTempCount() int64
	GetDate() time.Time
}

var _ ReportCountIntf = (*ReportCount)(nil)

func (r *ReportCount) GetOrganizationID() string {
	if r == nil {
		return ""
	}
	return r.OrganizationID
}

func (r *ReportCount) GetTenantID() string {
	if r == nil {
		return ""
	}
	return r.TenantID
}

func (r *ReportCount) GetCustomerCount() int64 {
	if r == nil {
		return 0
	}
	return r.CustomerCount
}

func (r *ReportCount) GetTempCount() int64 {
	if r == nil {
		return 0
	}
	return r.TempCount
}

func (r *ReportCount) GetDate() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.Date
}
