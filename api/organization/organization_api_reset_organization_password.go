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

// 重置密码
func (s *OrganizationAPIHandler) ResetOrganizationPassword(ctx context.Context, req *pb.ResetOrganizationPasswordRequest, rsp *pb.ResetOrganizationPasswordResponse) error {
	err := validateResetOrganizationPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.userAPI.ResetOrganizationPassword(ctx, &userv1.ResetOrganizationPasswordRequest{
		// 手机号
		Phone: req.GetPhone(),
		// 新明文密码
		NewPlainPassword: req.GetNewPlainPassword(),
		// 凭证 ID
		TxId: req.GetTxId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	return nil
}

// 验证request
func validateResetOrganizationPasswordRequest(req *pb.ResetOrganizationPasswordRequest) error {
	if req.GetNewPlainPassword() == "" {
		return gerr.New("new plain password should not be empty")
	}
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx_id should not be empty")
	}
	return nil
}
