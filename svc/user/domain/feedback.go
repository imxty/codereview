package domain

func (f *Feedback) GetFeedbackID() string {
	if f == nil {
		return ""
	}
	return f.FeedbackID
}

func (f *Feedback) GetTenantID() string {
	if f == nil {
		return ""
	}
	return f.TenantID
}

func (f *Feedback) GetPhone() string {
	if f == nil {
		return ""
	}
	return f.Phone
}

func (f *Feedback) GetContent() string {
	if f == nil {
		return ""
	}
	return f.Content
}

type Feedback struct {
	FeedbackID string
	TenantID   string
	Phone      string
	Content    string
}

type FeedbackMapper interface {
	ToDomainFeedbackMapper
	FromDomainFeedbackMapper
}

type ToDomainFeedbackMapper interface {
	ToDomainFeedback() FeedbackIntf
}

type FromDomainFeedbackMapper interface {
	FromDomainFeedback(FeedbackIntf)
}

type FeedbackIntf interface {
	GetFeedbackID() string
	GetTenantID() string
	GetPhone() string
	GetContent() string
}

var _ FeedbackIntf = (*Feedback)(nil)
