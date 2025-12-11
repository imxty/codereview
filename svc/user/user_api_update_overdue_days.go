package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 修改复查天数
func (u *UserAPIHandler) UpdateOverdueDays(ctx context.Context, req *pb.UpdateOverdueDaysRequest, rsp *pb.UpdateOverdueDaysResponse) error {
	err := validateUpdateOverdueDaysRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商户
	tenant, err := u.userStore.GetTenantEntity(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tenant == nil {
		return errors.Errorf(ErrTenantNotFound, "tenant[%s] not found", req.GetTenantId())
	}

	// 更新商户复查天数
	err = u.userStore.UpdateTenantReview(ctx, tenant.GetTenantID(), req.GetDays())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	rsp.Days = req.GetDays()

	return nil
}

// 验证 request
func validateUpdateOverdueDaysRequest(req *pb.UpdateOverdueDaysRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetDays() < 0 {
		return gerr.New("invalid days")
	}
	return nil
}
