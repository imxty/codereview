package domain

const (
	// 无效的方案配置类型
	TreatmentItemTypeInvalid int32 = 0
	// 未设置的方案配置类型
	TreatmentItemTypeUnset int32 = 1
	// 风险疾病类型
	TreatmentItemTypeRiskyDisease int32 = 2
	// 脏腑辩证类型
	TreatmentItemTypeDirtyDialectic int32 = 3
	// 理疗类型
	TreatmentItemTypePhysicalTherapy int32 = 4
	// 体质类型
	TreatmentItemTypePhysical int32 = 5
)

type TreatmentItem struct {
	TreatmentItemID      string
	TreatmentRevID       string
	TreatmentItemType    int32
	Symptom              string
	RecommendedProductID string
}

type TreatmentItemMapper interface {
	ToDomainTreatmentItemMapper
	FromDomainTreatmentItemMapper
}

type ToDomainTreatmentItemMapper interface {
	ToDomainTreatmentItem() TreatmentItemIntf
}

type FromDomainTreatmentItemMapper interface {
	FromDomainTreatmentItem(TreatmentItemIntf)
}

type TreatmentItemIntf interface {
	GetTreatmentItemID() string
	GetTreatmentRevID() string
	GetTreatmentItemType() int32
	GetSymptom() string
	GetRecommendedProductID() string
}

var _ TreatmentItemIntf = (*TreatmentItem)(nil)

func (t *TreatmentItem) GetTreatmentItemID() string {
	if t == nil {
		return ""
	}
	return t.TreatmentItemID
}

func (t *TreatmentItem) GetTreatmentRevID() string {
	if t == nil {
		return ""
	}
	return t.TreatmentRevID
}

func (t *TreatmentItem) GetTreatmentItemType() int32 {
	if t == nil {
		return 0
	}
	return t.TreatmentItemType
}

func (t *TreatmentItem) GetSymptom() string {
	if t == nil {
		return ""
	}
	return t.Symptom
}

func (t *TreatmentItem) GetRecommendedProductID() string {
	if t == nil {
		return ""
	}
	return t.RecommendedProductID
}
