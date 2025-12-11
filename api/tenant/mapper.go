package tenant

import (
	"errors"
	"net/url"
	"path"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
)

var (
	// 系统名称
	RiskSystemName = []string{"cardiovascular_system", "emotional_stress_system", "digestive_system", "skeletal_system"}
	// 平台风险预估系统
	RiskSystem = map[string][]string{
		// 心血管系统
		"cardiovascular_system": {"血压", "血糖", "血脂", "心梗"},
		// 情绪压力系统
		"emotional_stress_system": {"睡眠", "抑郁", "疲劳", "焦虑"},
		// 消化系统
		"digestive_system": {"胃动力", "免疫力", "胃炎", "咽炎"},
		// 骨骼系统
		"skeletal_system": {"颈椎", "脑供血", "脊柱", "肾功能"},
	}
)

// toAppReviewStatus
func toAppReviewStatus(s userpb.TenantReviewStatus) pb.TenantReviewStatus {
	switch s {
	case userpb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEWING:
		return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEWING
	case userpb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEW_FAIL:
		return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEW_FAIL
	case userpb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEW_SUCCESS:
		return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEW_SUCCESS
	}
	return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_UNSET
}

// toAppStatus
func toAppStatus(s userpb.TenantStatus) pb.TenantStatus {
	switch s {
	// 未认证
	case userpb.TenantStatus_TENANT_STATUS_UNAUTH:
		return pb.TenantStatus_TENANT_STATUS_UNAUTH
	case userpb.TenantStatus_TENANT_STATUS_PENDING:
		return pb.TenantStatus_TENANT_STATUS_PENDING
	case userpb.TenantStatus_TENANT_STATUS_USING:
		return pb.TenantStatus_TENANT_STATUS_USING
	}
	return pb.TenantStatus_TENANT_STATUS_UNSET
}

// toAppTenantEntity
func toAppTenantEntity(entity *userpb.Entity, s3Domain string) *pb.Entity {
	if entity == nil {
		return nil
	}
	businessUrl := ""
	if entity.GetBusinessLicenseUrl() != "" {
		link, _ := url.Parse(s3Domain)
		link.Path = path.Join(link.Path, entity.GetBusinessLicenseUrl())
		businessUrl = link.String()
	}
	logoUrl := ""
	if entity.GetLogoUrl() != "" {
		link, _ := url.Parse(s3Domain)
		link.Path = path.Join(link.Path, entity.GetLogoUrl())
		logoUrl = link.String()
	}
	return &pb.Entity{
		// 租户ID
		TenantId: entity.GetTenantId(),
		// 租户名称
		EntityName: entity.GetEntityName(),
		// 地址
		Address: toAppAddressFromUser(entity.GetAddress()),
		// 联系电话
		ContactPhone: entity.GetContactPhone(),
		// 联系人姓名
		ContactName: entity.GetContactName(),
		// 营业执照地址
		BusinessLicenseUrl: businessUrl,
		// 社会信用代码
		SocialCreditCode: entity.GetSocialCreditCode(),
		LogoUrl:          logoUrl,
		// 组织ID
		OrganizationId: entity.GetOrganizationId(),
		// 商户审核状态
		TenantReviewStatus: toAppReviewStatus(entity.GetTenantReviewStatus()),
		// 商户状态
		TenantStatus: toAppStatus(entity.GetTenantStatus()),
		// 如果审核失败,失败原因
		FailReason: entity.GetFailReason(),
		// 安全手机号
		SafePhone: entity.GetSafePhone(),
	}
}

// toAppAddressFromUser
func toAppAddressFromUser(address *userpb.Address) *pb.Address {
	if address == nil {
		return nil
	}
	return &pb.Address{
		// 省
		Province: address.GetProvince(),
		// 市
		City: address.GetCity(),
		// 区
		District: address.GetDistrict(),
		// 街道
		Street: address.GetStreet(),
	}
}

// toAppCustomer
func toAppCustomer(c *customerpb.Customer) *pb.Customer {
	return &pb.Customer{
		// 常客id
		CustomerId: c.GetCustomerId(),
		// 添加该常客员工id
		StaffId: c.GetStaffId(),
		// 常客昵称
		Nickname: c.GetNickname(),
		// 常客首字母
		Initial: c.GetInitial(),
		// 常客性别
		Gender: toAppCustomerGender(c.GetGender()),
		// 常客手机号
		Phone: c.GetPhone(),
		// 常客生日
		Birthday: toAppCustomerDate(c.GetBirthday()),
		// 常客年龄
		Age: c.GetAge(),
		// 常客身高
		Height: c.GetHeight(),
		// 常客体重
		Weight: c.GetWeight(),
		// 常客既往病史
		Pmh: c.GetPmh(),
		// 常客其他备注
		Remarks: c.GetRemarks(),
	}
}

