package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	reviewv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取审核通知列表
func (s *OrganizationAPIHandler) ListReviewNotifications(ctx context.Context, req *pb.ListReviewNotificationsRequest, rsp *pb.ListReviewNotificationsResponse) error {
	err := validateListReviewNotificationsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 发送请求
	listRsp, err := s.reviewAPI.ListReviewNotifications(ctx, &reviewv1.ListReviewNotificationsRequest{
		OrganizationId: req.GetOrganizationId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	notifications := listRsp.GetReviewNotifications()
	// 转化结构
	appNotifications := make([]*pb.ReviewNotification, len(notifications))
	for k, v := range notifications {
		appNotifications[k] = toAppNotification(v)
	}

	rsp.ReviewNotifications = appNotifications

	return nil
}

// 验证request
func validateListReviewNotificationsRequest(req *pb.ListReviewNotificationsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	return nil
}
