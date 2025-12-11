package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	reviewv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) ConfirmReviewNotification(ctx context.Context, req *pb.ConfirmReviewNotificationRequest, rsp *pb.ConfirmReviewNotificationResponse) error {

	err := validateConfirmReviewNotificationRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	_, err = s.reviewAPI.ConfirmReviewNotification(ctx, &reviewv1.ConfirmReviewNotificationRequest{
		NotificationId: req.GetNotificationId(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateConfirmReviewNotificationRequest(req *pb.ConfirmReviewNotificationRequest) error {
	if req.GetNotificationId() == "" {
		return gerr.New("notification id should not be empty")
	}
	return nil
}