// toAppCustomerGender
func toAppCustomerGender(c customerpb.Gender) pb.Gender {
	switch c {
	case customerpb.Gender_GENDER_FEMALE:
		return pb.Gender_GENDER_FEMALE
	case customerpb.Gender_GENDER_MALE:
		return pb.Gender_GENDER_MALE
	case customerpb.Gender_GENDER_UNSET:
		return pb.Gender_GENDER_UNSET
	default:
		return pb.Gender_GENDER_INVALID
	}
}

// toAppCustomerDate
func toAppCustomerDate(c *customerpb.Date) *pb.Date {
	return &pb.Date{
		Year:  c.GetYear(),
		Month: c.GetMonth(),
		Day:   c.GetDay(),
	}
}

// toUserImage
func toUserImage(image *pb.UploadingImage) *userpb.UploadingImage {
	if image == nil {
		return nil
	}
	return &userpb.UploadingImage{
		// image/jpg, image/png,image/webp
		Mime: image.GetMime(),
		// Image file as bytes
		Image: image.GetImage(),
		// 文件名
		Filename: image.GetFilename(),
	}
}

// toAppGender
func toAppGender(h reportpb.Gender) pb.Gender {
	switch h {
	case reportpb.Gender_GENDER_FEMALE:
		return pb.Gender_GENDER_FEMALE
	case reportpb.Gender_GENDER_MALE:
		return pb.Gender_GENDER_MALE
	default:
		return pb.Gender_GENDER_UNSET
	}
}

