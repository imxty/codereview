package user

import (
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
)

// Faq 常见问题
type Faq struct {
	FaqID    string `gorm:"primary_key;column:faq_id"`
	Question string `gorm:"column:question"`
	Answer   string `gorm:"column:answer"`
}

func (f Faq) TableName() string {
	return "faq"
}

// domain->db
func (f *Faq) FromDomainFaq(d domain.FaqIntf) {
	if f == nil || d == nil {
		return
	}

	f.FaqID = d.GetFaqID()
	f.Question = d.GetQuestion()
	f.Answer = d.GetAnswer()
}

// db->domain
func (f *Faq) ToDomainFaq() domain.FaqIntf {
	if f == nil {
		return nil
	}

	p := domain.Faq{
		FaqID:    f.FaqID,
		Question: f.Question,
		Answer:   f.Answer,
	}
	return &p
}
