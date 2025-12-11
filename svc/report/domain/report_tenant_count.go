package domain

type ReportTenantCount struct {
	OrganizationID string
	TenantID       string
	MonthlyCount   int32
	YearOnYear     float64
	MonthOnMonth   float64
	Date           string
}

type ReportTenantCountMapper interface {
	ToDomainReportTenantCountMapper
	FromDomainReportTenantCountMapper
}

type ToDomainReportTenantCountMapper interface {
	ToDomainReportTenantCount() ReportTenantCountIntf
}

type FromDomainReportTenantCountMapper interface {
	FromDomainReportTenantCount(ReportTenantCountIntf)
}

type ReportTenantCountIntf interface {
	GetOrganizationID() string
	GetTenantID() string
	GetMonthlyCount() int32
	GetYearOnYear() float64
	GetMonthOnMonth() float64
	GetDate() string
}

var _ ReportTenantCountIntf = (*ReportTenantCount)(nil)

func (r *ReportTenantCount) GetOrganizationID() string {
	if r == nil {
		return ""
	}
	return r.OrganizationID
}

func (r *ReportTenantCount) GetTenantID() string {
	if r == nil {
		return ""
	}
	return r.TenantID
}

func (r *ReportTenantCount) GetMonthlyCount() int32 {
	if r == nil {
		return 0
	}
	return r.MonthlyCount
}

func (r *ReportTenantCount) GetYearOnYear() float64 {
	if r == nil {
		return 0.0
	}
	return r.YearOnYear
}

func (r *ReportTenantCount) GetMonthOnMonth() float64 {
	if r == nil {
		return 0.0
	}
	return r.MonthOnMonth
}

func (r *ReportTenantCount) GetDate() string {
	if r == nil {
		return ""
	}
	return r.Date
}
