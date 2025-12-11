package organization

import (
	"errors"
	"net/url"
	"path"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	productv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reportv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
)

// toApiTenant
func toApiTenant(t *userv1.TenantEntity, s3Domain string) *pb.TenantEntity {
	if t == nil {
		return nil
	}

	return &pb.TenantEntity{
		TenantId:           t.GetTenantId(),
		OrganizationId:     t.GetOrganizationId(),
		Logo:               toApiImage(t.GetLogo()),
		LogoUrl:            toApiImageUrl(t.GetLogoUrl(), s3Domain),
		Name:               t.GetName(),
		Address:            toApiAddress(t.GetAddress()),
		BusinessLicense:    toApiImage(t.GetBusinessLicense()),
		BusinessLicenseUrl: toApiImageUrl(t.GetBusinessLicenseUrl(), s3Domain),
		SocialCreditCode:   t.GetSocialCreditCode(),
		ContactName:        t.GetContactName(),
		ContactPhone:       t.GetContactPhone(),
		TenantReviewStatus: toApiTenantReviewStatus(t.GetTenantReviewStatus()),
		FailReason:         t.GetFailReason(),
		TenantStatus:       toTenantStatus(t.GetTenantStatus()),
		CreateAt:           t.GetCreateAt(),
		SafePhone:          t.GetSafePhone(),
	}
}

// toTenantStatus
func toTenantStatus(t userv1.TenantStatus) pb.TenantStatus {
	switch t {
	case userv1.TenantStatus_TENANT_STATUS_UNSET:
		return pb.TenantStatus_TENANT_STATUS_UNSET
	case userv1.TenantStatus_TENANT_STATUS_UNAUTH:
		return pb.TenantStatus_TENANT_STATUS_UNAUTH
	case userv1.TenantStatus_TENANT_STATUS_PENDING:
		return pb.TenantStatus_TENANT_STATUS_PENDING
	case userv1.TenantStatus_TENANT_STATUS_USING:
		return pb.TenantStatus_TENANT_STATUS_USING
	}
	return pb.TenantStatus_TENANT_STATUS_INVALID
}

// toApiTenantReviewStatus
func toApiTenantReviewStatus(t userv1.TenantReviewStatus) pb.TenantReviewStatus {
	switch t {
	case userv1.TenantReviewStatus_TENANT_REVIEW_STATUS_UNSET:
		return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_UNSET
	case userv1.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEWING:
		return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEWING
	case userv1.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEW_SUCCESS:
		return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEW_SUCCESS
	case userv1.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEW_FAIL:
		return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEW_FAIL
	}
	return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_INVALID
}

// toApiImage
func toApiImage(u *userv1.UploadingImage) *pb.UploadingImage {
	if u == nil {
		return nil
	}

	return &pb.UploadingImage{
		Mime:     u.GetMime(),
		Image:    u.GetImage(),
		Filename: u.GetFilename(),
	}
}

// toApiAddress
func toApiAddress(a *userv1.Address) *pb.Address {
	if a == nil {
		return nil
	}

	return &pb.Address{
		Province: a.GetProvince(),
		City:     a.GetCity(),
		District: a.GetDistrict(),
		Street:   a.GetStreet(),
	}
}

// toSvcTenant
func toSvcTenant(p *pb.TenantEntity) *userv1.TenantEntity {
	if p == nil {
		return nil
	}

	return &userv1.TenantEntity{
		TenantId:           p.GetTenantId(),
		OrganizationId:     p.GetOrganizationId(),
		Logo:               toSvcImage(p.GetLogo()),
		LogoUrl:            p.GetLogoUrl(),
		Name:               p.GetName(),
		Address:            toSvcAddress(p.GetAddress()),
		BusinessLicense:    toSvcImage(p.GetBusinessLicense()),
		BusinessLicenseUrl: p.GetBusinessLicenseUrl(),
		SocialCreditCode:   p.GetSocialCreditCode(),
		ContactName:        p.GetContactName(),
		ContactPhone:       p.GetContactPhone(),
		SafePhone:          p.GetSafePhone(),
	}
}

// toApiImageUrl
func toApiImageUrl(imageUrl, s3Domain string) string {
	if imageUrl == "" {
		return ""
	}
	link, _ := url.Parse(s3Domain)
	link.Path = path.Join(link.Path, imageUrl)
	return link.String()
}

