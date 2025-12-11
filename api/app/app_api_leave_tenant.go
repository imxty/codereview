package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
)

func (s *AppAPIHandler) LeaveTenant(ctx context.Context, req *pb.LeaveTenantRequest, rsp *pb.LeaveTenantResponse) error {

	// 1.验证request
	err := validateLeaveTenantRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 发送离开当前租户的请求
	_, err = s.userAPI.LeaveTenant(ctx, &userpb.LeaveTenantRequest{
		StaffId:  req.GetStaffId(),
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateLeaveTenantRequest(req *pb.LeaveTenantRequest) error {
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
