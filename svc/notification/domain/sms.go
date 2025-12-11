package domain

import "time"

const (
	// SendSucceed 发送成功
	SmsStatusSentSuccessful int32 = 2
	// SendFailed 发送失败
	SmsStatusSendFailed int32 = 3
)

type Sms struct {
	SmsID          string
	TxID           string
	Phone          string
	SmsStatus      int32
	TemplateAction string
	PlatformType   string
	TemplateParam  string
	Language       string
	SerialNumber   string
	SmsErrorLog    string
	Rev            int32
	CreatedAt      time.Time // 创建时间
}

type SmsTemplateParam struct {
	Code string `json:"code"`
}

type SmsMapper interface {
	ToDomainSmsMapper
	FromDomainSmsMapper
}

type ToDomainSmsMapper interface {
	ToDomainSms() SmsIntf
}

type FromDomainSmsMapper interface {
	FromDomainSms(SmsIntf)
}

type SmsIntf interface {
	GetSmsID() string
	GetTxID() string
	GetPhone() string
	GetSmsStatus() int32
	GetTemplateAction() string
	GetPlatformType() string
	GetTemplateParam() string
	GetLanguage() string
	GetSerialNumber() string
	GetSmsErrorLog() string
	GetRev() int32
	GetCreatedAt() time.Time
}

var _ SmsIntf = (*Sms)(nil)

func (s *Sms) GetSmsID() string {
	if s == nil {
		return ""
	}
	return s.SmsID
}

func (s *Sms) GetTxID() string {
	if s == nil {
		return ""
	}
	return s.TxID
}

func (s *Sms) GetPhone() string {
	if s == nil {
		return ""
	}
	return s.Phone
}

func (s *Sms) GetSmsStatus() int32 {
	if s == nil {
		return 0
	}
	return s.SmsStatus
}

func (s *Sms) GetTemplateAction() string {
	if s == nil {
		return ""
	}
	return s.TemplateAction
}

func (s *Sms) GetPlatformType() string {
	if s == nil {
		return ""
	}
	return s.PlatformType
}

func (s *Sms) GetTemplateParam() string {
	if s == nil {
		return ""
	}
	return s.TemplateParam
}

func (s *Sms) GetLanguage() string {
	if s == nil {
		return ""
	}
	return s.Language
}

func (s *Sms) GetSerialNumber() string {
	if s == nil {
		return ""
	}
	return s.SerialNumber
}

func (s *Sms) GetSmsErrorLog() string {
	if s == nil {
		return ""
	}
	return s.SmsErrorLog
}

func (s *Sms) GetRev() int32 {
	if s == nil {
		return 0
	}
	return s.Rev
}

func (s *Sms) GetCreatedAt() time.Time {
	if s == nil {
		return time.Time{}
	}
	return s.CreatedAt
}
