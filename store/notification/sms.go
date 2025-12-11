package notification

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/notification/domain"
	"gorm.io/gorm"
)

type Sms struct {
	SmsID          string          `gorm:"primary_key;column:sms_id"`
	TxID           string          `gorm:"column:tx_id"`
	Phone          string          `gorm:"column:phone"`
	SmsStatus      int32           `gorm:"column:sms_status"`
	TemplateAction string          `gorm:"column:template_action"`
	PlatformType   string          `gorm:"column:platform_type"`
	TemplateParam  string          `gorm:"column:template_param"`
	Language       string          `gorm:"column:language"`
	SerialNumber   string          `gorm:"column:serial_number"`
	SmsErrorLog    string          `gorm:"column:sms_error_log"`
	Rev            int32           `gorm:"column:rev"`
	CreatedAt      time.Time       // 创建时间
	UpdatedAt      time.Time       // 更新时间
	DeletedAt      *gorm.DeletedAt // 删除时间
}

func (s Sms) TableName() string {
	return "sms"
}

// domain->db
func (s *Sms) FromDomainSms(d domain.SmsIntf) {
	if s == nil || d == nil {
		return
	}

	s.SmsID = d.GetSmsID()
	s.TxID = d.GetTxID()
	s.Phone = d.GetPhone()
	s.SmsStatus = d.GetSmsStatus()
	s.TemplateAction = d.GetTemplateAction()
	s.PlatformType = d.GetPlatformType()
	s.TemplateParam = d.GetTemplateParam()
	s.Language = d.GetLanguage()
	s.SerialNumber = d.GetSerialNumber()
	s.SmsErrorLog = d.GetSmsErrorLog()
	s.Rev = d.GetRev()
}

// db->domain
func (s *Sms) ToDomainSms() domain.SmsIntf {
	if s == nil {
		return nil
	}

	p := domain.Sms{
		SmsID:          s.SmsID,
		TxID:           s.TxID,
		Phone:          s.Phone,
		SmsStatus:      s.SmsStatus,
		TemplateAction: s.TemplateAction,
		PlatformType:   s.PlatformType,
		TemplateParam:  s.TemplateParam,
		Language:       s.Language,
		SerialNumber:   s.SerialNumber,
		SmsErrorLog:    s.SmsErrorLog,
		Rev:            s.Rev,
	}
	return &p
}
