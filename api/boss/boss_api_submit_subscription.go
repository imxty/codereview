package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 提交订阅
func (s *BossAPIHandler) SubmitSubscription(ctx context.Context, req *pb.SubmitSubscriptionRequest, rsp *pb.SubmitSubscriptionResponse) error {
	err := validateSubmitSubscriptionRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.SubmitSubscription(ctx, &userv1.SubmitSubscriptionRequest{
		TenantId:  req.GetTenantId(),
		UserId:    req.GetUserId(),
		TimeCount: req.GetTimeCount(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

// 验证request
func validateSubmitSubscriptionRequest(req *pb.SubmitSubscriptionRequest) error {
	if req.GetTenantId() == nil {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	return nil
}
