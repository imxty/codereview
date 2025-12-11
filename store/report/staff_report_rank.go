package report

import (
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/rs/xid"
)

// StaffReportRank 员工测量排行
type StaffReportRank struct {
	StaffID          xid.ID `gorm:"column:staff_id"`
	TodayMeasurement int64  `gorm:"column:today_measurement"`
	TotalMeasurement int64  `gorm:"column:total_measurement"`
}

func (r StaffReportRank) TableName() string {
	return "report"
}

// domain->db
func (s *StaffReportRank) FromDomainStaffReportRank(d domain.StaffReportRankIntf) {
	if s == nil || d == nil {
		return
	}

	s.StaffID = d.GetStaffID()
	s.TodayMeasurement = d.GetTodayMeasurement()
	s.TotalMeasurement = d.GetTotalMeasurement()
}

// db->domain
func (s *StaffReportRank) ToDomainStaffReportRank() domain.StaffReportRankIntf {
	if s == nil {
		return nil
	}

	p := domain.StaffReportRank{
		StaffID:          s.StaffID,
		TodayMeasurement: s.TodayMeasurement,
		TotalMeasurement: s.TotalMeasurement,
	}
	return &p
}
