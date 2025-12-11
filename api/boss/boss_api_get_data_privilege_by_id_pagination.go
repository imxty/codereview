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

func (s *BossAPIHandler) GetDataPrivilegeByIDPagination(ctx context.Context, req *pb.GetDataPrivilegeByIDPaginationRequest, rsp *pb.GetDataPrivilegeByIDPaginationResponse) error {
	err := validateGetDataPrivilegeByIDPaginationRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.userAPI.GetDataPrivilegeByIDPagination(ctx, &userv1.GetDataPrivilegeByIDPaginationRequest{
		PrivilegeId: req.GetPrivilegeId(),
		Pagination:  toUserPagination(req.GetPagination()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.Datas = getRsp.GetDatas()
	rsp.TotalCount = getRsp.GetTotalCount()

	return nil
}

func validateGetDataPrivilegeByIDPaginationRequest(req *pb.GetDataPrivilegeByIDPaginationRequest) error {
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
