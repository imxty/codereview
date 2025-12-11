package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取组织下所有认证商户
func (s *UserAPIHandler) ListPrivilegeGroupsPagination(ctx context.Context, req *pb.ListPrivilegeGroupsPaginationRequest, rsp *pb.ListPrivilegeGroupsPaginationResponse) error {
	err := validateListPrivilegeGroupsPaginationRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	groups, count, err := s.userStore.ListPrivilegeGroupsPagination(ctx, req.GetPagination().GetSize(), req.GetPagination().GetOffset())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回响应
	results := make([]*pb.PrivilegeGroup, len(groups))
	for i, each := range groups {
		results[i] = toProtoPrivilegeGroup(each)
	}
	rsp.PrivilegeGroups = results
	rsp.TotalCount = int32(count)

	return nil
}

func validateListPrivilegeGroupsPaginationRequest(req *pb.ListPrivilegeGroupsPaginationRequest) error {
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
