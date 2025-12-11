package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) UpdateStaffNickname(ctx context.Context, req *pb.UpdateStaffNicknameRequest, rsp *pb.UpdateStaffNicknameResponse) error {
	err := validateUpdateStaffNicknameRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 发送更新请求
	_, err = s.userAPI.UpdateStaffNickname(ctx, &userpb.UpdateStaffNicknameRequest{
		TenantId:    req.GetTenantId(),
		StaffId:     req.GetStaffId(),
		NewNickname: req.GetNewNickname(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateUpdateStaffNicknameRequest(req *pb.UpdateStaffNicknameRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetNewNickname() == "" {
		return gerr.New("nickname should not be empty")
	}
	return nil
}
