package report

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"gorm.io/gorm"
)

type InquiryDiagnosis struct {
	InquiryID           string          `gorm:"column:inquiry_id"`
	Content             string          `gorm:"column:content"`
	IsMultipleSelection bool            `gorm:"column:is_multiple_selection"`
	Selections          string          `gorm:"column:selections"`
	CreatedAt           time.Time       // 创建时间
	UpdatedAt           time.Time       // 更新时间
	DeletedAt           *gorm.DeletedAt // 删除时间
}

func (i InquiryDiagnosis) TableName() string {
	return "inquiry_diagnosis"
}

// domain->db
func (i *InquiryDiagnosis) FromDomainInquiryDiagnosis(d domain.InquiryDiagnosisIntf) {
	if i == nil || d == nil {
		return
	}

	i.InquiryID = d.GetInquiryID()
	i.Content = d.GetContent()
	i.IsMultipleSelection = d.GetIsMultipleSelection()
	i.Selections = d.GetSelections()
}

// db->domain
func (i *InquiryDiagnosis) ToDomainInquiryDiagnosis() domain.InquiryDiagnosisIntf {
	if i == nil {
		return nil
	}

	p := domain.InquiryDiagnosis{
		InquiryID:           i.InquiryID,
		Content:             i.Content,
		IsMultipleSelection: i.IsMultipleSelection,
		Selections:          i.Selections,
	}
	return &p
}
