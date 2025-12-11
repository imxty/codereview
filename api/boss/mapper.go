package boss

import (
	"net/url"
	"path"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	notificationv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	reviewv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
)

func toApiReviewType(t reviewv1.ReviewType) pb.ReviewType {
	switch t {
	case reviewv1.ReviewType_REVIEW_TYPE_CERTIFICATE:
		return pb.ReviewType_REVIEW_TYPE_CERTIFICATE
	case reviewv1.ReviewType_REVIEW_TYPE_PRODUCT:
		return pb.ReviewType_REVIEW_TYPE_PRODUCT
	default:
		return pb.ReviewType_REVIEW_TYPE_INVALID
	}
}

// toAppTenantEntity
func toAppTenantEntity(s3Domain string, entity *reviewv1.TenantEntity) *pb.TenantEntity {
	if entity == nil {
		return nil
	}
	logoUrl, businessUrl := "", ""
	if entity.GetLogoUrl() != "" {
		link, _ := url.Parse(s3Domain)
		link.Path = path.Join(link.Path, entity.GetLogoUrl())
		logoUrl = link.String()
	}
	if entity.GetBusinessLicenseUrl() != "" {
		link, _ := url.Parse(s3Domain)
		link.Path = path.Join(link.Path, entity.GetBusinessLicenseUrl())
		businessUrl = link.String()
	}
	return &pb.TenantEntity{
		// 租户ID
		TenantId: entity.GetTenantId(),
		// 选填，商铺logo
		LogoUrl: logoUrl,
		// 商铺名称
		EntityName: entity.GetEntityName(),
		// 地址
		Address: toAppAddress(entity.GetAddress()),
		// 联系人姓名
		ContactName: entity.GetContactName(),
		// 联系人电话
		ContactPhone: entity.GetContactPhone(),
		// 营业执照
		BusinessLicenseUrl: businessUrl,
		// 社会信用代码
		SocialCreditCode: entity.GetSocialCreditCode(),
	}
}

// toAppAddress
func toAppAddress(ad *reviewv1.Address) *pb.Address {
	if ad == nil {
		return nil
	}
	return &pb.Address{
		// 省
		Province: ad.GetProvince(),
		// 市
		City: ad.GetCity(),
		// 区
		District: ad.GetDistrict(),
		// 街道
		Street: ad.GetStreet(),
	}
}

// toAppProductMap
func toAppProductMap(s3Domain string, m *reviewv1.SymptomProductMap) *pb.SymptomProductMap {
	if m == nil {
		return nil
	}
	ps := make([]*pb.Product, len(m.GetProducts()))
	for k, v := range m.GetProducts() {
		ps[k] = toAppProduct(s3Domain, v)
	}
	return &pb.SymptomProductMap{
		// 症候
		Symptom: m.GetSymptom(),
		// 商品列表
		Products: ps,
	}
}

// toAppProduct
func toAppProduct(s3Domain string, t *reviewv1.Product) *pb.Product {
	if t == nil {
		return nil
	}
	imgUrl := ""
	if t.GetImageUrl() != "" {
		link, _ := url.Parse(s3Domain)
		link.Path = path.Join(link.Path, t.GetImageUrl())
		imgUrl = link.String()
	}
	return &pb.Product{
		// 商品id
		ProductId: t.GetProductId(),
		// 商品类型
		ProductType: tAppProductType(t.GetProductType()),
		// 商品名称
		ProductName: t.GetProductName(),
		// 商品状态
		ProductStatus: toAppProductStatus(t.GetProductStatus()),
		// 商品介绍
		Description: t.GetDescription(),
		// 商品图片
		ImageUrl: imgUrl,
		// 商品备注
		Remarks: t.GetRemarks(),
		// 药品名称(商品类型为理疗服务方案不用传)
		DrugName: t.GetDrugName(),
		// 药品是非处方药 (商品类型为理疗服务方案不用传)
		IsOtc: t.GetIsOtc(),
		// 准字号 (商品类型为理疗服务方案不用传)
		ApprovedNumber: t.GetApprovedNumber(),
		// 药品准效期 (商品类型为理疗服务方案不用传)
		DrugValidityPeriod: t.GetDrugValidityPeriod(),
	}
}