// toSvcImage
func toSvcImage(p *pb.UploadingImage) *userv1.UploadingImage {
	if p == nil {
		return nil
	}

	return &userv1.UploadingImage{
		Mime:     p.GetMime(),
		Image:    p.GetImage(),
		Filename: p.GetFilename(),
	}
}

// toSvcAddress
func toSvcAddress(p *pb.Address) *userv1.Address {
	if p == nil {
		return nil
	}

	return &userv1.Address{
		Province: p.GetProvince(),
		City:     p.GetCity(),
		District: p.GetDistrict(),
		Street:   p.GetStreet(),
	}
}

// toApiOrganization
func toApiOrganization(o *userv1.Organization) *pb.Organization {
	if o == nil {
		return nil
	}

	return &pb.Organization{
		OrganizationId:           o.GetOrganizationId(),
		Name:                     o.GetName(),
		Phone:                    o.GetPhone(),
		HasTenant:                o.GetHasTenant(),
		OrganizationContactPhone: o.GetOrganizationContactPhone(),
	}
}

// toAppTreatment
func toAppTreatment(t *productv1.Treatment, s3Domain string) *pb.Treatment {
	f := func(r []*productv1.TreatmentItem) []*pb.TreatmentItem {
		res := make([]*pb.TreatmentItem, len(r))
		for i, v := range r {
			res[i] = toAppTreatmentItem(v, s3Domain)
		}
		return res
	}
	return &pb.Treatment{
		// 方案id
		TreatmentId: t.GetTreatmentId(),
		// 方案名称
		TreatmentName: t.GetTreatmentName(),
		// 审核是否通过
		ReviewPass: t.GetReviewPass(),
		// 方案审核意见
		ReviewComment: t.GetReviewComment(),
		// 组织ID
		OrganizationId: t.GetOrganizationId(),
		// 是否发布
		IsPublished: t.GetIsPublished(),
		// 方案状态
		TreatmentStatus: toAppTreatmentStatus(t.GetTreatmentStatus()),
		// 风险疾病方案配置
		TreatmentItemsRiskyDisease: f(t.GetTreatmentItemsRiskyDisease()),
		// 脏腑辩证方案配置
		TreatmentItemsDirtyDialectic: f(t.GetTreatmentItemsDirtyDialectic()),
		// 理疗方案配置
		TreatmentItemsPhysicalTherapy: f(t.GetTreatmentItemsPhysicalTherapy()),
		// 体质方案配置
		TreatmentItemsPhysicalDialectics: f(t.GetTreatmentItemsPhysicalDialectics()),
		// 创建时间
		CreatedTime: t.GetCreatedTime(),
		PublishTime: t.GetPublishTime(),
	}
}

// toAppTreatmentItem
func toAppTreatmentItem(ti *productv1.TreatmentItem, s3Domain string) *pb.TreatmentItem {
	products := make([]*pb.Product, len(ti.GetProducts()))
	for i, v := range ti.GetProducts() {
		products[i] = toAppProduct(v, s3Domain)
	}
	return &pb.TreatmentItem{
		Symptom:  ti.GetSymptom(),
		Products: products,
	}
}

// toAppProduct
func toAppProduct(p *productv1.Product, s3Domain string) *pb.Product {
	if p == nil {
		return nil
	}
	return &pb.Product{
		ProductId:          p.GetProductId(),
		ProductType:        toAppProductType(p.GetProductType()),
		ProductName:        p.GetProductName(),
		ProductStatus:      toAppProductStatus(p.GetProductStatus()),
		Description:        p.GetDescription(),
		ImageUrl:           toAppImageUrl(p.GetImageUrl(), s3Domain),
		Remarks:            p.GetRemarks(),
		DrugName:           p.GetDrugName(),
		IsOtc:              p.GetIsOtc(),
		ApprovedNumber:     p.GetApprovedNumber(),
		DrugValidityPeriod: p.GetDrugValidityPeriod(),
		SymptomKeys:        p.GetSymptomKeys(),
	}
}

