package product

import (
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
)

type SymptomCount struct {
	Symptom string `gorm:"column:symptom"`
	Count   int32  `gorm:"column:exposure_count"`
}

func (s SymptomCount) TableName() string {
	return "symptom_count"
}

// domain->db
func (s *SymptomCount) FromDomainSymptomCount(d domain.SymptomCountIntf) {
	if s == nil || d == nil {
		return
	}

	s.Symptom = d.GetSymptom()
	s.Count = d.GetCount()
}

// db->domain
func (s *SymptomCount) ToDomainSymptomCount() domain.SymptomCountIntf {
	if s == nil {
		return nil
	}

	p := domain.SymptomCount{
		Symptom: s.Symptom,
		Count:   s.Count,
	}
	return &p
}