// tAppProductType
func tAppProductType(productType reviewv1.ProductType) pb.ProductType {
	switch productType {
	case reviewv1.ProductType_PRODUCT_TYPE_UNSET:
		return pb.ProductType_PRODUCT_TYPE_UNSET
	case reviewv1.ProductType_PRODUCT_TYPE_SERVICE:
		return pb.ProductType_PRODUCT_TYPE_SERVICE
	case reviewv1.ProductType_PRODUCT_TYPE_CPD:
		return pb.ProductType_PRODUCT_TYPE_CPD
	case reviewv1.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT:
		return pb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT
	case reviewv1.ProductType_PRODUCT_TYPE_NUTRITION:
		return pb.ProductType_PRODUCT_TYPE_NUTRITION
	}
	return pb.ProductType_PRODUCT_TYPE_INVALID
}

// toAppProductStatus
func toAppProductStatus(protoStatus reviewv1.ProductStatus) pb.ProductStatus {
	switch protoStatus {
	case reviewv1.ProductStatus_PRODUCT_STATUS_UNSET:
		return pb.ProductStatus_PRODUCT_STATUS_UNSET
	case reviewv1.ProductStatus_PRODUCT_STATUS_USING:
		return pb.ProductStatus_PRODUCT_STATUS_USING
	case reviewv1.ProductStatus_PRODUCT_STATUS_TO_BE_USED:
		return pb.ProductStatus_PRODUCT_STATUS_TO_BE_USED
	case reviewv1.ProductStatus_PRODUCT_STATUS_UNUSED:
		return pb.ProductStatus_PRODUCT_STATUS_UNUSED
	case reviewv1.ProductStatus_PRODUCT_STATUS_DELETED:
		return pb.ProductStatus_PRODUCT_STATUS_DELETED
	}
	return pb.ProductStatus_PRODUCT_STATUS_INVALID
}

// toSvcReviewStatus
func toSvcReviewStatus(status pb.ReviewStatus) reviewv1.ReviewStatus {
	switch status {
	// 待审核
	case pb.ReviewStatus_REVIEW_STATUS_PENDING_REVIEW:
		return reviewv1.ReviewStatus_REVIEW_STATUS_PENDING_REVIEW
		// 审核中
	case pb.ReviewStatus_REVIEW_STATUS_REVIEWING:
		return reviewv1.ReviewStatus_REVIEW_STATUS_REVIEWING
	// 审核成功
	case pb.ReviewStatus_REVIEW_STATUS_SUCCESS:
		return reviewv1.ReviewStatus_REVIEW_STATUS_SUCCESS
		// 审核失败
	case pb.ReviewStatus_REVIEW_STATUS_FAIL:
		return reviewv1.ReviewStatus_REVIEW_STATUS_FAIL
	default:
		return reviewv1.ReviewStatus_REVIEW_STATUS_UNSET
	}
}

// toApiReviewStatus
func toApiReviewStatus(status reviewv1.ReviewStatus) pb.ReviewStatus {
	switch status {
	// 待审核
	case reviewv1.ReviewStatus_REVIEW_STATUS_PENDING_REVIEW:
		return pb.ReviewStatus_REVIEW_STATUS_PENDING_REVIEW
		// 审核中
	case reviewv1.ReviewStatus_REVIEW_STATUS_REVIEWING:
		return pb.ReviewStatus_REVIEW_STATUS_REVIEWING
	// 审核成功
	case reviewv1.ReviewStatus_REVIEW_STATUS_SUCCESS:
		return pb.ReviewStatus_REVIEW_STATUS_SUCCESS
		// 审核失败
	case reviewv1.ReviewStatus_REVIEW_STATUS_FAIL:
		return pb.ReviewStatus_REVIEW_STATUS_FAIL
	default:
		return pb.ReviewStatus_REVIEW_STATUS_UNSET
	}
}

// toReviewPagination
func toReviewPagination(p *pb.Pagination) *reviewv1.Pagination {
	if p == nil {
		return nil
	}
	return &reviewv1.Pagination{
		Offset: p.GetOffset(),
		Size:   p.GetSize(),
	}
}

// toUserPagination
func toUserPagination(p *pb.Pagination) *userv1.Pagination {
	if p == nil {
		return nil
	}
	return &userv1.Pagination{
		Offset: p.GetOffset(),
		Size:   p.GetSize(),
	}
}

// toAppReviewIssue
func toAppReviewIssue(issue *reviewv1.ReviewIssue) *pb.ReviewIssue {
	if issue == nil {
		return nil
	}

	return &pb.ReviewIssue{
		// 工单ID
		ReviewIssueId: issue.GetReviewIssueId(),
		// 提审者ID
		SubmitterOrganizationId: issue.GetSubmitterOrganizationId(),
		SubmitterTenantId:       issue.GetSubmitterTenantId(),
		// 提审时间
		SubmitTime: issue.GetSubmitTime(),
		// 审核人名称
		ReviewNickName: issue.GetReviewNickName(),
		// 工单状态
		Status: toApiReviewStatus(issue.GetStatus()),
		// 审核类型
		ReviewType: toApiReviewType(issue.GetReviewType()),
		// 是否被用户取消
		IsCancled: issue.GetIsCancled(),
		// 是否通过
		IsPassed: issue.GetIsPassed(),
		// 未通过原因
		FailReason: issue.GetFailReason(),
		// 是否是工单的接单者
		IsOwner: issue.GetIsOwner(),
		// 审核人ID
		ReviewerUserId:   issue.GetReviewerUserId(),
		OrganizationName: issue.GetOrganizationName(),
	}
}