// toAppProductStatus
func toAppProductStatus(ps productv1.ProductStatus) pb.ProductStatus {
	switch ps {
	case productv1.ProductStatus_PRODUCT_STATUS_UNSET:
		return pb.ProductStatus_PRODUCT_STATUS_UNSET
	case productv1.ProductStatus_PRODUCT_STATUS_TO_BE_USED:
		return pb.ProductStatus_PRODUCT_STATUS_TO_BE_USED
	case productv1.ProductStatus_PRODUCT_STATUS_USING:
		return pb.ProductStatus_PRODUCT_STATUS_USING
	case productv1.ProductStatus_PRODUCT_STATUS_UNUSED:
		return pb.ProductStatus_PRODUCT_STATUS_UNUSED
	case productv1.ProductStatus_PRODUCT_STATUS_DELETED:
		return pb.ProductStatus_PRODUCT_STATUS_DELETED
	default:
		return pb.ProductStatus_PRODUCT_STATUS_INVALID
	}
}

// toSvcProductType
func toSvcProductType(pt pb.ProductType) productv1.ProductType {
	switch pt {
	case pb.ProductType_PRODUCT_TYPE_UNSET:
		return productv1.ProductType_PRODUCT_TYPE_UNSET
	case pb.ProductType_PRODUCT_TYPE_CPD:
		return productv1.ProductType_PRODUCT_TYPE_CPD
	case pb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT:
		return productv1.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT
	case pb.ProductType_PRODUCT_TYPE_NUTRITION:
		return productv1.ProductType_PRODUCT_TYPE_NUTRITION
	case pb.ProductType_PRODUCT_TYPE_SERVICE:
		return productv1.ProductType_PRODUCT_TYPE_SERVICE
	default:
		return productv1.ProductType_PRODUCT_TYPE_INVALID
	}
}

// toAppProductType
func toAppProductType(pt productv1.ProductType) pb.ProductType {
	switch pt {
	case productv1.ProductType_PRODUCT_TYPE_UNSET:
		return pb.ProductType_PRODUCT_TYPE_UNSET
	case productv1.ProductType_PRODUCT_TYPE_CPD:
		return pb.ProductType_PRODUCT_TYPE_CPD
	case productv1.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT:
		return pb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT
	case productv1.ProductType_PRODUCT_TYPE_NUTRITION:
		return pb.ProductType_PRODUCT_TYPE_NUTRITION
	case productv1.ProductType_PRODUCT_TYPE_SERVICE:
		return pb.ProductType_PRODUCT_TYPE_SERVICE
	default:
		return pb.ProductType_PRODUCT_TYPE_INVALID
	}
}

// toSvcProductStatus
func toSvcProductStatus(ps pb.ProductStatus) productv1.ProductStatus {
	switch ps {
	case pb.ProductStatus_PRODUCT_STATUS_UNSET:
		return productv1.ProductStatus_PRODUCT_STATUS_UNSET
	case pb.ProductStatus_PRODUCT_STATUS_TO_BE_USED:
		return productv1.ProductStatus_PRODUCT_STATUS_TO_BE_USED
	case pb.ProductStatus_PRODUCT_STATUS_USING:
		return productv1.ProductStatus_PRODUCT_STATUS_USING
	case pb.ProductStatus_PRODUCT_STATUS_UNUSED:
		return productv1.ProductStatus_PRODUCT_STATUS_UNUSED
	case pb.ProductStatus_PRODUCT_STATUS_DELETED:
		return productv1.ProductStatus_PRODUCT_STATUS_DELETED
	default:
		return productv1.ProductStatus_PRODUCT_STATUS_INVALID
	}
}

// toProductProduct
func toProductProduct(p *pb.Product) *productv1.Product {
	if p == nil {
		return nil
	}
	return &productv1.Product{
		ProductId:          p.GetProductId(),
		ProductType:        toSvcProductType(p.GetProductType()),
		ProductName:        p.GetProductName(),
		ProductStatus:      toSvcProductStatus(p.GetProductStatus()),
		Description:        p.GetDescription(),
		ImageUrl:           p.GetImageUrl(),
		Remarks:            p.GetRemarks(),
		DrugName:           p.GetDrugName(),
		IsOtc:              p.GetIsOtc(),
		ApprovedNumber:     p.GetApprovedNumber(),
		DrugValidityPeriod: p.GetDrugValidityPeriod(),
		SymptomKeys:        p.GetSymptomKeys(),
	}
}