// toSvcCustomerPagination
func toSvcCustomerPagination(p *pb.Pagination) *customerpb.Pagination {
	return &customerpb.Pagination{
		Offset: p.GetOffset(),
		Size:   p.GetSize(),
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

// toAppStaff
func toAppStaff(staff *userpb.Staff) *pb.Staff {
	if staff == nil {
		return nil
	}
	return &pb.Staff{
		// 员工ID
		StaffId: staff.GetStaffId(),
		// 员工姓名
		Name: staff.GetName(),
		// 员工手机号
		Phone: staff.GetPhone(),
		// 员工状态
		IsActivated: staff.GetIsActivated(),
	}
}

// toSvcLanguage
func toSvcLanguage() notificationpb.Language {
	// TODO: 目前只支持中文
	return notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
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
		// 组织修改手机号
	case pb.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE, nil
		// 商户重置密码
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD, nil
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS, nil
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD, nil
	}
	return notificationpb.TemplateAction_TEMPLATE_ACTION_INVALID, errors.New("invalid template_action")
}

// toSvcCustomer
func toSvcCustomer(c *pb.Customer) *customerpb.Customer {
	return &customerpb.Customer{
		// 常客id
		CustomerId: c.GetCustomerId(),
		// 添加该常客员工id
		StaffId: c.GetStaffId(),
		// 常客昵称
		Nickname: c.GetNickname(),
		// 常客首字母
		Initial: c.GetInitial(),
		// 常客性别
		Gender: toSvcCustomerGender(c.GetGender()),
		// 常客手机号
		Phone: c.GetPhone(),
		// 常客生日
		Birthday: toSvcCustomerDate(c.GetBirthday()),
		// 常客年龄
		Age: c.GetAge(),
		// 常客身高
		Height: c.GetHeight(),
		// 常客体重
		Weight: c.GetWeight(),
		// 常客既往病史
		Pmh: c.GetPmh(),
		// 常客其他备注
		Remarks: c.GetRemarks(),
	}
}

// toSvcCustomerDate
func toSvcCustomerDate(c *pb.Date) *customerpb.Date {
	return &customerpb.Date{
		Year:  c.GetYear(),
		Month: c.GetMonth(),
		Day:   c.GetDay(),
	}
}

// toSvcCustomerGender
func toSvcCustomerGender(c pb.Gender) customerpb.Gender {
	switch c {
	case pb.Gender_GENDER_FEMALE:
		return customerpb.Gender_GENDER_FEMALE
	case pb.Gender_GENDER_MALE:
		return customerpb.Gender_GENDER_MALE
	case pb.Gender_GENDER_UNSET:
		return customerpb.Gender_GENDER_UNSET
	default:
		return customerpb.Gender_GENDER_INVALID
	}
}

// toReportDate
func toReportDate(t *pb.Date) *reportpb.Date {
	if t == nil {
		return nil
	}

	return &reportpb.Date{
		Year:  t.GetYear(),
		Month: t.GetMonth(),
		Day:   t.GetDay(),
	}
}

// toApiReport
func toApiReport(r *reportpb.SummaryReport) *pb.SummaryReport {
	if r == nil {
		return nil
	}

	return &pb.SummaryReport{
		ReportId:        r.GetReportId(),
		CreateTime:      r.GetCreateTime(),
		IsCustomer:      r.GetIsCustomer(),
		CustomerName:    r.GetCustomerName(),
		CustomerPhone:   r.GetCustomerPhone(),
		Gender:          toAppGender(r.GetGender()),
		DirtyDialectics: r.GetDirtyDialectics(),
		TenantName:      r.GetTenantName(),
		TenantId:        r.GetTenantId(),
	}
}

// toApiStaffRank
func toApiStaffRank(u *userpb.StaffRank) *pb.StaffRank {
	if u == nil {
		return nil
	}

	return &pb.StaffRank{
		Name:             u.GetName(),
		TodayMeasurement: u.GetTodayMeasurement(),
		TodayCustomer:    u.GetTodayCustomer(),
		TotalMeasurement: u.GetTotalMeasurement(),
		TotalCustomer:    u.GetTotalCustomer(),
	}
}

// toCustomerPagination
func toCustomerPagination(p *pb.Pagination) *customerpb.Pagination {
	if p == nil {
		return nil
	}

	return &customerpb.Pagination{
		Offset: p.GetOffset(),
		Size:   p.GetSize(),
	}
}

// toReportCustomerType
func toReportCustomerType(p pb.CustomerType) reportpb.CustomerType {
	switch p {
	case pb.CustomerType_CUSTOMER_TYPE_BOTH:
		return reportpb.CustomerType_CUSTOMER_TYPE_BOTH
	case pb.CustomerType_CUSTOMER_TYPE_CUSTOMER:
		return reportpb.CustomerType_CUSTOMER_TYPE_CUSTOMER
	case pb.CustomerType_CUSTOMER_TYPE_TEMP:
		return reportpb.CustomerType_CUSTOMER_TYPE_TEMP
	default:
		return reportpb.CustomerType_CUSTOMER_TYPE_UNSET
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

// toUserAddress
func toUserAddress(address *pb.Address) *userpb.Address {
	if address == nil {
		return nil
	}
	return &userpb.Address{
		// 省
		Province: address.GetProvince(),
		// 市
		City: address.GetCity(),
		// 区
		District: address.GetDistrict(),
		// 街道
		Street: address.GetStreet(),
	}
}

// toAppEntity
func toAppEntity(entity *userpb.TenantEntity, s3Domain string) *pb.Entity {
	if entity == nil {
		return nil
	}
	businessUrl := ""
	if entity.GetBusinessLicenseUrl() != "" {
		link, _ := url.Parse(s3Domain)
		link.Path = path.Join(link.Path, entity.GetBusinessLicenseUrl())
		businessUrl = link.String()
	}
	logoUrl := ""
	if entity.GetLogoUrl() != "" {
		link, _ := url.Parse(s3Domain)
		link.Path = path.Join(link.Path, entity.GetLogoUrl())
		logoUrl = link.String()
	}
	return &pb.Entity{
		// 租户ID
		TenantId: entity.GetTenantId(),
		// 租户名称
		EntityName: entity.GetName(),
		// 地址
		Address: toAppAddressFromUser(entity.GetAddress()),
		// 联系电话
		ContactPhone: entity.GetContactPhone(),
		// 联系人姓名
		ContactName: entity.GetContactName(),
		// 营业执照地址
		BusinessLicenseUrl: businessUrl,
		// 社会信用代码
		SocialCreditCode: entity.GetSocialCreditCode(),
		// 安全手机号
		SafePhone: entity.GetSafePhone(),
		LogoUrl:   logoUrl,
	}
}