// toAppSystemUser
func toAppSystemUser(u *userv1.SystemUser) *pb.SystemUser {
	if u == nil {
		return nil
	}
	return &pb.SystemUser{
		// 用户ID
		UserId: u.GetUserId(),
		// 手机号
		Phone: u.GetPhone(),
		// 姓名
		Nickname: u.GetNickname(),
		RoleType: u.GetRoleType(),
		Remark:   u.GetRemark(),
	}
}

// toAppReportCount
func toAppReportCount(u *reportpb.ReportCount) *pb.ReportCount {
	if u == nil {
		return nil
	}
	return &pb.ReportCount{
		// 组织名称
		OrganizationName: u.GetOrganizationName(),
		// 商户名称
		TenantName: u.GetTenantName(),
		// 体验测量次数
		TempCustomerReportCount: u.GetTempCustomerReportCount(),
		// vip测量次数
		CustomerReportCount: u.GetCustomerReportCount(),
		// 开始时间
		StartTime: u.GetStartTime(),
	}
}

// toAppReportStaffCount
func toAppReportStaffCount(u *reportpb.StaffReport) *pb.StaffReport {
	if u == nil {
		return nil
	}
	return &pb.StaffReport{
		// 员工名称
		Nickname: u.GetNickname(),
		// phone
		Phone: u.GetPhone(),
		// bool
		IsActivated: u.GetIsActivated(),
		// 体验测量次数
		TempCustomerReportCount: u.GetTempCustomerReportCount(),
		// vip测量次数
		CustomerReportCount: u.GetCustomerReportCount(),
		// 开始时间
		StartTime: u.GetStartTime(),
	}
}

// toReportPagination
func toReportPagination(p *pb.Pagination) *reportpb.Pagination {
	if p == nil {
		return nil
	}
	return &reportpb.Pagination{
		Offset: p.GetOffset(),
		Size:   p.GetSize(),
	}
}

// toSvcPagePrivileges
func toSvcPagePrivileges(p pb.PagePrivilege) userv1.PagePrivilege {
	switch p {
	case pb.PagePrivilege_PAGE_PRIVILEGE_UNSET:
		return userv1.PagePrivilege_PAGE_PRIVILEGE_UNSET
	case pb.PagePrivilege_PAGE_PRIVILEGE_OPERATIONAL:
		return userv1.PagePrivilege_PAGE_PRIVILEGE_OPERATIONAL
	case pb.PagePrivilege_PAGE_PRIVILEGE_REVIEW:
		return userv1.PagePrivilege_PAGE_PRIVILEGE_REVIEW
	case pb.PagePrivilege_PAGE_PRIVILEGE_TENANT:
		return userv1.PagePrivilege_PAGE_PRIVILEGE_TENANT
	case pb.PagePrivilege_PAGE_PRIVILEGE_SUBSCRIPTION:
		return userv1.PagePrivilege_PAGE_PRIVILEGE_SUBSCRIPTION
	case pb.PagePrivilege_PAGE_PRIVILEGE_LIST:
		return userv1.PagePrivilege_PAGE_PRIVILEGE_LIST
	}
	return userv1.PagePrivilege_PAGE_PRIVILEGE_INVALID
}

// toApiPagePrivilege
func toApiPagePrivilege(p userv1.PagePrivilege) pb.PagePrivilege {
	switch p {
	case userv1.PagePrivilege_PAGE_PRIVILEGE_UNSET:
		return pb.PagePrivilege_PAGE_PRIVILEGE_UNSET
	case userv1.PagePrivilege_PAGE_PRIVILEGE_OPERATIONAL:
		return pb.PagePrivilege_PAGE_PRIVILEGE_OPERATIONAL
	case userv1.PagePrivilege_PAGE_PRIVILEGE_REVIEW:
		return pb.PagePrivilege_PAGE_PRIVILEGE_REVIEW
	case userv1.PagePrivilege_PAGE_PRIVILEGE_TENANT:
		return pb.PagePrivilege_PAGE_PRIVILEGE_TENANT
	case userv1.PagePrivilege_PAGE_PRIVILEGE_SUBSCRIPTION:
		return pb.PagePrivilege_PAGE_PRIVILEGE_SUBSCRIPTION
	case userv1.PagePrivilege_PAGE_PRIVILEGE_LIST:
		return pb.PagePrivilege_PAGE_PRIVILEGE_LIST
	}
	return pb.PagePrivilege_PAGE_PRIVILEGE_INVALID
}

