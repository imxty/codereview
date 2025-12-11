package report

import (
	"encoding/json"
	gerr "errors"
	"regexp"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/huimaibao-service/svc/summary"
	platformpb "github.com/jinmukeji/platform-proto/gen/go/platform/report/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	YoungAge = "16-44岁"
	OldAge   = "45岁以上"
)

var (
	StressKeyMap = map[string]string{
		"has_done_sports":          "运动",
		"has_drunk_wine":           "饮酒",
		"has_had_cold":             "感冒",
		"has_rhinitis_episode":     "鼻炎",
		"has_abdominal_pain":       "腹痛腹泻",
		"has_viral_infection":      "病毒",
		"has_physiological_period": "生理期",
		"has_ovulation":            "排卵期",
		"has_pregnant":             "怀孕",
		"has_lactation":            "哺乳",
	}
)

var (
	// 脏腑辩证提取正则
	DirtyDialecticReg = regexp.MustCompile(`\[([\w\p{Han}（）]+)\]\(#[\w\.]+#\)`)
)

// mapReportBodyPartToPulseTest 把 report 里面定义的 bodyPart 类型转换为算法所使用的 finger 类型
func mapReportHandToPulseTest(hand pb.Hand) (platformpb.BodyPart, error) {
	switch hand {
	// 左手
	case pb.Hand_HAND_LEFT:
		return platformpb.BodyPart_BODY_PART_FINGER_LEFT_4, nil
	// 右手
	case pb.Hand_HAND_RIGHT:
		return platformpb.BodyPart_BODY_PART_FINGER_RIGHT_4, nil
	}
	return platformpb.BodyPart_BODY_PART_INVALID, gerr.New("invalid finger")
}

// toDomainHand
func toDomainHand(hand pb.Hand) int32 {
	switch hand {
	case pb.Hand_HAND_LEFT:
		return domain.HandLeft
	case pb.Hand_HAND_RIGHT:
		return domain.HandRight
	}
	return domain.HandRight
}

// toDomainGender
func toDomainGender(gender platformpb.Gender) int32 {
	switch gender {
	case platformpb.Gender_GENDER_FEMALE:
		return domain.GenderFemale
	case platformpb.Gender_GENDER_MALE:
		return domain.GenderMale
	}
	return domain.GenderUnset
}

// toPlatformGender
func toPlatformGender(gender customerpb.Gender) platformpb.Gender {
	switch gender {
	case customerpb.Gender_GENDER_FEMALE:
		return platformpb.Gender_GENDER_FEMALE
	case customerpb.Gender_GENDER_MALE:
		return platformpb.Gender_GENDER_MALE
	case customerpb.Gender_GENDER_UNSET:
		return platformpb.Gender_GENDER_UNSET
	default:
		return platformpb.Gender_GENDER_INVALID
	}
}

// mapTempCustomerAge 转化临时用户的年龄
func mapTempCustomerAge(age string) int {
	switch age {
	case OldAge:
		return 66
	case YoungAge:
		return 34
	}
	return 0
}

// toPlatformGenderFromReport
func toPlatformGenderFromReport(gender pb.Gender) platformpb.Gender {
	switch gender {
	case pb.Gender_GENDER_FEMALE:
		return platformpb.Gender_GENDER_FEMALE
	case pb.Gender_GENDER_MALE:
		return platformpb.Gender_GENDER_MALE
	case pb.Gender_GENDER_UNSET:
		return platformpb.Gender_GENDER_UNSET
	default:
		return platformpb.Gender_GENDER_INVALID
	}
}

// toProtoSimpleReport
func toProtoSimpleReport(r domain.ReportIntf) (*pb.SimpleHealthReport, error) {
	createdAt := timestamppb.New(r.GetCreatedAt())
	ss := ""
	if r.GetIsStressState() {
		var stressState *domain.StressState
		err := json.Unmarshal([]byte(r.GetStressState()), &stressState)
		if err != nil {
			return nil, err
		}
		ss = getStressState(stressState)
	}
	return &pb.SimpleHealthReport{
		// 报告ID
		ReportId: r.GetReportID(),
		// 创建时间
		CreateTime: createdAt,
		// 员工备注
		StaffRemarks: r.GetStaffRemarks(),
		// 性别
		Gender: toProtoGender(r.GetGender()),
		// 年龄
		Age: r.GetAge(),
		// 测量手
		Hand: toProtoHand(r.GetHand()),
		// 是否是常客报告
		IsCustomer:          r.GetIsCustomer(),
		StressStateJudgment: ss,
	}, nil
}

