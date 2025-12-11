package domain

func (t *TenantTreatment) GetTenantTreatmentID() string {
	if t == nil {
		return ""
	}
	return t.TenantTreatmentID
}

func (t *TenantTreatment) GetTreatmentID() string {
	if t == nil {
		return ""
	}
	return t.TreatmentID
}

func (t *TenantTreatment) GetTenantID() string {
	if t == nil {
		return ""
	}
	return t.TenantID
}

func (t *TenantTreatment) GetRev() int32 {
	if t == nil {
		return 0
	}
	return t.Rev
}

type TenantTreatment struct {
	TenantTreatmentID string
	TreatmentID       string
	TenantID          string
	Rev               int32
}

type TenantTreatmentMapper interface {
	ToDomainTenantTreatmentMapper
	FromDomainTenantTreatmentMapper
}

type ToDomainTenantTreatmentMapper interface {
	ToDomainTenantTreatment() TenantTreatmentIntf
}

type FromDomainTenantTreatmentMapper interface {
	FromDomainTenantTreatment(TenantTreatmentIntf)
}

type TenantTreatmentIntf interface {
	GetTenantTreatmentID() string
	GetTreatmentID() string
	GetTenantID() string
	GetRev() int32
}

var _ TenantTreatmentIntf = (*TenantTreatment)(nil)
