package domain

type TreatmentRevItem struct {
	TreatmentName        string
	TreatmentStatus      int32
	TreatmentItemID      string
	TreatmentRevID       string
	TreatmentItemType    int32
	Symptom              string
	RecommendedProductID string
}

type TreatmentRevItemMapper interface {
	ToDomainTreatmentRevItemMapper
	FromDomainTreatmentRevItemMapper
}

type ToDomainTreatmentRevItemMapper interface {
	ToDomainTreatmentRevItem() TreatmentRevItemIntf
}

type FromDomainTreatmentRevItemMapper interface {
	FromDomainTreatmentRevItem(TreatmentRevItemIntf)
}

type TreatmentRevItemIntf interface {
	GetTreatmentName() string
	GetTreatmentStatus() int32
	GetTreatmentItemID() string
	GetTreatmentRevID() string
	GetTreatmentItemType() int32
	GetSymptom() string
	GetRecommendedProductID() string
}

var _ TreatmentRevItemIntf = (*TreatmentRevItem)(nil)

func (t *TreatmentRevItem) GetTreatmentName() string {
	if t == nil {
		return ""
	}
	return t.TreatmentName
}

func (t *TreatmentRevItem) GetTreatmentStatus() int32 {
	if t == nil {
		return 0
	}
	return t.TreatmentStatus
}

func (t *TreatmentRevItem) GetTreatmentItemID() string {
	if t == nil {
		return ""
	}
	return t.TreatmentItemID
}

func (t *TreatmentRevItem) GetTreatmentRevID() string {
	if t == nil {
		return ""
	}
	return t.TreatmentRevID
}

func (t *TreatmentRevItem) GetTreatmentItemType() int32 {
	if t == nil {
		return 0
	}
	return t.TreatmentItemType
}

func (t *TreatmentRevItem) GetSymptom() string {
	if t == nil {
		return ""
	}
	return t.Symptom
}

func (t *TreatmentRevItem) GetRecommendedProductID() string {
	if t == nil {
		return ""
	}
	return t.RecommendedProductID
}
