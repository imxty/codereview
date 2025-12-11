package app

import (
	gerr "errors"
	"net/url"
	"path"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	devicepb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
)

// mapper app层数据结构和service映射的工具
// rule: 1.从app到service如果数据结构的name相同，func的名字为toSvcXXXX
//       2.从app到service如果数据结构的name不相同,app数据结构的名称为A，service名称为B，func的名字为toSvcBFromA
//       3.反之，从service到app一样，前缀为toAppXXXX

// toAppDevice
func toAppDevice(d *devicepb.Device) *pb.Device {
	return &pb.Device{ // 设备ID
		DeviceId: d.GetDeviceId(),
		// 设备mac地址
		Mac: d.GetMac(),
		// 设备型号
		Model: d.GetModel(),
	}
}

// toAppStaff
func toAppStaff(s *userpb.Staff, s3Domain string) *pb.Staff {
	if s == nil {
		return nil
	}
	return &pb.Staff{
		// 员工ID
		StaffId: s.GetStaffId(),
		// 手机号
		Phone: s.GetPhone(),
		// 员工姓名
		Name: s.GetName(),
		// 租户信息
		Tenant: toAppTenant(s.GetTenant(), s3Domain),
	}
}

// toAppTenant
func toAppTenant(t *userpb.Tenant, s3Domain string) *pb.Tenant {
	if t == nil {
		return nil
	}
	logoUrl := ""
	if t.GetLogoUrl() != "" {
		link, _ := url.Parse(s3Domain)
		link.Path = path.Join(link.Path, t.GetLogoUrl())
		logoUrl = link.String()
	}
	return &pb.Tenant{
		// 租户ID
		TenantId: t.GetTenantId(),
		// 租户名称
		TenantName: t.GetTenantName(),
		// logo地址
		LogoUrl: logoUrl,
		// 地址
		Address: toAppAddress(t.GetAddress()),
		// 联系电话
		ContactPhone: t.GetContactPhone(),
	}
}

// toAppAddress
func toAppAddress(ad *userpb.Address) *pb.Address {
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

// toSvcCustomer
func toSvcCustomer(c *pb.Customer) (*customerpb.Customer, error) {
	gender, _ := toSvcCustomerGender(c.GetGender())
	birthday := toSvcBirthday(c.GetBirthday())
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
		Gender: gender,
		// 常客手机号
		Phone: c.GetPhone(),
		// 常客生日
		Birthday: birthday,
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
	}, nil
}

