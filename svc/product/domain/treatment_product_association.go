package domain

import (
	"time"
)

type TreatmentProductAssociation struct {
	TreatmentItemID   string
	TreatmentRevID    string
	TreatmentItemType int32
	// item匹配的symptom
	Symptom              string
	RecommendedProductID string
	OrganizationID       string
	ProductImageUrl      string
	ProductName          string
	ProductType          int32
	ProductStatus        int32
	ProductIntroduction  string
	Remarks              string
	DrugName             string
	DrugClassification   int32
	ApprovedNumber       string
	DrugValidityPeriod   time.Time
}

type TreatmentProductAssociationMapper interface {
	ToDomainTreatmentProductAssociationMapper
	FromDomainTreatmentProductAssociationMapper
}

type ToDomainTreatmentProductAssociationMapper interface {
	ToDomainTreatmentProductAssociation() TreatmentProductAssociationIntf
}

type FromDomainTreatmentProductAssociationMapper interface {
	FromDomainTreatmentProductAssociation(TreatmentProductAssociationIntf)
}

type TreatmentProductAssociationIntf interface {
	GetTreatmentItemID() string
	GetTreatmentRevID() string
	GetTreatmentItemType() int32
	GetSymptom() string
	GetRecommendedProductID() string
	GetOrganizationID() string
	GetProductImageUrl() string
	GetProductName() string
	GetProductType() int32
	GetProductStatus() int32
	GetProductIntroduction() string
	GetRemarks() string
	GetDrugName() string
	GetDrugClassification() int32
	GetApprovedNumber() string
	GetDrugValidityPeriod() time.Time
}

var _ TreatmentProductAssociationIntf = (*TreatmentProductAssociation)(nil)

func (t *TreatmentProductAssociation) GetTreatmentItemID() string {
	if t == nil {
		return ""
	}
	return t.TreatmentItemID
}

func (t *TreatmentProductAssociation) GetTreatmentRevID() string {
	if t == nil {
		return ""
	}
	return t.TreatmentRevID
}

func (t *TreatmentProductAssociation) GetTreatmentItemType() int32 {
	if t == nil {
		return 0
	}
	return t.TreatmentItemType
}

func (t *TreatmentProductAssociation) GetSymptom() string {
	if t == nil {
		return ""
	}
	return t.Symptom
}

func (t *TreatmentProductAssociation) GetRecommendedProductID() string {
	if t == nil {
		return ""
	}
	return t.RecommendedProductID
}

func (t *TreatmentProductAssociation) GetOrganizationID() string {
	if t == nil {
		return ""
	}
	return t.OrganizationID
}

func (t *TreatmentProductAssociation) GetProductImageUrl() string {
	if t == nil {
		return ""
	}
	return t.ProductImageUrl
}

func (t *TreatmentProductAssociation) GetProductName() string {
	if t == nil {
		return ""
	}
	return t.ProductName
}

func (t *TreatmentProductAssociation) GetProductType() int32 {
	if t == nil {
		return 0
	}
	return t.ProductType
}

func (t *TreatmentProductAssociation) GetProductStatus() int32 {
	if t == nil {
		return 0
	}
	return t.ProductStatus
}

func (t *TreatmentProductAssociation) GetProductIntroduction() string {
	if t == nil {
		return ""
	}
	return t.ProductIntroduction
}

func (t *TreatmentProductAssociation) GetRemarks() string {
	if t == nil {
		return ""
	}
	return t.Remarks
}

func (t *TreatmentProductAssociation) GetDrugName() string {
	if t == nil {
		return ""
	}
	return t.DrugName
}

func (t *TreatmentProductAssociation) GetDrugClassification() int32 {
	if t == nil {
		return 0
	}
	return t.DrugClassification
}

func (t *TreatmentProductAssociation) GetApprovedNumber() string {
	if t == nil {
		return ""
	}
	return t.ApprovedNumber
}

func (t *TreatmentProductAssociation) GetDrugValidityPeriod() time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.DrugValidityPeriod
}
