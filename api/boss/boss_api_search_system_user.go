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

// 查询系统用户
func (s *BossAPIHandler) SearchSystemUser(ctx context.Context, req *pb.SearchSystemUserRequest, rsp *pb.SearchSystemUserResponse) error {
	err := validateSearchSystemUserRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	searchRsp, err := s.userAPI.SearchSystemUser(ctx, &userv1.SearchSystemUserRequest{
		Nickname:           req.GetNickname(),
		Phone:              req.GetPhone(),
		PrivilegeGroupName: req.GetPrivilegeGroupName(),
		IsActivated:        req.GetIsActivated(),
		Pagination:         toUserPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 转换
	users := make([]*pb.SystemUser, len(searchRsp.GetUsers()))
	for k, v := range searchRsp.GetUsers() {
		users[k] = toApiSystemUser(v)
	}

	// 返回响应
	rsp.Users = users
	rsp.TotalCount = searchRsp.GetTotalCount()

	return nil
}

// 验证request
func validateSearchSystemUserRequest(req *pb.SearchSystemUserRequest) error {
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	return nil
}