// toProtoGender
func toProtoGender(gender int32) pb.Gender {
	switch gender {
	case domain.GenderFemale:
		return pb.Gender_GENDER_FEMALE
	case domain.GenderMale:
		return pb.Gender_GENDER_MALE
	default:
		return pb.Gender_GENDER_INVALID
	}
}

// toProtoHand
func toProtoHand(hand int32) pb.Hand {
	switch hand {
	case domain.HandLeft:
		return pb.Hand_HAND_LEFT
	case domain.HandRight:
		return pb.Hand_HAND_RIGHT
	default:
		return pb.Hand_HAND_INVALID
	}
}

// toProtoCustomer
func toProtoCustomer(c *customerpb.Customer) *pb.ReportCustomer {
	if c == nil {
		return nil
	}
	return &pb.ReportCustomer{
		// 常客id
		CustomerId: c.GetCustomerId(),
		// 常客昵称
		Nickname: c.GetNickname(),
		// 常客首字母
		Initial: c.GetInitial(),
		// 常客性别
		Gender: toProtoGenderFromCustomer(c.GetGender()),
		// 常客手机号
		Phone: c.GetPhone(),
		// 年龄
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

// toProtoGender
func toProtoGenderFromCustomer(gender customerpb.Gender) pb.Gender {
	switch gender {
	case customerpb.Gender_GENDER_FEMALE:
		return pb.Gender_GENDER_FEMALE
	case customerpb.Gender_GENDER_MALE:
		return pb.Gender_GENDER_MALE
	default:
		return pb.Gender_GENDER_INVALID
	}
}

// toProtoTenant userpb->reportpb
func toProtoTenant(s *userpb.TenantEntity) *pb.Tenant {
	if s == nil {
		return nil
	}
	return &pb.Tenant{
		// 租户ID
		TenantId: s.GetTenantId(),
		// 租户名称
		TenantName: s.GetName(),
		// 地址
		Address: toProtoAddress(s.GetAddress()),
		// 联系电话
		Phone:   s.GetContactPhone(),
		LogoUrl: s.GetLogoUrl(),
	}
}

// toProtoAddress userpb->reportpb
func toProtoAddress(ad *userpb.Address) *pb.Address {
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

// toProtoLoopUp
func toProtoLoopUp(lookup *domain.LookUp) *pb.Lookup {
	if lookup == nil {
		return nil
	}
	content := lookup.Content
	res := DirtyDialecticReg.FindStringSubmatch(lookup.Content)
	if len(res) > 1 {
		content = res[1]
	}
	return &pb.Lookup{
		// 序号
		Key: lookup.Key,
		// 显示内容
		Content: content,
		// 得分或排序值
		Score: float64(lookup.Score),
	}
}

// toProtoLookUp
func toProtoLookUp(lookup *domain.LookUp) *pb.Lookup {
	if lookup == nil {
		return nil
	}
	content := lookup.Content
	res := DirtyDialecticReg.FindStringSubmatch(lookup.Content)
	if len(res) > 1 {
		content = res[1]
	}
	return &pb.Lookup{
		// 序号
		Key: lookup.Key,
		// 显示内容
		Content: content,
		// 得分或排序值
		Score: float64(lookup.Score),
	}
}

// toPhysiqueProtoLookUp
func toPhysiqueProtoLookUp(lookup *domain.LookUp) *pb.Lookup {
	if lookup == nil {
		return nil
	}
	content := lookup.Content
	if c, ok := summary.PhysiqueMap[lookup.Key]; ok {
		content = c.Label
	}
	return &pb.Lookup{
		// 序号
		Key: lookup.Key,
		// 显示内容
		Content: content,
		// 得分或排序值
		Score: float64(lookup.Score),
	}
}

// toProtoProduct
func toProtoProduct(product *productpb.Product) *pb.Product {
	if product == nil {
		return nil
	}
	// 症候Key
	symptomKey := ""
	if len(product.GetSymptomKeys()) != 0 {
		symptomKey = product.GetSymptomKeys()[0]
	}
	return &pb.Product{
		// 图片连接
		ImageUrl: product.GetImageUrl(),
		// 名称
		Name: product.GetProductName(),
		// 简介
		Description: product.GetDescription(),
		ProductId:   product.GetProductId(),
		// 症候
		SymptomKey: symptomKey,
	}
}

// toCustomer
func toCustomer(c *pb.Customer) *customerpb.Customer {
	if c == nil {
		return nil
	}
	return &customerpb.Customer{
		// 添加该常客员工id
		StaffId: c.GetStaffId(),
		// 常客昵称
		Nickname: c.GetNickname(),
		// 常客性别
		Gender: toCustomerGender(c.GetGender()),
		// 常客手机号
		Phone: c.GetPhone(),
		// 常客生日
		Birthday: &customerpb.Date{
			Year:  c.Birthday.Year,
			Month: c.Birthday.Month,
			Day:   c.Birthday.Day,
		},
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

// toCustomerGender
func toCustomerGender(gender pb.Gender) customerpb.Gender {
	switch gender {
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

// toGenderFromCustomer
func toGenderFromCustomer(gender customerpb.Gender) pb.Gender {
	switch gender {
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

// toReportSummary
func toReportSummary(r domain.ReportIntf) (*pb.SummaryReport, error) {
	if r == nil {
		return nil, nil
	}
	keys, err := getDirtyDialecticsKeys(r.GetDirtyDialectic())
	if err != nil {
		return nil, err
	}
	return &pb.SummaryReport{
		// 报告ID
		ReportId: r.GetReportID(),
		// 创建时间
		CreateTime: timestamppb.New(r.GetCreatedAt()),
		// 是否是常客报告
		IsCustomer: r.GetIsCustomer(),
		// 性别
		Gender: toProtoGender(r.GetGender()),
		// 脏腑辨证
		DirtyDialectics: keys,
		TenantId:        r.GetTenantID(),
	}, nil
}

// 获取脏腑辨证的key
func getDirtyDialecticsKeys(dd string) ([]string, error) {
	if dd == "" {
		return nil, nil
	}
	var dirtyDialect *domain.DirtyDialectic
	err := json.Unmarshal([]byte(dd), &dirtyDialect)
	if err != nil {
		return nil, err
	}
	keys := make([]string, len(dirtyDialect.LookUps))
	for k, v := range dirtyDialect.LookUps {
		keys[k] = v.Key
	}
	return keys, nil
}

// getStressState 生成应急态模块
func getStressState(s *domain.StressState) string {
	if !s.HasStressState {
		return ""
	}
	if s.HasDoneSports {
		return "运动"
	}
	if s.HasDrinkedWine {
		return "饮酒"
	}
	if s.HasHadCold {
		return "感冒"
	}
	if s.HasRhinitisEpisode {
		return "鼻炎"
	}
	if s.HasAbdominalPain {
		return "腹痛"
	}
	if s.HasViralInfection {
		return "病毒"
	}
	if s.HasPhysiologicalPeriod {
		return "生理期"
	}
	if s.HasOvulation {
		return "排卵期"
	}
	if s.HasPregnant {
		return "怀孕"
	}
	if s.HasLactation {
		return "哺乳期"
	}
	return ""
}

func toAliTfGender(gender int32) string {
	switch gender {
	case 2:
		return "男"
	case 3:
		return "女"
	default:
		return "男"
	}
}

func toApiInquiryAnswerItem(i domain.InquiryAnswerItem) *pb.InquiryAnswerItem {
	return &pb.InquiryAnswerItem{
		AnswerId:  i.AnswerID,
		Content:   i.Content,
		Exclusive: i.Exclusive,
	}
}
