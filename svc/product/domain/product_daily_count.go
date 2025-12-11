package domain

import (
	"time"
)

func (p *ProductDailyCount) GetExposedAt() time.Time {
	if p == nil {
		return time.Time{}
	}
	return p.ExposedAt
}

func (p *ProductDailyCount) GetCount() int32 {
	if p == nil {
		return 0
	}
	return p.Count
}

type ProductDailyCount struct {
	ExposedAt time.Time
	Count     int32
}

type ProductDailyCountMapper interface {
	ToDomainProductDailyCountMapper
	FromDomainProductDailyCountMapper
}

type ToDomainProductDailyCountMapper interface {
	ToDomainProductDailyCount() ProductDailyCountIntf
}

type FromDomainProductDailyCountMapper interface {
	FromDomainProductDailyCount(ProductDailyCountIntf)
}

type ProductDailyCountIntf interface {
	GetExposedAt() time.Time
	GetCount() int32
}

var _ ProductDailyCountIntf = (*ProductDailyCount)(nil)
