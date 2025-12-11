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

// 员工离职
func (s *BossAPIHandler) DeleteSystemUser(ctx context.Context, req *pb.DeleteSystemUserRequest, rsp *pb.DeleteSystemUserResponse) error {
	err := validateDeleteSystemUserRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.DeleteStaff(ctx, &userv1.DeleteStaffRequest{
		StaffId: req.GetUserId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

// 验证request
func validateDeleteSystemUserRequest(req *pb.DeleteSystemUserRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	return nil
}
