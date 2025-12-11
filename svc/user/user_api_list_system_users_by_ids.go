package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 通过用户 ID 获取系统用户信息
func (s *UserAPIHandler) ListSystemUsersByIds(ctx context.Context, req *pb.ListSystemUsersByIdsRequest, rsp *pb.ListSystemUsersByIdsResponse) error {
	err := validateListSystemUsersByIdsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取系统用户
	sus, err := s.userStore.ListSystemUsers(ctx, req.GetSystemUserIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	suMap := make(map[string]*pb.SystemUser)
	for _, v := range sus {
		suMap[v.GetUserID()] = toAppSystemUser(v)
	}
	rsp.SystemUsers = suMap
	return nil
}

func validateListSystemUsersByIdsRequest(req *pb.ListSystemUsersByIdsRequest) error {
	if len(req.GetSystemUserIds()) == 0 {
		return gerr.New("system user ids should not be nil")
	}
	return nil
}
