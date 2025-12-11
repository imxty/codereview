package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 根据组织名查询组织名
func (u *UserAPIHandler) SearchOrganizationByName(ctx context.Context, req *pb.SearchOrganizationByNameRequest, rsp *pb.SearchOrganizationByNameResponse) error {
	err := validateSearchOrganizationByNameRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询组织信息
	organizations, err := u.userStore.SearchOrganizationByName(ctx, req.GetOrganizationName())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 查询所有组织
	if req.GetPrivilegeId() == "all" {
		results := make(map[string]string)
		for _, each := range organizations {
			results[each.GetOrganizationID()] = each.GetUsername()
		}

		rsp.Organizations = results

		return nil
	}

	// 查询数据权限
	ps, err := u.userStore.GetDataPrivilegeByID(ctx, req.GetPrivilegeId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 如果是所有权限
	if len(ps) == 1 && ps[0] == "all" {
		results := make(map[string]string)
		for _, each := range organizations {
			results[each.GetOrganizationID()] = each.GetUsername()
		}

		rsp.Organizations = results
		return nil
	}

	// 不是所有权限
	p_map := make(map[string]bool)
	for _, p := range ps {
		p_map[p] = true
	}

	// 筛选数据权限内的组织
	results := make(map[string]string)
	for _, each := range organizations {
		if _, ok := p_map[each.GetOrganizationID()]; ok {
			results[each.GetOrganizationID()] = each.GetUsername()
		}
	}

	// 返回响应
	rsp.Organizations = results

	return nil
}

// 验证 request
func validateSearchOrganizationByNameRequest(req *pb.SearchOrganizationByNameRequest) error {
	if req.GetPrivilegeId() == "" {
		return gerr.New("privilege_id should not be empty")
	}
	return nil
}
