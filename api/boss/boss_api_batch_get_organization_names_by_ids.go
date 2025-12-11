package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"

	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *BossAPIHandler) BatchGetOrganizationNamesByIDs(ctx context.Context, req *pb.BatchGetOrganizationNamesByIDsRequest, rsp *pb.BatchGetOrganizationNamesByIDsResponse) error {
	err := validateBatchGetOrganizationNamesByIDsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 发送请求
	searchRsp, err := s.userAPI.BatchGetOrganizationNamesByIDs(ctx, &userpb.BatchGetOrganizationNamesByIDsRequest{
		// 组织ID
		OrganizationId: req.GetOrganizationId(),
	})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.Organizations = searchRsp.GetOrganizations()
	return nil
}

// 验证request
func validateBatchGetOrganizationNamesByIDsRequest(req *pb.BatchGetOrganizationNamesByIDsRequest) error {
	if req.GetOrganizationId() == nil {
		return gerr.New("organization_id should not be nil")
	}
	return nil
}
