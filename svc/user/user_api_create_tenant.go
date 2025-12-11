package user

import (
	"bytes"
	"context"
	gerr "errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/image"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

const (
	// ErrTenantExceedLimit
	ErrTenantExceedLimit = 5309
)

const (
	StaffCountQuota = 6
)

// 创建商户
func (u *UserAPIHandler) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest, rsp *pb.CreateTenantResponse) error {
	err := validateCreateTenantRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationpb.VerifyPhoneVerificationCodeRequest{
		// 手机号
		Phone: req.GetTenant().GetSafePhone(),
		// 短信验证码
		SmsCode: req.GetSmsCode(),
		// 凭证 ID
		TxId: req.GetTxId(),
		// 短信验证码类型
		Action: notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE,
	})
	if err != nil {
		return err
	}

	// 查询组织
	organization, err := u.userStore.GetOrganizationByOrganizationId(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if organization == nil {
		return errors.Errorf(ErrOrganizationNotExist, "organization not found [%s]", req.GetOrganizationId())
	}

	// 查询组织下商户数量
	ots, err := u.userStore.ListOrganizationExistTenants(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if len(ots) >= int(organization.GetTenantLimit()) {
		return errors.Errorf(ErrTenantExceedLimit, "organization[%s] tenant exceed limit", req.GetOrganizationId())
	}

	// 通过手机号查询商户
	userExist, err := u.userStore.GetTenantUserByPhone(ctx, req.GetTenant().GetContactPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if userExist != nil {
		// 手机号存在
		return errors.Errorf(ErrPhoneHasBeenUsed, "tenant phone[%s] has been used", req.GetTenant().GetContactPhone())
	}

	// 获取组织下所有商户
	tids, err := u.userStore.ListAllOrganizationTenantIds(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	tenantIdMap := make(map[string]bool)
	for _, v := range tids {
		tenantIdMap[v] = true
	}

	// 生成商户ID
	exist := true
	times := 0
	tid := generateTenantID(req.GetOrganizationId())
	for exist {
		// 如果循环500次还有冲突则报错
		if times == 500 {
			return errors.Error(codes.InvalidOperation, "generate tenant id failed")
		}
		// 查询tid是否存在
		if _, ok := tenantIdMap[tid]; !ok {
			// 存在
			exist = false
		}
		times++
	}

	// 获取商户资质路径和 logo
	businessUrl := ""
	logoUrl := ""
	storeName := req.GetTenant().GetName()
	province := req.GetTenant().GetAddress().GetProvince()
	city := req.GetTenant().GetAddress().GetCity()
	district := req.GetTenant().GetAddress().GetDistrict()
	street := req.GetTenant().GetAddress().GetStreet()
	socialCreditCode := req.GetTenant().GetSocialCreditCode()
	// 提交商户资质
	if req.GetTenant().GetBusinessLicense() != nil {
		businessUrl, err = u.uploadImage(req.GetTenant().GetBusinessLicense())
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
	} else {
		businessUrl = req.GetTenant().GetBusinessLicenseUrl()
	}

	// 提交 logo
	logoUrl = req.GetTenant().GetLogoUrl()

	// 1.生成商家组织关系
	// 2.生成商户(user)
	// 3.提审
	ctx = u.userStore.BeginTx(ctx)
	orgTenant := &domain.OrganizationTenant{
		OrganizationTenantID: xid.New().String(),
		OrganizationID:       req.GetOrganizationId(),
		TenantID:             tid,
		// 创建就是审核状态
		ReviewStatus: domain.TenantReviewStatusReviewing,
		// 刚创建的商户是未认证状态
		IsActivated: false,
	}
	err = u.userStore.CreateOrganizationTenant(ctx, orgTenant)
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 生成 hash 密码
	hashedPassword := generateHashSHA256(req.GetPlainPassword())
	// 初始化商户 user
	user := &domain.User{
		UserID:         xid.New().String(),
		OrganizationID: req.GetOrganizationId(),
		TenantID:       tid,
		Nickname:       req.GetTenant().GetContactName(),
		Phone:          req.GetTenant().GetSafePhone(),
		RoleType:       domain.RoleTypeTenant,
		HashedPassword: hashedPassword,
		IsActivated:    true,
	}
	err = u.userStore.CreateUser(ctx, user)
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 初始化 tenant_entity
	te := &domain.TenantEntity{
		TenantID:            tid,
		StaffCountQuota:     StaffCountQuota,
		ReportSharingStatus: true,
		ConstitutionStatus:  true,
		Overdue:             7,
	}
	err = u.userStore.CreateTenantEntity(ctx, te)
	if err != nil {
		u.userStore.RollbackTx(ctx)
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 创建一个新的 revision
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
		OrganizationId:         req.GetOrganizationId(),
		TenantId:               tid,
		TenantEntityRevisionId: entityApproving.GetTenantEntityRevisionID(),
		IsOrganization:         true,
	})
	if err != nil {
		return err
	}

	rsp.Tenant = &pb.TenantEntity{
		TenantId:           tid,
		OrganizationId:     req.GetOrganizationId(),
		Address:            req.GetTenant().GetAddress(),
		Name:               req.GetTenant().GetName(),
		LogoUrl:            req.GetTenant().GetLogoUrl(),
		BusinessLicenseUrl: req.GetTenant().GetBusinessLicenseUrl(),
		ContactName:        req.GetTenant().GetContactName(),
		ContactPhone:       req.GetTenant().GetContactPhone(),
		SocialCreditCode:   req.GetTenant().GetSocialCreditCode(),
		TenantReviewStatus: pb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEWING,
		TenantStatus:       pb.TenantStatus_TENANT_STATUS_UNAUTH,
		SafePhone:          req.GetTenant().GetSafePhone(),
	}

	return nil
}

// 生成商户ID
func generateTenantID(orgId string) string {
	letter := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := len(letter)

	// 设置随机种子
	rand.NewSource(time.Now().UnixNano())

	// 生成4个字符的字符串
	result := ""
	for i := 0; i < 4; i++ {
		index := rand.Intn(length)
		result += string(letter[index])
	}
	return fmt.Sprintf("%s%s", orgId, result)
}

// 验证 request
func validateCreateTenantRequest(req *pb.CreateTenantRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	// 验证相关信息是否存在
	tenant := req.GetTenant()
	if tenant.GetAddress() == nil {
		if tenant.GetAddress() == nil {
			return gerr.New("address should not be empty")
		}
	}
	return nil
}

// 上传图片
func (u *UserAPIHandler) uploadImage(logo *pb.UploadingImage) (string, error) {
	const (
		ImageSizeLimit = 2 * 1024 * 1024
	)
	logoUrl := ""
	if logo == nil {
		return "", nil
	}
	// 检查 mime 和大小
	if len(logo.GetImage()) > ImageSizeLimit {
		return "", gerr.New("image should not exceed 2M")
	}
	// 检测 mime
	suffix, err := image.GetImageSuffix(logo.GetMime())
	if err != nil {
		return "", err
	}
	// 如果传了 logo 就去上传 logo
	// 生成图片的 ID
	logoUrl = fmt.Sprintf("%s.%s", xid.New().String(), suffix)
	path, err := u.s3Store.Save(logoUrl, logo.GetMime(), byteReader(logo.GetImage()))
	if err != nil {
		return "", err
	}
	return path, nil
}

// byteReader 将字节格式化为可读文本流
func byteReader(img []byte) io.ReadSeeker {
	return bytes.NewReader(img)
}
