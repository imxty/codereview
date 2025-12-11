package review

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/svc/review/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取审核通知列表
func (s *ReviewAPIHandler) ListReviewNotifications(ctx context.Context, req *pb.ListReviewNotificationsRequest, rsp *pb.ListReviewNotificationsResponse) error {
	err := validateListReviewNotificationsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询通知
	// 查询组织通知
	var notifications []domain.ReviewNotificationIntf
	if req.GetOrganizationId() != "" {
		notifications, err = s.reviewStore.ListOrganizationNotifications(ctx, req.GetOrganizationId())
		if err != nil {
			return errors.Error(codes.InvalidRequest, err.Error())
		}
	} else {
		notifications, err = s.reviewStore.ListTenantNotifications(ctx, req.GetTenantId())
		if err != nil {
			return errors.Error(codes.InvalidRequest, err.Error())
		}
	}

	// 转换结构
	pNotifications := make([]*pb.ReviewNotification, len(notifications))
	for k, v := range notifications {
		pNotifications[k] = toProtoNotification(v)
	}

	rsp.ReviewNotifications = pNotifications
	return nil
}

// 验证request
func validateListReviewNotificationsRequest(req *pb.ListReviewNotificationsRequest) error {
	if req.GetOrganizationId() == "" && req.GetTenantId() == "" {
		return gerr.New("organization id and tenant id should not be empty")
	}
	return nil
}
