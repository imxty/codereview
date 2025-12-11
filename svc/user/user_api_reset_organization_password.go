package user

import (
	"context"
	gerr "errors"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrInvalidTxId
	ErrInvalidTxId = 5021
)

func (u *UserAPIHandler) ResetOrganizationPassword(ctx context.Context, req *pb.ResetOrganizationPasswordRequest, rsp *pb.ResetOrganizationPasswordResponse) error {
	err := validateResetOrganizationPasswordRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询组织
	organization, err := u.userStore.GetOrganizationByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if organization == nil {
		return errors.Errorf(ErrOrganizationNotExist, "organization[%s] not found", req.GetPhone())
	}

	// 验证手机号
	checkRsp, err := u.notificationAPI.CheckTxIdUsage(ctx, &notificationpb.CheckTxIdUsageRequest{
		TxId: req.GetTxId(),
	})
	if err != nil {
		return err
	}
	if checkRsp.GetAction() != notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD {
		return errors.Errorf(ErrInvalidTxId, "tx id[%s] is not aim to reset organization password", req.GetTxId())
	}

	// 修改密码
	hashedPassword := generateHashSHA256(req.GetNewPlainPassword())
	err = u.userStore.UpdateOrganizationPassword(ctx, organization.GetOrganizationID(), hashedPassword, organization.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证 request
func validateResetOrganizationPasswordRequest(req *pb.ResetOrganizationPasswordRequest) error {
	if req.GetNewPlainPassword() == "" {
		return gerr.New("password should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	return nil
}
