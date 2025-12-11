package domain

type SymptomCount struct {
	Symptom string
	Count   int32
}

type SymptomCountMapper interface {
	ToDomainSymptomCountMapper
	FromDomainSymptomCountMapper
}

type ToDomainSymptomCountMapper interface {
	ToDomainSymptomCount() SymptomCountIntf
}

type FromDomainSymptomCountMapper interface {
	FromDomainSymptomCount(SymptomCountIntf)
}

type SymptomCountIntf interface {
	GetSymptom() string
	GetCount() int32
}

var _ SymptomCountIntf = (*SymptomCount)(nil)

func (s *SymptomCount) GetSymptom() string {
	if s == nil {
		return ""
	}
	return s.Symptom
}

func (s *SymptomCount) GetCount() int32 {
	if s == nil {
		return 0
	}
	return s.Count
}
