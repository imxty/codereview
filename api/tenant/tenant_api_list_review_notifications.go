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

// 获取审核通知列表
func (s *TenantAPIHandler) ListReviewNotifications(ctx context.Context, req *pb.ListReviewNotificationsRequest, rsp *pb.ListReviewNotificationsResponse) error {
	err := validateListReviewNotificationsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 发送获取审核通知列表请求
	listRsp, err := s.reviewAPI.ListReviewNotifications(ctx, &reviewv1.ListReviewNotificationsRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 转化结构
	notifications := listRsp.GetReviewNotifications()
	appNotifications := make([]*pb.ReviewNotification, len(notifications))
	for k, v := range notifications {
		appNotifications[k] = toAppNotification(v)
	}

	// 返回响应
	rsp.ReviewNotifications = appNotifications

	return nil
}

// 验证request
func validateListReviewNotificationsRequest(req *pb.ListReviewNotificationsRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
