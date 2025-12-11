package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

// SubmitFeedback 提交反馈
func (u *UserAPIHandler) SubmitFeedback(ctx context.Context, req *pb.SubmitFeedbackRequest, rsp *pb.SubmitFeedbackResponse) error {
	err := validateSubmitFeedback(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	err = u.userStore.SubmitFeedback(ctx, &domain.Feedback{
		FeedbackID: xid.New().String(),
		TenantID:   req.GetTenantId(),
		Phone:      req.GetPhone(),
		Content:    req.GetContent(),
	})
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

func validateSubmitFeedback(req *pb.SubmitFeedbackRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetContent() == "" {
		return gerr.New("content should not be empty")
	}
	return nil
}
