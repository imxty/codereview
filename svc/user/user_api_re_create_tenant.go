package user

import (
	"context"
	gerr "errors"
	"strings"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

func (u *UserAPIHandler) ReCreateTenant(ctx context.Context, req *pb.ReCreateTenantRequest, rsp *pb.ReCreateTenantResponse) error {
	// 验证request
	err := validateReCreateTenantRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationpb.VerifyPhoneVerificationCodeRequest{
		// 手机号
		Phone: req.GetTenant().GetSafePhone(),
		// 短信验证码
		SmsCode: req.GetSmsCode(),
		// 凭证ID
		TxId: req.GetTxId(),
		// 短信验证码类型
		Action: notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE,
	})
	if err != nil {
		return err
	}

	// 获取租户ID
	tid := req.GetTenant().GetTenantId()

	// 获取商户user
	t, err := u.userStore.GetTenantUser(ctx, tid)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if t == nil {
		return errors.Errorf(ErrTenantNotFound, "tenant[%s] not found", tid)
	}

	ot, err := u.userStore.GetOrganizationTenant(ctx, t.GetOrganizationID(), tid)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if ot == nil {
		return errors.Errorf(ErrOrganizationTenantNotExist, "organization tenant[%s] not found", tid)
	}

	// 如果手机号不同需要验证
	if t.GetPhone() != req.GetTenant().GetContactPhone() {
		// 通过手机号查询商户
		userExist, err := u.userStore.GetTenantUserByPhone(ctx, req.GetTenant().GetSafePhone())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		if userExist != nil {
			// 手机号存在
			return errors.Errorf(ErrPhoneHasBeenUsed, "tenant phone[%s] has been used", req.GetTenant().GetSafePhone())
		}
	}

	// 获取最新的副本更新情况
	revision, err := u.userStore.GetLatestTenantEntityRevision(ctx, tid)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if revision == nil {
		return errors.Errorf(ErrTenantNotFound, "revision[%s] not found", tid)
	}

	// 获取商户资质路径和logo
	businessUrl := ""
	logoUrl := ""
	storeName := req.GetTenant().GetName()
	province := req.GetTenant().GetAddress().GetProvince()
	city := req.GetTenant().GetAddress().GetCity()
	district := req.GetTenant().GetAddress().GetDistrict()
	street := req.GetTenant().GetAddress().GetStreet()
	socialCreditCode := req.GetTenant().GetSocialCreditCode()
	reqTenant := req.GetTenant()

	// 上传营业执照
	// 如果没有传就用原来的
	if reqTenant.GetBusinessLicenseUrl() == "" {
		// 如果有营业执照，上传获取URL
		if reqTenant.GetBusinessLicense() != nil {
			path, err := u.uploadImage(reqTenant.GetBusinessLicense())
			if err != nil {
				return errors.Error(codes.InvalidOperation, err.Error())
			}
			businessUrl = path
		} else {
			businessUrl = revision.GetBusinessLicenseUrl()
		}
	} else {
		businessUrl = checkLogoUrl(reqTenant.GetBusinessLicenseUrl())
		if businessUrl == "" {
			businessUrl = revision.GetBusinessLicenseUrl()
		}
	}

	// 上传logo
	// 如果没有传就用原来的
	logoUrl = checkLogoUrl(req.GetTenant().GetLogoUrl())
	if logoUrl == "" {
		// 如果有，则用最新提交的
		logoUrl = revision.GetLogoUrl()
	}

	// 1.生成商户组织关系
	// 2.生成商户(user)
	// 3.提审
	ctx = u.userStore.BeginTx(ctx)
	// 生成hash密码
	hashedPassword := generateHashSHA256(req.GetPlainPassword())
	// 构建更新数据
	updateMap := map[string]interface{}{
		"nickname":        req.GetTenant().GetContactName(),
		"phone":           req.GetTenant().GetContactPhone(),
		"hashed_password": hashedPassword,
	}
	err = u.userStore.UpdateTenantUser(ctx, tid, updateMap)
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 将状态变为待审核
	err = u.userStore.ChangeOrganizationTenantReviewStatus(ctx, tid, domain.TenantReviewStatusReviewing, ot.GetRev())
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 创建一个新的revision
	// 如果通过就创建新的副本，没有通过，就属于之前的副本组
	// 创建一个租户资质的副本
	entityApproving := &domain.TenantEntityRevision{
		TenantEntityRevisionID: xid.New().String(),
		TenantID:               tid,
		StoreName:              storeName,
		Province:               province,
		City:                   city,
		District:               district,
		Street:                 street,
		SafePhone:              req.GetTenant().GetSafePhone(),
		ContactName:            req.GetTenant().GetContactName(),
		ContactPhone:           req.GetTenant().GetContactPhone(),
		SocialCreditCode:       socialCreditCode,
		BusinessLicenseUrl:     businessUrl,
		LogoUrl:                logoUrl,
	}
	err = u.userStore.CreateTenantEntityRevision(ctx, entityApproving)
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	u.userStore.CommitTx(ctx)

	// 提交审核请求
	_, err = u.reviewAPI.CommitEntityCertificate(ctx, &reviewpb.CommitEntityCertificateRequest{
		OrganizationId:         req.GetTenant().GetOrganizationId(),
		TenantId:               tid,
		TenantEntityRevisionId: entityApproving.GetTenantEntityRevisionID(),
		// 如果有有组织ID就是组织提的
		IsOrganization: req.GetTenant().GetOrganizationId() != "",
	})
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return err
	}

	return nil
}

// 验证request
func validateReCreateTenantRequest(req *pb.ReCreateTenantRequest) error {
	if req.GetTenant().GetTenantId() == "" {
		return gerr.New("tenant id should not be empty")
	}
	tenant := req.GetTenant()
	if tenant.GetAddress() == nil {
		return gerr.New("address should not be empty")
	}
	if tenant.GetSocialCreditCode() == "" {
		return gerr.New("social credit code should not be empty")
	}
	if tenant.GetAddress().GetCity() == "" {
		return gerr.New("city should not be empty")
	}
	if tenant.GetAddress().GetProvince() == "" {
		return gerr.New("province should not be empty")
	}
	if tenant.GetAddress().GetStreet() == "" {
		return gerr.New("street should not be empty")
	}
	if tenant.GetName() == "" {
		return gerr.New("name should not be empty")
	}
	if tenant.GetContactPhone() == "" {
		return gerr.New("contact phone should not be empty")
	}
	if tenant.GetContactName() == "" {
		return gerr.New("contact name should not be empty")
	}
	if tenant.GetSafePhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx id should not be empty")
	}
	if req.GetPlainPassword() == "" {
		return gerr.New("password should not be empty")
	}
	return nil
}

// checkLogoUrl
// 如果是https开头的路径说明没有修改
func checkLogoUrl(url string) string {
	if strings.HasPrefix(url, "https") {
		return ""
	}
	return url
}
