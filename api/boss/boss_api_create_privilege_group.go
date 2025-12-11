package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 创建权限组
func (s *BossAPIHandler) CreatePrivilegeGroup(ctx context.Context, req *pb.CreatePrivilegeGroupRequest, rsp *pb.CreatePrivilegeGroupResponse) error {
	err := validateCreatePrivilegeGroupRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	pagePrivileges := make([]userv1.PagePrivilege, len(req.GetPagePrivileges()))
	for k, v := range req.GetPagePrivileges() {
		pagePrivileges[k] = toSvcPagePrivileges(v)
	}

	_, err = s.userAPI.CreatePrivilegeGroup(ctx, &userv1.CreatePrivilegeGroupRequest{
		UserId:         req.GetUserId(),
		PrivilegeName:  req.GetPrivilegeName(),
		Remark:         req.GetRemark(),
		PagePrivileges: pagePrivileges,
		DataPrivileges: req.GetDataPrivileges(),
	})
	if err != nil {
		return err
	}

	return nil
}

// 验证request
func validateCreatePrivilegeGroupRequest(req *pb.CreatePrivilegeGroupRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetPrivilegeName() == "" {
		return gerr.New("privilege_name should not be empty")
	}
	return nil
}
