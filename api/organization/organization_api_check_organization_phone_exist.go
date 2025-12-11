package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 判断组织手机号是否存在
func (s *OrganizationAPIHandler) CheckOrganizationPhoneExist(ctx context.Context, req *pb.CheckOrganizationPhoneExistRequest, rsp *pb.CheckOrganizationPhoneExistResponse) error {
	err := validateCheckOrganizationPhoneExistRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 检查手机号是否存在
	checkRsp, err := s.userAPI.CheckOrganizationPhoneExist(ctx, &userv1.CheckOrganizationPhoneExistRequest{
		Phone: req.GetPhone(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回响应
	rsp.Exist = checkRsp.GetExist()

	return nil
}

// 验证request
func validateCheckOrganizationPhoneExistRequest(req *pb.CheckOrganizationPhoneExistRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	return nil
}
