package domain

import "github.com/rs/xid"

type StaffReportRank struct {
	StaffID          xid.ID
	TodayMeasurement int64
	TotalMeasurement int64
}

type StaffReportRankMapper interface {
	ToDomainStaffReportRankMapper
	FromDomainStaffReportRankMapper
}

type ToDomainStaffReportRankMapper interface {
	ToDomainStaffReportRank() StaffReportRankIntf
}

type FromDomainStaffReportRankMapper interface {
	FromDomainStaffReportRank(StaffReportRankIntf)
}

type StaffReportRankIntf interface {
	GetStaffID() xid.ID
	GetTodayMeasurement() int64
	GetTotalMeasurement() int64
}

var _ StaffReportRankIntf = (*StaffReportRank)(nil)

func (s *StaffReportRank) GetStaffID() xid.ID {
	if s == nil {
		return xid.NilID()
	}
	return s.StaffID
}

func (s *StaffReportRank) GetTodayMeasurement() int64 {
	if s == nil {
		return 0
	}
	return s.TodayMeasurement
}

func (s *StaffReportRank) GetTotalMeasurement() int64 {
	if s == nil {
		return 0
	}
	return s.TotalMeasurement
}
