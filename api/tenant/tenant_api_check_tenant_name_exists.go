package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 查询商户名是否存在
func (s *TenantAPIHandler) CheckTenantNameExists(ctx context.Context, req *pb.CheckTenantNameExistsRequest, rsp *pb.CheckTenantNameExistsResponse) error {
	err := validateCheckTenantNameExistsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	checkRsp, err := s.userAPI.CheckTenantNameExists(ctx, &userv1.CheckTenantNameExistsRequest{
		TenantName: req.GetTenantName(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.IsExists = checkRsp.GetIsExists()

	return nil
}

// 验证request
func validateCheckTenantNameExistsRequest(req *pb.CheckTenantNameExistsRequest) error {
	if req.GetTenantName() == "" {
		return gerr.New("tenant_name should not be empty")
	}
	return nil
}
