package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) ResetAdminPassword(ctx context.Context, req *pb.ResetAdminPasswordRequest, rsp *pb.ResetAdminPasswordResponse) error {
	err := validateResetAdminPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 检测密码强度
	ok := checkPassword(req.GetNewPlainPassword())
	if !ok {
		return errors.Error(codes.InvalidRequest, "ez password")
	}
	// 发送重置商户密码请求
	_, err = s.userAPI.ResetAdminPassword(ctx, &userpb.ResetAdminPasswordRequest{
		Phone:            req.GetPhone(),
		SmsTxId:          req.GetSmsTxId(),
		NewPlainPassword: req.GetNewPlainPassword(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateResetAdminPasswordRequest(req *pb.ResetAdminPasswordRequest) error {
	if req.GetNewPlainPassword() == "" {
		return gerr.New("new_plain_password should not be empty")
	}
	if req.GetSmsTxId() == "" {
		return gerr.New("sms_tx_id should not be empty")
	}
	if req.GetPhone() == "" {
		return gerr.New("invalid tenant_id")
	}
	return nil
}

// 检测密码是否满足强度要求
func checkPassword(pass string) bool {
	if len(pass) < 8 || len(pass) > 16 {
		return false
	}
	return true
}
