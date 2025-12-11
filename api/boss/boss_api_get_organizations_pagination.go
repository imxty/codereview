package boss

import (
	"context"
	gerr "errors"
	"sort"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 创建系统用户
func (s *BossAPIHandler) GetOrganizationsPagination(ctx context.Context, req *pb.GetOrganizationsPaginationRequest, rsp *pb.GetOrganizationsPaginationResponse) error {
	err := validateGetOrganizationsPagination(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 搜索组织
	searchRsp, err := s.userAPI.SearchOrganizationByName(ctx, &userv1.SearchOrganizationByNameRequest{
		OrganizationName: "",
		PrivilegeId:      "all",
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	organizations := searchRsp.GetOrganizations()

	// 返回分页数据
	pageData := make(map[string]string)
	keys := make([]string, len(pageData))

	for k := range organizations {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	start := int(req.GetPagination().GetOffset())
	end := int(req.GetPagination().GetOffset() + req.GetPagination().GetSize())

	if end > len(organizations) {
		end = len(organizations)
	}

	for i := start; i < end; i++ {
		key := keys[i]
		pageData[key] = organizations[key]
	}

	rsp.TotalCount = int32(len(organizations))
	rsp.Organizations = pageData

	return nil
}

func validateGetOrganizationsPagination(req *pb.GetOrganizationsPaginationRequest) error {
	if req.GetUserId() == "" {
		return gerr.New("user_id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