// toAppImageUrl
func toAppImageUrl(imageUrl, s3Domain string) string {
	if imageUrl == "" {
		return ""
	}
	link, _ := url.Parse(s3Domain)
	link.Path = path.Join(link.Path, imageUrl)
	return link.String()
}

// toAppTreatmentStatus
func toAppTreatmentStatus(ts productv1.TreatmentStatus) pb.TreatmentStatus {
	switch ts {
	case productv1.TreatmentStatus_TREATMENT_STATUS_UNSET:
		return pb.TreatmentStatus_TREATMENT_STATUS_UNSET
	case productv1.TreatmentStatus_TREATMENT_STATUS_APPROVED:
		return pb.TreatmentStatus_TREATMENT_STATUS_APPROVED
	case productv1.TreatmentStatus_TREATMENT_STATUS_UNAPPROVED:
		return pb.TreatmentStatus_TREATMENT_STATUS_UNAPPROVED
	case productv1.TreatmentStatus_TREATMENT_STATUS_UNDER_REVIEW:
		return pb.TreatmentStatus_TREATMENT_STATUS_UNDER_REVIEW
	case productv1.TreatmentStatus_TREATMENT_STATUS_DRAFT:
		return pb.TreatmentStatus_TREATMENT_STATUS_DRAFT
	default:
		return pb.TreatmentStatus_TREATMENT_STATUS_INVALID
	}
}

func toProductUploadImage(image *pb.UploadingImage) *productv1.UploadingImage {
	if image == nil {
		return nil
	}
	return &productv1.UploadingImage{
		// image/jpg, image/png,image/webp
		Mime: image.GetMime(),
		// Image file as bytes
		Image: image.GetImage(),
		// 文件名
		Filename: image.GetFilename(),
	}
}

// toProductTreatment
func toProductTreatment(t *pb.Treatment) *productv1.Treatment {
	if t == nil {
		return nil
	}
	return &productv1.Treatment{
		TreatmentId: t.GetTreatmentId(),
		// 方案名称
		TreatmentName: t.GetTreatmentName(),
		// 风险疾病方案配置
		TreatmentItemsRiskyDisease: toProductTreatmentItem(t.GetTreatmentItemsRiskyDisease()),
		// 脏腑辩证方案配置
		TreatmentItemsDirtyDialectic: toProductTreatmentItem(t.GetTreatmentItemsDirtyDialectic()),
		// 理疗方案配置
		TreatmentItemsPhysicalTherapy: toProductTreatmentItem(t.GetTreatmentItemsPhysicalTherapy()),
		// 体质方案配置
		TreatmentItemsPhysicalDialectics: toProductTreatmentItem(t.GetTreatmentItemsPhysicalDialectics()),
		// 组织ID
		OrganizationId: t.GetOrganizationId(),
	}
}

// toProductTreatmentItem
func toProductTreatmentItem(t []*pb.TreatmentItem) []*productv1.TreatmentItem {
	if t == nil {
		return nil
	}
	pt := make([]*productv1.TreatmentItem, len(t))
	for k, v := range t {
		ps := make([]*productv1.Product, len(v.GetProducts()))
		for i, j := range v.GetProducts() {
			ps[i] = toProductProduct(j)
		}
		pt[k] = &productv1.TreatmentItem{
			// 症候
			Symptom: v.GetSymptom(),
			// 商品列表
			Products: ps,
		}
	}
	return pt
}

// toAppSymptomExposure
func toAppSymptomExposure(se *productv1.SymptomExposure) *pb.SymptomExposure {
	if se == nil {
		return nil
	}
	return &pb.SymptomExposure{
		SymptomName:   se.GetSymptomName(),
		ExposureCount: se.GetExposureCount(),
	}
}

// toAppProductExposure
func toAppProductExposure(pe *productv1.ProductExposure) *pb.ProductExposure {
	if pe == nil {
		return nil
	}
	return &pb.ProductExposure{
		ExposedDate:   toAppProductDate(pe.GetExposedDate()),
		ExposureCount: pe.GetExposureCount(),
	}
}

// toAppProductDate
func toAppProductDate(c *productv1.Date) *pb.Date {
	if c == nil {
		return nil
	}
	return &pb.Date{
		Year:  c.GetYear(),
		Month: c.GetMonth(),
		Day:   c.GetDay(),
	}
}

