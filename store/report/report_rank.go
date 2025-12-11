package report

import "github.com/jinmukeji/huimaibao-service/svc/report/domain"

type ReportRank struct {
	TenantID    string `gorm:"column:tenant_id"`
	ReportCount int64  `gorm:"column:report_count"`
}

// domain->db
func (r *ReportRank) FromDomainReportRank(d domain.ReportRankIntf) {
	if r == nil || d == nil {
		return
	}

	r.TenantID = d.GetTenantID()
	r.ReportCount = d.GetReportCount()
}

// db->domain
func (r *ReportRank) ToDomainReportRank() domain.ReportRankIntf {
	if r == nil {
		return nil
	}

	p := domain.ReportRank{
		TenantID:    r.TenantID,
		ReportCount: r.ReportCount,
	}
	return &p
}