// toReportSvcCustomer
func toReportSvcCustomer(c *pb.Customer) *reportpb.Customer {
	gender := toSvcGender(c.GetGender())
	birthday := toReportSvcBirthday(c.GetBirthday())
	return &reportpb.Customer{
		// 添加该常客员工id
		StaffId: c.GetStaffId(),
		// 常客昵称
		Nickname: c.GetNickname(),
		// 常客性别
		Gender: gender,
		// 常客手机区号
		AreaCode: c.GetAreaCode(),
		// 常客手机号
		Phone: c.GetPhone(),
		// 常客生日
		Birthday: birthday,
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

// toAppCustomer
func toAppCustomer(c *customerpb.Customer) (*pb.Customer, error) {
	if c == nil {
		return nil, nil
	}
	gender, err := toAppCustomerGender(c.GetGender())
	if err != nil {
		return nil, err
	}
	birthday := toAppBirthday(c.GetBirthday())
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
		Gender: gender,
		// 常客手机号
		Phone: c.GetPhone(),
		// 常客生日
		Birthday: birthday,
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
	}, nil
}

// toSvcCustomerGender
func toSvcCustomerGender(gender pb.Gender) (customerpb.Gender, error) {
	switch gender {
	case pb.Gender_GENDER_INVALID:
		return customerpb.Gender_GENDER_INVALID, gerr.New("invalid gender")
	case pb.Gender_GENDER_UNSET:
		return customerpb.Gender_GENDER_UNSET, gerr.New("invalid gender")
	case pb.Gender_GENDER_MALE:
		return customerpb.Gender_GENDER_MALE, nil
	case pb.Gender_GENDER_FEMALE:
		return customerpb.Gender_GENDER_FEMALE, nil
	}
	return customerpb.Gender_GENDER_INVALID, gerr.New("invalid gender")
}

// toAppCustomerGender
func toAppCustomerGender(gender customerpb.Gender) (pb.Gender, error) {
	switch gender {
	case customerpb.Gender_GENDER_INVALID:
		return pb.Gender_GENDER_INVALID, gerr.New("invalid gender")
	case customerpb.Gender_GENDER_UNSET:
		return pb.Gender_GENDER_UNSET, gerr.New("invalid gender")
	case customerpb.Gender_GENDER_MALE:
		return pb.Gender_GENDER_MALE, nil
	case customerpb.Gender_GENDER_FEMALE:
		return pb.Gender_GENDER_FEMALE, nil
	}
	return pb.Gender_GENDER_INVALID, gerr.New("invalid gender")
}

// toSvcBirthday
func toSvcBirthday(b *pb.Date) *customerpb.Date {
	return &customerpb.Date{
		Year:  b.GetYear(),
		Month: b.GetMonth(),
		Day:   b.GetDay(),
	}
}

// toReportSvcBirthday
func toReportSvcBirthday(b *pb.Date) *reportpb.Date {
	return &reportpb.Date{
		Year:  b.GetYear(),
		Month: b.GetMonth(),
		Day:   b.GetDay(),
	}
}

// toAppBirthday
func toAppBirthday(b *customerpb.Date) *pb.Date {
	return &pb.Date{
		Year:  b.GetYear(),
		Month: b.GetMonth(),
		Day:   b.GetDay(),
	}
}

// toSvcCustomerPagination
func toSvcCustomerPagination(p *pb.Pagination) *customerpb.Pagination {
	return &customerpb.Pagination{
		Offset: p.GetOffset(),
		Size:   p.GetSize(),
	}
}

// toSvcGender
func toSvcGender(g pb.Gender) reportpb.Gender {
	switch g {
	case pb.Gender_GENDER_FEMALE:
		return reportpb.Gender_GENDER_FEMALE
	case pb.Gender_GENDER_MALE:
		return reportpb.Gender_GENDER_MALE
	default:
		return reportpb.Gender_GENDER_UNSET
	}
}

// toSvcLanguage
func toSvcLanguage(l pb.Language) notificationpb.Language {
	switch l {
	case pb.Language_LANGUAGE_SIMPLIFIED_CHINESE:
		return notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	case pb.Language_LANGUAGE_ENGLISH:
		return notificationpb.Language_LANGUAGE_ENGLISH
	case pb.Language_LANGUAGE_TRADITIONAL_CHINESE:
		return notificationpb.Language_LANGUAGE_TRADITIONAL_CHINESE
	default:
		return notificationpb.Language_LANGUAGE_UNSET
	}
}

// toSvcTemplate
func toSvcTemplate(l pb.TemplateAction) notificationpb.TemplateAction {
	switch l {
	// 未设置类型
	case pb.TemplateAction_TEMPLATE_ACTION_UNSET:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_UNSET
		// 组织注册
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION
		// 组织登陆
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_ORGANIZATION:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_ORGANIZATION
		// 组织重置密码
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD
		// 组织绑定手机号
	case pb.TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE
		// 组织账单
	case pb.TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING
		// 商户创建成功
	case pb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS
		// 商户创建失败
	case pb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL
		// 商户绑定手机号
	case pb.TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE
		// 商户登陆
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_TENANT:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_TENANT
		// 商户修改手机号
	case pb.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE
		// 商户重置密码
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_APP:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_APP
	case pb.TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS
	case pb.TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD:
		return notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD
	}
	return notificationpb.TemplateAction_TEMPLATE_ACTION_INVALID
}

// toSvcPagination
func toSvcPagination(p *pb.Pagination) *reportpb.Pagination {
	if p == nil {
		return nil
	}
	return &reportpb.Pagination{
		Size:   p.GetSize(),
		Offset: p.GetOffset(),
	}
}

// toAppHand
func toAppHand(h reportpb.Hand) pb.Hand {
	switch h {
	case reportpb.Hand_HAND_LEFT:
		return pb.Hand_HAND_LEFT
	case reportpb.Hand_HAND_RIGHT:
		return pb.Hand_HAND_RIGHT
	default:
		return pb.Hand_HAND_UNSET
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

// toAppHealthReport
func toAppHealthReport(r *reportpb.SimpleHealthReport) *pb.HealthReport {
	if r == nil {
		return nil
	}
	return &pb.HealthReport{
		// 报告ID
		ReportId: r.GetReportId(),
		// 创建时间
		CreateTime: r.GetCreateTime(),
		// 员工名称
		StaffName: r.GetStaffName(),
		// 员工备注
		StaffRemarks: r.GetStaffRemarks(),
		// 性别
		Gender: toAppGender(r.GetGender()),
		// 年龄
		Age: r.GetAge(),
		// 测量手
		Hand: toAppHand(r.GetHand()),
		// 是否是常客报告
		IsCustomer:          r.GetIsCustomer(),
		StressStateJudgment: r.GetStressStateJudgment(),
	}
}

// toReportTempCustomer
func toReportTempCustomer(c *pb.TempCustomer) *reportpb.TempCustomer {
	if c == nil {
		return nil
	}
	return &reportpb.TempCustomer{
		Age:    c.GetAge(),
		Gender: toSvcGender(c.GetGender()),
	}
}

// toSrvHand
func toSrvHand(hand pb.Hand) reportpb.Hand {
	switch hand {
	case pb.Hand_HAND_LEFT:
		return reportpb.Hand_HAND_LEFT
	case pb.Hand_HAND_RIGHT:
		return reportpb.Hand_HAND_RIGHT
	default:
		return reportpb.Hand_HAND_UNSET
	}
}

// toCustomerGender
func toCustomerGender(c customerpb.Gender) pb.Gender {
	switch c {
	case customerpb.Gender_GENDER_UNSET:
		return pb.Gender_GENDER_UNSET
	case customerpb.Gender_GENDER_FEMALE:
		return pb.Gender_GENDER_FEMALE
	case customerpb.Gender_GENDER_MALE:
		return pb.Gender_GENDER_MALE
	}
	return pb.Gender_GENDER_INVALID
}

// toAppCustomerLastStatus
func toAppCustomerLastStatus(c *customerpb.CustomerLastStatus) *pb.CustomerLastStatus {
	if c == nil {
		return nil
	}

	return &pb.CustomerLastStatus{
		TenantName:    c.GetTenantName(),
		CustomerName:  c.GetCustomerName(),
		CustomerPhone: c.GetCustomerPhone(),
		OverdueCount:  c.GetOverdueCount(),
		Gender:        toCustomerGender(c.GetGender()),
		Initial:       c.GetInitial(),
	}
}

// toSvcImage
func toSvcImage(p *pb.UploadingImage) *reportpb.UploadingImage {
	if p == nil {
		return nil
	}

	return &reportpb.UploadingImage{
		Mime:     p.GetMime(),
		Image:    p.GetImage(),
		Filename: p.GetFilename(),
	}
}

func toAppInquiryQuestion(i *reportpb.InquiryQuestion) *pb.InquiryQuestion {
	if i == nil {
		return nil
	}

	items := make([]*pb.InquiryAnswerItem, len(i.GetItems()))
	for k, v := range i.GetItems() {
		items[k] = toAppInquiryAnswerItem(v)
	}

	return &pb.InquiryQuestion{
		InquiryId:           i.GetInquiryId(),
		Content:             i.GetContent(),
		IsMultipleSelection: i.GetIsMultipleSelection(),
		Items:               items,
	}
}

func toAppInquiryAnswerItem(i *reportpb.InquiryAnswerItem) *pb.InquiryAnswerItem {
	if i == nil {
		return nil
	}

	return &pb.InquiryAnswerItem{
		AnswerId:  i.GetAnswerId(),
		Content:   i.GetContent(),
		Exclusive: i.GetExclusive(),
	}
}

func toSvcAnswer(a *pb.InquiryAnswer) *reportpb.InquiryAnswer {
	if a == nil {
		return nil
	}

	return &reportpb.InquiryAnswer{
		InquiryId: a.GetInquiryId(),
		Answers:   a.GetAnswers(),
	}
}
