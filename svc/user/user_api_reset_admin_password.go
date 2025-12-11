package user

import (
	"context"
	gerr "errors"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 重置商户密码
func (u *UserAPIHandler) ResetAdminPassword(ctx context.Context, req *pb.ResetAdminPasswordRequest, rsp *pb.ResetAdminPasswordResponse) error {
	// 验证 request
	err := validateResetAdminPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询 entity_user
	entityUser, err := u.userStore.GetTenantUserByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if entityUser == nil {
		return errors.Errorf(ErrTenantNotFound, "entity[%s] not found", req.GetPhone())
	}

	// 验证手机号
	checkRsp, err := u.notificationAPI.CheckTxIdUsage(ctx, &notificationpb.CheckTxIdUsageRequest{
		TxId: req.GetSmsTxId(),
	})
	if err != nil {
		return err
	}
	if checkRsp.GetAction() != notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD {
		return errors.Errorf(ErrInvalidTxId, "tx id[%s] is not aim to reset tenant password", req.GetSmsTxId())
	}

	// 修改密码
	hashedPassword := generateHashSHA256(req.GetNewPlainPassword())
	err = u.userStore.UpdateTenantPassword(ctx, entityUser.GetTenantID(), hashedPassword, entityUser.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证 request
func validateResetAdminPasswordRequest(req *pb.ResetAdminPasswordRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("tenant id should not be empty")
	}
	if req.GetSmsTxId() == "" {
		return gerr.New("sms tx id should not be empty")
	}
	if req.GetNewPlainPassword() == "" {
		return gerr.New("password should not be empty")
	}
	return nil
}