// toReportPagination
func toReportPagination(p *pb.Pagination) *reportv1.Pagination {
	if p == nil {
		return nil
	}
	return &reportv1.Pagination{
		// 分页偏移量
		Offset: p.GetOffset(),
		// 期望分页查询记录的数量
		Size: p.GetSize(),
	}
}

// toAppTenantReportCount
func toAppTenantReportCount(p *reportv1.TenantReportCount) *pb.TenantReportCount {
	if p == nil {
		return nil
	}
	return &pb.TenantReportCount{
		// 商户名称
		TenantName: p.GetTenantName(),
		// 当月测量次数
		MonthlyCustomerMeasurementCount: p.GetMonthlyCustomerMeasurementCount(),
		// 同比
		YearOnYear: p.GetYearOnYear(),
		// 环比
		MonthOnMonth: p.GetMonthOnMonth(),
		// 时间
		Time: p.GetTime(),
	}
}

// toProductPagination
func toProductPagination(p *pb.Pagination) *productv1.Pagination {
	if p == nil {
		return nil
	}
	return &productv1.Pagination{
		// 分页偏移量
		Offset: p.GetOffset(),
		// 期望分页查询记录的数量
		Size: p.GetSize(),
	}
}

// toApiPagination
func toApiPagination(p *pb.Pagination) *userv1.Pagination {
	if p == nil {
		return nil
	}
	return &userv1.Pagination{
		// 分页偏移量
		Offset: p.GetOffset(),
		// 期望分页查询记录的数量
		Size: p.GetSize(),
	}
}

