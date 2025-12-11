package report

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"gorm.io/gorm"
)

// SharedReport 分享报告
type SharedReport struct {
	Token     string          `gorm:"primary_key;column:token"`
	ReportID  string          `gorm:"column:report_id"`
	CreatedAt time.Time       // 创建时间
	UpdatedAt time.Time       // 更新时间
	DeletedAt *gorm.DeletedAt // 删除时间
}

func (s SharedReport) TableName() string {
	return "shared_report"
}

// domain->db
func (s *SharedReport) FromDomainSharedReport(d domain.SharedReportIntf) {
	if s == nil || d == nil {
		return
	}

	s.Token = d.GetToken()
	s.ReportID = d.GetReportID()
}

// db->domain
func (s *SharedReport) ToDomainSharedReport() domain.SharedReportIntf {
	if s == nil {
		return nil
	}

	p := domain.SharedReport{
		Token:    s.Token,
		ReportID: s.ReportID,
	}
	return &p
}
