package domain

import (
	"time"
)

type RecommendedProductStatistics struct {
	RecommendedProductStatisticsID string
	OrganizationID                 string
	TenantID                       string
	ReportID                       string
	RecommendedProductID           string
	Symptom                        string
	ExposedAt                      time.Time
}

type RecommendedProductStatisticsMapper interface {
	ToDomainRecommendedProductStatisticsMapper
	FromDomainRecommendedProductStatisticsMapper
}

type ToDomainRecommendedProductStatisticsMapper interface {
	ToDomainRecommendedProductStatistics() RecommendedProductStatisticsIntf
}

type FromDomainRecommendedProductStatisticsMapper interface {
	FromDomainRecommendedProductStatistics(RecommendedProductStatisticsIntf)
}

type RecommendedProductStatisticsIntf interface {
	GetRecommendedProductStatisticsID() string
	GetOrganizationID() string
	GetTenantID() string
	GetReportID() string
	GetRecommendedProductID() string
	GetSymptom() string
	GetExposedAt() time.Time
}

var _ RecommendedProductStatisticsIntf = (*RecommendedProductStatistics)(nil)

func (r *RecommendedProductStatistics) GetRecommendedProductStatisticsID() string {
	if r == nil {
		return ""
	}
	return r.RecommendedProductStatisticsID
}

func (r *RecommendedProductStatistics) GetOrganizationID() string {
	if r == nil {
		return ""
	}
	return r.OrganizationID
}

func (r *RecommendedProductStatistics) GetTenantID() string {
	if r == nil {
		return ""
	}
	return r.TenantID
}

func (r *RecommendedProductStatistics) GetReportID() string {
	if r == nil {
		return ""
	}
	return r.ReportID
}

func (r *RecommendedProductStatistics) GetRecommendedProductID() string {
	if r == nil {
		return ""
	}
	return r.RecommendedProductID
}

func (r *RecommendedProductStatistics) GetSymptom() string {
	if r == nil {
		return ""
	}
	return r.Symptom
}

func (r *RecommendedProductStatistics) GetExposedAt() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.ExposedAt
}
