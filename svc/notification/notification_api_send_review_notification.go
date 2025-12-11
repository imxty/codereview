package notification

import (
	"context"
	gerr "errors"
	"strconv"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/svc/notification/client"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	AliPhoneLimit = 100
)

func (s *NotificationAPIHandler) SendReviewNotification(ctx context.Context, req *pb.SendReviewNotificationRequest, rsp *pb.SendReviewNotificationResponse) error {
	err := validateSendReviewNotificationRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	params := make(map[string]string)
	params["issue"] = strconv.Itoa(int(req.GetIssueNumber()))
	// 获取手机号
	var phone []string
	index := 0
	// 手机数量上限 100 个，超过 100 分批发送
	for _, v := range req.GetPhone() {
		phone = append(phone, v)
		index++
		// 上限 100
		if index == AliPhoneLimit {
			index = 0
			_, err = s.aliYunSms.SendNotification(phone, s.signName, client.TemplateAction_TEMPLATE_ACTION_REVIEW_NOTIFICATION, params)
			if err != nil {
				return errors.Error(codes.InvalidOperation, err.Error())
			}
			// 初始化 phone
			phone = []string{}
		}
	}

	if len(phone) != 0 {
		_, err = s.aliYunSms.SendNotification(phone, s.signName, client.TemplateAction_TEMPLATE_ACTION_REVIEW_NOTIFICATION, params)
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
	}

	return nil
}

// 验证 request
func validateSendReviewNotificationRequest(req *pb.SendReviewNotificationRequest) error {
	if len(req.GetPhone()) == 0 {
		return gerr.New("phone should not be nil")
	}
	return nil
}
