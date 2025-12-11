package notification

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/svc/notification/client"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *NotificationAPIHandler) SendBillingNotification(ctx context.Context, req *pb.SendBillingNotificationRequest, rsp *pb.SendBillingNotificationResponse) error {
	err := validateSendBillingNotificationRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取手机号
	var phone []string
	index := 0
	// 手机数量上限 100 个，超过 100 分批发送
	for _, v := range req.GetNotifications() {
		phone = append(phone, v.GetPhone())
		index++
		// 上限 100
		if index == AliPhoneLimit {
			index = 0
			_, err = s.aliYunSms.SendNotification(phone, s.signName, client.TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING, nil)
			if err != nil {
				return errors.Error(codes.InvalidOperation, err.Error())
			}
			// 初始化 phone
			phone = []string{}
		}
	}

	if len(phone) != 0 {
		_, err = s.aliYunSms.SendNotification(phone, s.signName, client.TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING, nil)
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
	}

	return nil
}

// 验证 request
func validateSendBillingNotificationRequest(req *pb.SendBillingNotificationRequest) error {
	if req.GetNotifications() == nil {
		return gerr.New("notifications should not be nil")
	}
	if len(req.GetNotifications()) == 0 {
		return gerr.New("notifications should not be nil")
	}
	return nil
}
