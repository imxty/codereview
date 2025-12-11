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
func (s *BossAPIHandler) DeletePrivilegeGroup(ctx context.Context, req *pb.DeletePrivilegeGroupRequest, rsp *pb.DeletePrivilegeGroupResponse) error {
	err := validateDeletePrivilegeGroupRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.DeletePrivilegeGroup(ctx, &userv1.DeletePrivilegeGroupRequest{
		UserId:      req.GetUserId(),
		PrivilegeId: req.GetPrivilegeId(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

func validateDeletePrivilegeGroupRequest(req *pb.DeletePrivilegeGroupRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
