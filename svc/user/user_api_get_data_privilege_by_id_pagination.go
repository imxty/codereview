package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (u *UserAPIHandler) GetDataPrivilegeByIDPagination(ctx context.Context, req *pb.GetDataPrivilegeByIDPaginationRequest, rsp *pb.GetDataPrivilegeByIDPaginationResponse) error {
	err := validateGetDataPrivilegeByIDPagination(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 分页获取数据权限
	os, count, err := u.userStore.GetDataPrivilegePagination(ctx, req.GetPrivilegeId(), req.GetPagination().GetSize(), req.GetPagination().GetOffset())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获取组织名
	onames, err := u.userStore.BatchGetOrganizationNamesByIDs(ctx, os)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	nameMap := make(map[string]string)
	for _, each := range onames {
		nameMap[each.GetOrganizationID()] = each.GetUsername()
	}

	rsp.Datas = nameMap
	rsp.TotalCount = count

	return nil
}

func validateGetDataPrivilegeByIDPagination(req *pb.GetDataPrivilegeByIDPaginationRequest) error {
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
