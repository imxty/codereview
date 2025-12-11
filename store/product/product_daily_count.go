package product

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
)

type ProductDailyCount struct {
	ExposedAt time.Time `gorm:"column:exposed_date"`
	Count     int32     `gorm:"column:exposed_count"`
}

func (p ProductDailyCount) TableName() string {
	return "product_daily_count"
}

// domain->db
func (p *ProductDailyCount) FromDomainProductDailyCount(d domain.ProductDailyCountIntf) {
	if p == nil || d == nil {
		return
	}

	p.ExposedAt = d.GetExposedAt()
	p.Count = d.GetCount()
}

// db->domain
func (p *ProductDailyCount) ToDomainProductDailyCount() domain.ProductDailyCountIntf {
	if p == nil {
		return nil
	}

	pr := domain.ProductDailyCount{
		ExposedAt: p.ExposedAt,
		Count:     p.Count,
	}
	return &pr
}
