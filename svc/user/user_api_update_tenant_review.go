package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 修改复查天数
func (u *UserAPIHandler) UpdateTenantReview(ctx context.Context, req *pb.UpdateTenantReviewRequest, rsp *pb.UpdateTenantReviewResponse) error {
	err := validateUpdateTenantReviewRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 修改复查天数
	err = u.userStore.UpdateTenantReview(ctx, req.GetTenantId(), req.GetDays())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证 request
func validateUpdateTenantReviewRequest(req *pb.UpdateTenantReviewRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetDays() < 1 {
		return gerr.New("invalid days")
	}
	return nil
}
