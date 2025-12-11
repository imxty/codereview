package domain

type InquiryQuestion struct {
	InquiryID           string              `json:"inquiry_id"`
	Content             string              `json:"content"`
	IsMultipleSelection bool                `json:"is_multiple_selection"`
	Items               []InquiryAnswerItem `json:"items"`
}

type InquiryAnswerItem struct {
	AnswerID  string   `json:"answer_id"`
	Content   string   `json:"content"`
	Exclusive []string `json:"exclusive"`
}

type InquiryAnswer struct {
	InquiryID string   `json:"inquiry_id"`
	Answers   []string `json:"answers"`
}

type InquiryDiagnosis struct {
	InquiryID           string
	Content             string
	IsMultipleSelection bool
	Selections          string
}

type InquiryDiagnosisMapper interface {
	ToDomainInquiryDiagnosisMapper
	FromDomainInquiryDiagnosisMapper
}

type ToDomainInquiryDiagnosisMapper interface {
	ToDomainInquiryDiagnosis() InquiryDiagnosisIntf
}

type FromDomainInquiryDiagnosisMapper interface {
	FromDomainInquiryDiagnosis(InquiryDiagnosisIntf)
}

type InquiryDiagnosisIntf interface {
	GetInquiryID() string
	GetContent() string
	GetIsMultipleSelection() bool
	GetSelections() string
}

var _ InquiryDiagnosisIntf = (*InquiryDiagnosis)(nil)

func (i *InquiryDiagnosis) GetInquiryID() string {
	if i == nil {
		return ""
	}

	return i.InquiryID
}

func (i *InquiryDiagnosis) GetIsMultipleSelection() bool {
	if i == nil {
		return false
	}

	return i.IsMultipleSelection
}

func (i *InquiryDiagnosis) GetContent() string {
	if i == nil {
		return ""
	}

	return i.Content
}

func (i *InquiryDiagnosis) GetSelections() string {
	if i == nil {
		return ""
	}

	return i.Selections
}
