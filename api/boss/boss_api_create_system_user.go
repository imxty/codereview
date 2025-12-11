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

// 创建系统用户
func (s *BossAPIHandler) CreateSystemUser(ctx context.Context, req *pb.CreateSystemUserRequest, rsp *pb.CreateSystemUserResponse) error {
	err := validateCreateSystemUserRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.CreateSystemUser(ctx, &userv1.CreateSystemUserRequest{
		StaffName:   req.GetStaffName(),
		StaffPhone:  req.GetStaffPhone(),
		PrivilegeId: req.GetPrivilegeId(),
		Remark:      req.GetRemark(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

// 验证request
func validateCreateSystemUserRequest(req *pb.CreateSystemUserRequest) error {
	if req.GetStaffName() == "" {
		return gerr.New("staff_name should not be empty")
	}
	if req.GetStaffPhone() == "" {
		return gerr.New("staff_phone should not be empty")
	}
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
