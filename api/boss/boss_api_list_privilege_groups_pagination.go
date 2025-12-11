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

func (s *BossAPIHandler) ListPrivilegeGroupsPagination(ctx context.Context, req *pb.ListPrivilegeGroupsPaginationRequest, rsp *pb.ListPrivilegeGroupsPaginationResponse) error {
	err := validateListPrivilegeGroupsPaginationRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	groupsRsp, err := s.userAPI.ListPrivilegeGroupsPagination(ctx, &userv1.ListPrivilegeGroupsPaginationRequest{
		Pagination: toUserPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 转换权限列表
	apiGroups := make([]*pb.PrivilegeGroup, len(groupsRsp.GetPrivilegeGroups()))
	for k, v := range groupsRsp.GetPrivilegeGroups() {
		apiGroups[k] = toApiPrivilegeGroup(v)
	}

	rsp.PrivilegeGroups = apiGroups
	rsp.TotalCount = groupsRsp.GetTotalCount()

	return nil
}

func validateListPrivilegeGroupsPaginationRequest(req *pb.ListPrivilegeGroupsPaginationRequest) error {
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
