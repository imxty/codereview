package domain

import (
	"time"
)

const (
	// 方案未发布
	TreatmentRevIsNotPublished = false
	// 方案发布
	TreatmentRevIsPublished = true
)

type TreatmentRev struct {
	TreatmentRevID  string
	TreatmentID     string
	TreatmentName   string
	Remarks         string
	OrganizationID  string
	IsPublished     bool
	TreatmentStatus int32
	Rev             int32
	CreatedAt       time.Time
}

type TreatmentRevMapper interface {
	ToDomainTreatmentRevMapper
	FromDomainTreatmentRevMapper
}

type ToDomainTreatmentRevMapper interface {
	ToDomainTreatmentRev() TreatmentRevIntf
}

type FromDomainTreatmentRevMapper interface {
	FromDomainTreatmentRev(TreatmentRevIntf)
}

type TreatmentRevIntf interface {
	GetTreatmentRevID() string
	GetTreatmentID() string
	GetTreatmentName() string
	GetRemarks() string
	GetOrganizationID() string
	GetIsPublished() bool
	GetTreatmentStatus() int32
	GetRev() int32
	GetCreatedAt() time.Time
}

var _ TreatmentRevIntf = (*TreatmentRev)(nil)

func (t *TreatmentRev) GetTreatmentRevID() string {
	if t == nil {
		return ""
	}
	return t.TreatmentRevID
}

func (t *TreatmentRev) GetTreatmentID() string {
	if t == nil {
		return ""
	}
	return t.TreatmentID
}

func (t *TreatmentRev) GetTreatmentName() string {
	if t == nil {
		return ""
	}
	return t.TreatmentName
}

func (t *TreatmentRev) GetRemarks() string {
	if t == nil {
		return ""
	}
	return t.Remarks
}

func (t *TreatmentRev) GetOrganizationID() string {
	if t == nil {
		return ""
	}
	return t.OrganizationID
}

func (t *TreatmentRev) GetIsPublished() bool {
	if t == nil {
		return false
	}
	return t.IsPublished
}

func (t *TreatmentRev) GetTreatmentStatus() int32 {
	if t == nil {
		return 0
	}
	return t.TreatmentStatus
}

func (t *TreatmentRev) GetRev() int32 {
	if t == nil {
		return 0
	}
	return t.Rev
}

func (r *TreatmentRev) GetCreatedAt() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.CreatedAt
}