// toApiOrganizationTenantCompare
func toApiOrganizationTenantCompare(u *userv1.OrganizationTenantCompare) *pb.OrganizationTenantCompare {
	if u == nil {
		return nil
	}

	return &pb.OrganizationTenantCompare{
		Date:         u.GetDate(),
		TenantName:   u.GetTenantName(),
		MonthCount:   u.GetMonthCount(),
		YearOnYear:   u.GetYearOnYear(),
		MonthOnMonth: u.GetMonthOnMonth(),
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

// toSvcPagination
func toSvcPagination(p *pb.Pagination) *userv1.Pagination {
	if p == nil {
		return nil
	}

	return &userv1.Pagination{
		Offset: p.GetOffset(),
		Size:   p.GetSize(),
	}
}

// toSvcTemplateAction
func toSvcTemplateAction(action pb.TemplateAction) (notificationpb.TemplateAction, error) {
	switch action {
	// 组织注册
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION, nil
		// 组织登陆
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_ORGANIZATION:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_ORGANIZATION, nil
		// 组织重置密码
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD, nil
		// 组织绑定手机号
	case pb.TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE, nil
		// 组织账单
	case pb.TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING, nil
		// 商户创建成功
	case pb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS, nil
		// 商户创建失败
	case pb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL, nil
		// 商户绑定手机号
	case pb.TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE, nil
		// 商户登陆
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_TENANT:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_TENANT, nil
		// 商户修改手机号
	case pb.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE, nil
		// 商户重置密码
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD, nil
		// 组织修改手机号
	case pb.TemplateAction_TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE, nil
	// 金姆平台登录
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS, nil
	// 重置员工登录密码
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD, nil
	}
	return notificationpb.TemplateAction_TEMPLATE_ACTION_INVALID, errors.New("invalid template_action")
}

// toReportCustomerType
func toReportCustomerType(c pb.CustomerType) reportv1.CustomerType {
	switch c {
	case pb.CustomerType_CUSTOMER_TYPE_BOTH:
		return reportv1.CustomerType_CUSTOMER_TYPE_BOTH
	case pb.CustomerType_CUSTOMER_TYPE_CUSTOMER:
		return reportv1.CustomerType_CUSTOMER_TYPE_CUSTOMER
	case pb.CustomerType_CUSTOMER_TYPE_TEMP:
		return reportv1.CustomerType_CUSTOMER_TYPE_TEMP
	}
	return reportv1.CustomerType_CUSTOMER_TYPE_BOTH
}

// toAppGender
func toAppGender(g reportv1.Gender) pb.Gender {
	switch g {
	case reportv1.Gender_GENDER_FEMALE:
		return pb.Gender_GENDER_FEMALE
	case reportv1.Gender_GENDER_MALE:
		return pb.Gender_GENDER_MALE
	}
	return pb.Gender_GENDER_FEMALE
}

// toAppSummaryReport
func toAppSummaryReport(c *reportv1.SummaryReport) *pb.SummaryReport {
	if c == nil {
		return nil
	}
	return &pb.SummaryReport{
		// 报告ID
		ReportId: c.GetReportId(),
		// 创建时间
		CreateTime: c.GetCreateTime(),
		// 是否是常客报告
		IsCustomer: c.GetIsCustomer(),
		// 客户名称
		CustomerName: c.GetCustomerName(),
		// 客户手机号
		CustomerPhone: c.GetCustomerPhone(),
		// 性别
		Gender: toAppGender(c.GetGender()),
		// 脏腑辨证
		DirtyDialectics: c.GetDirtyDialectics(),
		// 商户名称
		TenantName: c.GetTenantName(),
		TenantId:   c.GetTenantId(),
	}
}

// toAppNotification
func toAppNotification(n *reviewpb.ReviewNotification) *pb.ReviewNotification {
	if n == nil {
		return nil
	}
	return &pb.ReviewNotification{
		// 通知ID
		NotificationId: n.GetNotificationId(),
		// 通知标题
		NotificationTitle: n.GetNotificationTitle(),
		// 是否通过
		Pass: n.GetPass(),
		// 审核类型
		ReviewType: toAppReviewType(n.GetReviewType()),
		// 审核意见
		ReviewComment: n.GetReviewComment(),
		// 创建时间
		CreatedTime: n.GetCreatedTime(),
	}
}

// toAppReviewType
func toAppReviewType(r reviewpb.ReviewType) pb.ReviewType {
	switch r {
	case reviewpb.ReviewType_REVIEW_TYPE_CERTIFICATE:
		return pb.ReviewType_REVIEW_TYPE_CERTIFICATE
	case reviewpb.ReviewType_REVIEW_TYPE_PRODUCT:
		return pb.ReviewType_REVIEW_TYPE_PRODUCT
	default:
		return pb.ReviewType_REVIEW_TYPE_UNSET
	}
}

// toCustomerPagination
func toCustomerPagination(p *pb.Pagination) *customerv1.Pagination {
	if p == nil {
		return nil
	}

	return &customerv1.Pagination{
		Offset: p.GetOffset(),
		Size:   p.GetSize(),
	}
}

// toReportDate
func toReportDate(d *pb.Date) *reportv1.Date {
	if d == nil {
		return nil
	}

	return &reportv1.Date{
		Year:  d.GetYear(),
		Month: d.GetMonth(),
		Day:   d.GetDay(),
	}
}

func toUserTenantEntity(tenant *pb.TenantEntity) *userv1.TenantEntity {
	if tenant == nil {
		return nil
	}
	return &userv1.TenantEntity{
		// 商户ID
		TenantId: tenant.GetTenantId(),
		// 商户名称
		Name: tenant.GetName(),
		// 地址
		Address: toUserAddress(tenant.GetAddress()),
		// 安全电话
		SafePhone: tenant.GetSafePhone(),
		// 联系人姓名
		ContactName: tenant.GetContactName(),
		// 营业执照
		BusinessLicense: toUserUploadImage(tenant.GetBusinessLicense()),
		// 社会信用代码
		SocialCreditCode: tenant.GetSocialCreditCode(),
		// 联系人手机号
		ContactPhone: tenant.GetContactPhone(),
		// 组织ID
		OrganizationId: tenant.GetOrganizationId(),
		// logo_url
		LogoUrl:            tenant.GetLogoUrl(),
		BusinessLicenseUrl: tenant.GetBusinessLicenseUrl(),
	}
}

func toUserAddress(ad *pb.Address) *userv1.Address {
	if ad == nil {
		return nil
	}
	return &userv1.Address{
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

func toUserUploadImage(image *pb.UploadingImage) *userv1.UploadingImage {
	if image == nil {
		return nil
	}
	return &userv1.UploadingImage{
		// image/jpg, image/png,image/webp
		Mime: image.GetMime(),
		// Image file as bytes
		Image: image.GetImage(),
		// 文件名
		Filename: image.GetFilename(),
	}
}
