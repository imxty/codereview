package domain

const (
	// 无效的方案
	TreatmentStatusInvalid int32 = 0
	// 未设置的方案
	TreatmentStatusUnset int32 = 1
	// 草稿
	TreatmentStatusDraft int32 = 2
	// 审核中
	TreatmentStatusUnderReview int32 = 3
	// 已过审
	TreatmentStatusApproved int32 = 4
	// 未过审
	TreatmentStatusUnapproved int32 = 5
)

const (
	// 方案状态无效
	TreatmentPropertyInValid int32 = 0
	// 方案状态有效
	TreatmentPropertyValid int32 = 1
)

type Treatment struct {
	TreatmentID    string
	OrganizationID string
	LatestRev      string
	Rev            int32
}

type TreatmentMapper interface {
	ToDomainTreatmentMapper
	FromDomainTreatmentMapper
}

type ToDomainTreatmentMapper interface {
	ToDomainTreatment() TreatmentIntf
}

type FromDomainTreatmentMapper interface {
	FromDomainTreatment(TreatmentIntf)
}

type TreatmentIntf interface {
	GetTreatmentID() string
	GetOrganizationID() string
	GetLatestRev() string
	GetRev() int32
}

var _ TreatmentIntf = (*Treatment)(nil)

func (t *Treatment) GetTreatmentID() string {
	if t == nil {
		return ""
	}
	return t.TreatmentID
}

func (t *Treatment) GetOrganizationID() string {
	if t == nil {
		return ""
	}
	return t.OrganizationID
}

func (t *Treatment) GetLatestRev() string {
	if t == nil {
		return ""
	}
	return t.LatestRev
}

func (t *Treatment) GetRev() int32 {
	if t == nil {
		return 0
	}
	return t.Rev
}
