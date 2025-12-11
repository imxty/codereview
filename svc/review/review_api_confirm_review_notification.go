package review

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrNotificationNotFound 通知不存在
	ErrNotificationNotFound = 5038
)

func (s *ReviewAPIHandler) ConfirmReviewNotification(ctx context.Context, req *pb.ConfirmReviewNotificationRequest, rsp *pb.ConfirmReviewNotificationResponse) error {
	// 1.验证request
	err := validateConfirmReviewNotificationRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取通知
	notification, err := s.reviewStore.GetReviewNotification(ctx, req.GetNotificationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if notification == nil {
		return errors.Errorf(ErrNotificationNotFound, "review notification[%s] not found", req.GetNotificationId())
	}

	// 确认审核通知
	err = s.reviewStore.ConfirmReviewNotification(ctx, req.GetNotificationId(), notification.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
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
