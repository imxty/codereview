package client

type SmsSend interface {
	// SendSms 发送短信
	SendSms(phoneNumber, signName string, templateAction TemplateAction, language TemplateLanguage, templateParam map[string]string) (bool, error)
	// SendNotification 发送通知
	SendNotification(phoneNumber []string, signName string, templateAction TemplateAction, templateParam map[string]string) (bool, error)
}
