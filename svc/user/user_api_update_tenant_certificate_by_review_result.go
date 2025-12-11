package user

import (
	"context"
	gerr "errors"

	notificationv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 修改商户审核结果
func (u *UserAPIHandler) UpdateTenantCertificateByReviewResult(ctx context.Context, req *pb.UpdateTenantCertificateByReviewResultRequest, rsp *pb.UpdateTenantCertificateByReviewResultResponse) error {
	err := validateUpdateTenantCertificateByReviewResult(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取最新的提审信息
	lastRevision, err := u.userStore.GetLatestTenantEntityRevision(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 获取组织商户关系
	organizationTenant, err := u.userStore.GetOrganizationTenantByTenantId(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if organizationTenant == nil {
		return errors.Errorf(ErrTenantNotFound, "tenant[%s] not found", req.GetTenantId())
	}

	// 如果审核失败，更新提审信息状态并添加失败原因
	if !req.GetStatus() {
		ctx = u.userStore.BeginTx(ctx)
		// 添加失败原因
		err = u.userStore.UpdateTenantRevisionFailReason(ctx, lastRevision.GetTenantEntityRevisionID(), req.GetFailReason())
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 修改状态为审核失败
		err = u.userStore.UpdateOrganizationTenantStatus(ctx, req.GetTenantId(), domain.TenantReviewStatusReviewFailed, organizationTenant.GetRev())
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}

		u.userStore.CommitTx(ctx)

		// 如果是首次认证失败，发送短信通知
		if !organizationTenant.GetIsActivated() {
			// 异步发送短信，忽略错误
			notificationCtx := context.Background()
			go u.notificationAPI.SendVerificationCode(notificationCtx, &notificationv1.SendVerificationCodeRequest{
				Phone:          lastRevision.GetSafePhone(),
				TemplateAction: notificationv1.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL,
				Language:       notificationv1.Language_LANGUAGE_SIMPLIFIED_CHINESE,
			})
		}
		return nil
	}

	// 如果审核成功，修改商户审核状态并更新商户实体信息
	ctx = u.userStore.BeginTx(ctx)

	// 修改商户审核状态
	if organizationTenant.GetReviewStatus() != domain.TenantReviewStatusReviewSuccess {
		// 如果当前是未认证状态，将状态变更为待续期
		err = u.userStore.UpdateOrganizationTenantStatus(ctx, req.GetTenantId(), domain.TenantReviewStatusReviewSuccess, organizationTenant.GetRev())
		if err != nil {
			u.userStore.RollbackTx(ctx)
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	// 更新商户实体信息
	err = u.userStore.UpdateTenantEntity(ctx, lastRevision, lastRevision.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	u.userStore.CommitTx(ctx)

	// 如果是首次认证成功，发送短信通知
	if !organizationTenant.GetIsActivated() {
		// 异步发送短信，忽略错误
		notificationCtx := context.Background()
		go u.notificationAPI.SendVerificationCode(notificationCtx, &notificationv1.SendVerificationCodeRequest{
			Phone:          lastRevision.GetSafePhone(),
			TemplateAction: notificationv1.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS,
			Language:       notificationv1.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		})
	}

	return nil
}

// 验证 request
func validateUpdateTenantCertificateByReviewResult(req *pb.UpdateTenantCertificateByReviewResultRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	return nil
}
