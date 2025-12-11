package notification

import (
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/svc/notification/client"
	"github.com/jinmukeji/huimaibao-service/svc/notification/domain"
)

type NotificationAPIHandler struct {
	notificationStore domain.NotificationRepository
	aliYunSms         client.SmsSend
	signName          string
}

var _ pb.NotificationAPIHandler = (*NotificationAPIHandler)(nil)

func (svc *NotificationAPIHandler) Name() string {
	const name = "NotificationAPI"
	return name
}

var _ pb.NotificationAPIHandler = (*NotificationAPIHandler)(nil)

func NewNotificationAPIHandler(notificationStore domain.NotificationRepository, aliYunSms client.SmsSend, signName string) *NotificationAPIHandler {
	return &NotificationAPIHandler{
		notificationStore: notificationStore,
		aliYunSms:         aliYunSms,
		signName:          signName,
	}
}
