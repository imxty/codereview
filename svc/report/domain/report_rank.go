package domain

type ReportRank struct {
	TenantID    string
	ReportCount int64
}

type ReportRankMapper interface {
	ToDomainReportRankMapper
	FromDomainReportRankMapper
}

type ToDomainReportRankMapper interface {
	ToDomainReportRank() ReportRankIntf
}

type FromDomainReportRankMapper interface {
	FromDomainReportRank(ReportRankIntf)
}

type ReportRankIntf interface {
	GetTenantID() string
	GetReportCount() int64
}

var _ ReportRankIntf = (*ReportRank)(nil)

func (r *ReportRank) GetTenantID() string {
	if r == nil {
		return ""
	}
	return r.TenantID
}

func (r *ReportRank) GetReportCount() int64 {
	if r == nil {
		return 0
	}
	return r.ReportCount
}
