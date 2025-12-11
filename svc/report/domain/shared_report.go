package domain

type SharedReportMapper interface {
	ToDomainSharedReportMapper
	FromDomainSharedReportMapper
}

type ToDomainSharedReportMapper interface {
	ToDomainSharedReport() SharedReportIntf
}

type FromDomainSharedReportMapper interface {
	FromDomainSharedReport(SharedReportIntf)
}

type SharedReportIntf interface {
	GetToken() string
	GetReportID() string
}

func (s *SharedReport) GetToken() string {
	if s == nil {
		return ""
	}
	return s.Token
}

func (s *SharedReport) GetReportID() string {
	if s == nil {
		return ""
	}
	return s.ReportID
}

type SharedReport struct {
	Token    string
	ReportID string
}
