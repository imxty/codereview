package domain

type Faq struct {
	FaqID    string
	Question string
	Answer   string
}

type FaqMapper interface {
	ToDomainFaqMapper
	FromDomainFaqMapper
}

type ToDomainFaqMapper interface {
	ToDomainFaq() FaqIntf
}

type FromDomainFaqMapper interface {
	FromDomainFaq(FaqIntf)
}

type FaqIntf interface {
	GetFaqID() string
	GetQuestion() string
	GetAnswer() string
}

var _ FaqIntf = (*Faq)(nil)

func (f *Faq) GetFaqID() string {
	if f == nil {
		return ""
	}
	return f.FaqID
}

func (f *Faq) GetQuestion() string {
	if f == nil {
		return ""
	}
	return f.Question
}

func (f *Faq) GetAnswer() string {
	if f == nil {
		return ""
	}
	return f.Answer
}