// toApiPrivilegeGroup
func toApiPrivilegeGroup(p *userv1.PrivilegeGroup) *pb.PrivilegeGroup {
	if p == nil {
		return nil
	}

	// 转换页面权限
	pagePrivileges := make([]pb.PagePrivilege, len(p.GetPagePrivileges()))
	for k, v := range p.GetPagePrivileges() {
		pagePrivileges[k] = toApiPagePrivilege(v)
	}

	return &pb.PrivilegeGroup{
		PrivilegeId:    p.GetPrivilegeId(),
		PrivilegeName:  p.GetPrivilegeName(),
		Remark:         p.GetRemark(),
		PagePrivileges: pagePrivileges,
	}
}

// toApiSystemUser
func toApiSystemUser(s *userv1.SystemUser) *pb.SystemUser {
	if s == nil {
		return nil
	}

	return &pb.SystemUser{
		UserId:         s.GetUserId(),
		Phone:          s.GetPhone(),
		Nickname:       s.GetNickname(),
		RoleType:       s.GetRoleType(),
		PrivilegeGroup: toApiPrivilegeGroup(s.GetPrivilegeGroup()),
		IsDeleted:      s.GetIsDeleted(),
		Remark:         s.GetRemark(),
	}
}

// toApiSubscriptionPeriod
func toApiSubscriptionPeriod(u *userv1.SubscriptionPeriod) *pb.SubscriptionPeriod {
	if u == nil {
		return nil
	}

	return &pb.SubscriptionPeriod{
		TenantId:         u.GetTenantId(),
		CreatedAt:        u.GetCreatedAt(),
		OrganizationName: u.GetOrganizationName(),
		TenantName:       u.GetTenantName(),
		ContactName:      u.GetContactName(),
		ContactPhone:     u.GetContactPhone(),
		Years:            u.GetYears(),
		ExpiredAt:        u.GetExpiredAt(),
	}
}

// toSvcStatus
func toSvcStatus(p pb.TenantStatus) userv1.TenantStatus {
	switch p {
	case pb.TenantStatus_TENANT_STATUS_UNSET:
		return userv1.TenantStatus_TENANT_STATUS_UNSET
	case pb.TenantStatus_TENANT_STATUS_UNAUTH:
		return userv1.TenantStatus_TENANT_STATUS_UNAUTH
	case pb.TenantStatus_TENANT_STATUS_PENDING:
		return userv1.TenantStatus_TENANT_STATUS_PENDING
	case pb.TenantStatus_TENANT_STATUS_USING:
		return userv1.TenantStatus_TENANT_STATUS_USING
	}
	return userv1.TenantStatus_TENANT_STATUS_INVALID
}

// toSvcLanguage
func toSvcLanguage(l pb.Language) notificationv1.Language {
	switch l {
	case pb.Language_LANGUAGE_UNSET:
		return notificationv1.Language_LANGUAGE_UNSET
	case pb.Language_LANGUAGE_SIMPLIFIED_CHINESE:
		return notificationv1.Language_LANGUAGE_SIMPLIFIED_CHINESE
	case pb.Language_LANGUAGE_TRADITIONAL_CHINESE:
		return notificationv1.Language_LANGUAGE_TRADITIONAL_CHINESE
	case pb.Language_LANGUAGE_ENGLISH:
		return notificationv1.Language_LANGUAGE_ENGLISH
	}
	return notificationv1.Language_LANGUAGE_INVALID
}

// toSvcTemplateAction
func toSvcTemplateAction(t pb.TemplateAction) notificationv1.TemplateAction {
	switch t {
	case pb.TemplateAction_TEMPLATE_ACTION_UNSET:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_UNSET
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_ORGANIZATION:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_SIGNIN_ORGANIZATION
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD
	case pb.TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE
	case pb.TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING
	case pb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS
	case pb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL
	case pb.TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_TENANT:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_SIGNIN_TENANT
	case pb.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_APP:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_SIGNIN_APP
	case pb.TemplateAction_TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD:
		return notificationv1.TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD
	}
	return notificationv1.TemplateAction_TEMPLATE_ACTION_INVALID
}
