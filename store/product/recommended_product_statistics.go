package product

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
)

type RecommendedProductStatistics struct {
	RecommendedProductStatisticsID string    `gorm:"primary_key;column:recommended_product_statistics_id"`
	OrganizationID                 string    `gorm:"column:organization_id"`
	TenantID                       string    `gorm:"column:tenant_id"`
	ReportID                       string    `gorm:"column:report_id"`
	RecommendedProductID           string    `gorm:"column:recommended_product_id"`
	Symptom                        string    `gorm:"column:symptom"`
	ExposedAt                      time.Time `gorm:"column:exposed_at"`
	CreatedAt                      time.Time // 创建时间
}

func (r RecommendedProductStatistics) TableName() string {
	return "recommended_product_statistics"
}

// domain->db
func (r *RecommendedProductStatistics) FromDomainRecommendedProductStatistics(d domain.RecommendedProductStatisticsIntf) {
	if r == nil || d == nil {
		return
	}

	r.RecommendedProductStatisticsID = d.GetRecommendedProductStatisticsID()
	r.OrganizationID = d.GetOrganizationID()
	r.TenantID = d.GetTenantID()
	r.ReportID = d.GetReportID()
	r.RecommendedProductID = d.GetRecommendedProductID()
	r.Symptom = d.GetSymptom()
	r.ExposedAt = d.GetExposedAt()
}

// db->domain
func (r *RecommendedProductStatistics) ToDomainRecommendedProductStatistics() domain.RecommendedProductStatisticsIntf {
	if r == nil {
		return nil
	}

	p := domain.RecommendedProductStatistics{
		RecommendedProductStatisticsID: r.RecommendedProductStatisticsID,
		OrganizationID:                 r.OrganizationID,
		TenantID:                       r.TenantID,
		ReportID:                       r.ReportID,
		RecommendedProductID:           r.RecommendedProductID,
		Symptom:                        r.Symptom,
		ExposedAt:                      r.ExposedAt,
	}
	return &p
}
