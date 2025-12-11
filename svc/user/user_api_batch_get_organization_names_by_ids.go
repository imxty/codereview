package user

import (
	"context"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 通过组织 ID 批量获取组织名
func (u *UserAPIHandler) BatchGetOrganizationNamesByIDs(ctx context.Context, req *pb.BatchGetOrganizationNamesByIDsRequest, rsp *pb.BatchGetOrganizationNamesByIDsResponse) error {
	if req.GetOrganizationId() == nil {
		rsp.Organizations = make(map[string]string)
		return nil
	}

	// 获取组织信息
	organizations, err := u.userStore.BatchGetOrganizationNamesByIDs(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// map<string, string> 组织 ID 对应组织名
	results := make(map[string]string)
	for _, each := range organizations {
		results[each.GetOrganizationID()] = each.GetUsername()
	}

	// 返回响应
	rsp.Organizations = results

	return nil
}
