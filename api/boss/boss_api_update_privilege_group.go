package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// 创建系统用户
func (s *BossAPIHandler) UpdatePrivilegeGroup(ctx context.Context, req *pb.UpdatePrivilegeGroupRequest, rsp *pb.UpdatePrivilegeGroupResponse) error {
	err := validateUpdatePrivilegeGroupRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	pagePrivileges := make([]userv1.PagePrivilege, len(req.GetPagePrivileges()))
	for k, v := range req.GetPagePrivileges() {
		pagePrivileges[k] = toSvcPagePrivileges(v)
	}

	_, err = s.userAPI.UpdatePrivilegeGroup(ctx, &userv1.UpdatePrivilegeGroupRequest{
		UserId:         req.GetUserId(),
		PrivilegeName:  req.GetPrivilegeName(),
		Remark:         req.GetRemark(),
		PagePrivileges: pagePrivileges,
		DataPrivileges: req.GetDataPrivileges(),
		PrivilegeId:    req.GetPrivilegeId(),
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	searchRsp, err := s.userAPI.SearchSystemUser(ctx, &userv1.SearchSystemUserRequest{
		Nickname:           "",
		Phone:              "",
		PrivilegeGroupName: req.GetPrivilegeName(),
		IsActivated:        &wrapperspb.BoolValue{Value: true},
		Pagination: &userv1.Pagination{
			Offset: 0,
			Size:   2147483647,
		},
	})

	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	kickUsers := make([]string, 0)
	for _, each := range searchRsp.GetUsers() {
		if each.GetPrivilegeGroup().GetPrivilegeId() == req.GetPrivilegeId() {
			kickUsers = append(kickUsers, each.GetUserId())
		}
	}

	err = s.tokenStore.KickOutUsers(ctx, kickUsers)
	if err != nil {
		return errors.Error(api.ErrCreateAccessToken, api.ErrorChineseMsg(api.ErrKickedOut))
	}

	return nil
}

func validateUpdatePrivilegeGroupRequest(req *pb.UpdatePrivilegeGroupRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetPrivilegeName() == "" {
		return gerr.New("privilege_name should not be empty")
	}
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
