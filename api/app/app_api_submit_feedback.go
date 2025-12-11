package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) SubmitFeedback(ctx context.Context, req *pb.SubmitFeedbackRequest, rsp *pb.SubmitFeedbackResponse) error {
	err := validateSubmitFeedbackRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.SubmitFeedback(ctx, &userv1.SubmitFeedbackRequest{
		TenantId: req.GetTenantId(),
		Phone:    req.GetPhone(),
		Content:  req.GetContent(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

func validateSubmitFeedbackRequest(req *pb.SubmitFeedbackRequest) error {
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
