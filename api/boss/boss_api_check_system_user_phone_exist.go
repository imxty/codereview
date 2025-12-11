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

// 验证手机号是否存在
func (s *BossAPIHandler) CheckSystemUserPhoneExist(ctx context.Context, req *pb.CheckSystemUserPhoneExistRequest, rsp *pb.CheckSystemUserPhoneExistResponse) error {
	err := validateCheckSystemUserPhoneExistRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 验证手机号
	checkRsp, err := s.userAPI.CheckSystemUserPhoneExist(ctx, &userv1.CheckSystemUserPhoneExistRequest{
		Phone: req.GetPhone(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.Exist = checkRsp.GetExist()

	return nil
}

// 验证 request
func validateCheckSystemUserPhoneExistRequest(req *pb.CheckSystemUserPhoneExistRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	return nil
}
