package h5

import (
	"net/url"
	"path"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/h5/v1"
	productv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/svc/summary"
)

var (
	// 系统名称
	RiskSystemName = []string{"cardiovascular_system", "emotional_stress_system", "digestive_system", "skeletal_system"}
	// 平台风险预估系统
	RiskSystem = map[string][]string{
		// 心血管系统
		"cardiovascular_system": []string{"血糖", "血压", "血脂", "脊柱不正"},
		// 情绪压力系统
		"emotional_stress_system": []string{"焦虑", "颈椎病", "咽炎", "脑供血不足"},
		// 消化系统
		"digestive_system": []string{"免疫力", "疲劳", "肾功能减退", "心梗"},
		// 骨骼系统
		"skeletal_system": []string{"抑郁", "睡眠", "胃动力不足", "胃炎"},
	}
)

// domain->report AnswerList
func toReportAnswerList(s *pb.AnswerList) *reportpb.AnswerList {
	if s == nil {
		return nil
	}
	as := s.GetAnswers()
	if len(as) == 0 {
		return nil
	}
	reportAs := make([]*reportpb.Answer, len(as))
	for k, v := range as {
		reportAs[k] = toReportAnswer(v)
	}
	return &reportpb.AnswerList{
		Answers: reportAs,
	}
}

// app->report Answer
func toReportAnswer(s *pb.Answer) *reportpb.Answer {
	return &reportpb.Answer{
		QuestionKey: s.GetQuestionKey(),
		AnswerKeys:  s.GetAnswerKeys(),
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

// toAppUserProfile
func toAppUserProfile(profile *reportpb.UserProfile) *pb.UserProfile {
	if profile == nil {
		return nil
	}
	return &pb.UserProfile{
		// 是否是常客
		IsCustomer: profile.GetIsCustomer(),
		// 性别
		Gender: toAppGender(profile.GetGender()),
		// 手
		Hand: toAppHand(profile.GetHand()),
		// 名字
		Name: profile.GetName(),
		// 年龄
		Age: profile.GetAge(),
		// 既往病史
		Pmh: profile.GetPmh(),
	}
}

// toAppWesternMedicineReport
func toAppWesternMedicineReport(report *reportpb.WesternMedicineReport) *pb.WesternMedicineReport {
	if report == nil {
		return nil
	}
	// 构建风险系统
	estimateMap := report.GetRiskEstimate()
	var riskSystem []*pb.RiskSystem = nil
	goal := RiskSystem
	if len(estimateMap) != 0 {
		riskSystem = make([]*pb.RiskSystem, len(RiskSystemName))
		for k, v := range RiskSystemName {
			risks := goal[v]
			riskEstimate := make([]*pb.RiskEstimate, len(risks))
			for i, risk := range risks {
				riskEstimate[i] = &pb.RiskEstimate{
					// 风险名称
					RiskName: risk,
					// 风险预估预估
					RiskEstimate: estimateMap[risk].GetRiskEstimate(),
					Key:          estimateMap[risk].GetKey(),
				}
			}
			riskSystem[k] = &pb.RiskSystem{
				// 风险系统名称
				RiskSystemName: v,
				// 风险预估
				RiskEstimates: riskEstimate,
			}
		}
	}

	return &pb.WesternMedicineReport{
		// 心率(单位bpm)
		HeartRate: report.GetHeartRate(),
		// 是否有血氧
		HasSpo: report.GetHasSpo(),
		// 血氧值
		Spo: report.GetSpo(),
		// 波形图数据
		WaveData: report.GetWaveData(),
		// 风险预估模块
		RiskSystem: riskSystem,
		// 小结(富文本字符串)
		WesternMedicineSummaries: report.GetWesternMedicineSummaries(),
		// 是否是应急态，如果是，只显示血氧，心率和脉搏波数据
		IsStressState: report.GetIsStressState(),
	}
}

// toAppTCMReport
func toAppTCMReport(report *reportpb.TCMReport) *pb.TCMReport {
	if report == nil {
		return nil
	}
	//
	return &pb.TCMReport{
		// 心包经
		C0: report.GetC0(),
		// 肝经
		C1: report.GetC1(),
		// 肾经
		C2: report.GetC2(),
		// 脾经
		C3: report.GetC3(),
		// 肺经
		C4: report.GetC4(),
		// 胃经
		C5: report.GetC5(),
		// 胆经
		C6: report.GetC6(),
		// 膀胱经
		C7: report.GetC7(),
		// 脏腑辩证模块
		DirtyDialecticModule:     toAppDirtyDialecticModule(report.GetDirtyDialecticModule()),
		PhysiqueDialecticsModule: toAppPhysiqueModule(report.GetPhysiqueDialecticsModule()),
		// 中医报告小结(富文本字段)
		TcmSummaries: report.GetTcmSummaries(),
		// 如果is_stress_state为true，该模块才有值
		StressStateModule: toAppStressStateModule(report.GetStressStateModule()),
		// 是否是应急态，如果是，只显示柱状图和应急态状态
		IsStressState: report.GetIsStressState(),
	}
}

// toAppLookUp
func toAppLookUp(lookup *reportpb.Lookup) *pb.Lookup {
	if lookup == nil {
		return nil
	}
	return &pb.Lookup{
		// 序号
		Key: lookup.GetKey(),
		// 显示内容
		Content: lookup.GetContent(),
		// 得分或排序值
		Score: lookup.GetScore(),
		// 词条链接 Key
		LinkKey: lookup.GetLinkKey(),
	}
}

// toAppDirtyDialecticModule
func toAppDirtyDialecticModule(module *reportpb.DirtyDialecticModule) *pb.DirtyDialecticModule {
	if module == nil {
		return nil
	}
	lookups := make([]*pb.Lookup, len(module.GetLookups()))
	for k, v := range module.GetLookups() {
		lookups[k] = toAppLookUp(v)
	}
	return &pb.DirtyDialecticModule{
		Lookups: lookups,
	}
}

// toAppPhysiqueModule
func toAppPhysiqueModule(module *reportpb.PhysiqueDialecticsModule) *pb.PhysiqueDialecticsModule {
	if module == nil {
		return nil
	}
	lookups := make([]*pb.Lookup, len(module.GetLookups()))
	for k, v := range module.GetLookups() {
		lookups[k] = toAppLookUp(v)
	}
	return &pb.PhysiqueDialecticsModule{
		Lookups: lookups,
	}
}

// toAppStressStateModule
func toAppStressStateModule(module *reportpb.StressStateModule) *pb.StressStateModule {
	if module == nil {
		return nil
	}
	lookups := make([]*pb.Lookup, len(module.GetLookups()))
	for k, v := range module.GetLookups() {
		lookups[k] = toAppLookUp(v)
	}
	return &pb.StressStateModule{
		Lookups: lookups,
	}
}

// toAppPhysicalTherapyReport
func toAppPhysicalTherapyReport(report *reportpb.PhysicalTherapyReport) *pb.PhysicalTherapyReport {
	if report == nil {
		return nil
	}
	return &pb.PhysicalTherapyReport{
		// 阴虚指数
		F0: report.GetF0(),
		// 阳虚指数
		F1: report.GetF1(),
		// 湿气指数
		F2: report.GetF2(),
		// 血瘀指数
		F3: report.GetF3(),
		// 阴阳平衡指数
		F4: report.GetF4(),
	}
}

// toAppDietaryAdviceModule
func toAppDietaryAdviceModule(report *reportpb.DietaryAdviceModule) *pb.DietaryAdviceModule {
	if report == nil {
		return nil
	}
	lookups := make([]*pb.Lookup, len(report.GetLookups()))
	for k, v := range report.GetLookups() {
		lookups[k] = toAppLookUp(v)
	}
	return &pb.DietaryAdviceModule{
		Lookups: lookups,
	}
}

// toAppChineseMedicineAdviceModule
func toAppChineseMedicineAdviceModule(report *reportpb.ChineseMedicineAdviceModule) *pb.ChineseMedicineAdviceModule {
	if report == nil {
		return nil
	}
	lookups := make([]*pb.Lookup, len(report.GetLookups()))
	for k, v := range report.GetLookups() {
		lookups[k] = toAppLookUp(v)
	}
	return &pb.ChineseMedicineAdviceModule{
		Lookups: lookups,
	}
}

// toAppPhysicalAdviceModule
func toAppPhysicalAdviceModule(report *reportpb.PhysicalTherapyAdviceModule) *pb.PhysicalTherapyAdviceModule {
	if report == nil {
		return nil
	}
	lookups := make([]*pb.Lookup, len(report.GetLookups()))
	for k, v := range report.GetLookups() {
		lookups[k] = toAppLookUp(v)
	}
	return &pb.PhysicalTherapyAdviceModule{
		Lookups: lookups,
	}
}

// toAppSportsAdviceModule
func toAppSportsAdviceModule(report *reportpb.SportsAdviceModule) *pb.SportsAdviceModule {
	if report == nil {
		return nil
	}
	lookups := make([]*pb.Lookup, len(report.GetLookups()))
	for k, v := range report.GetLookups() {
		lookups[k] = toAppLookUp(v)
	}
	return &pb.SportsAdviceModule{
		Lookups: lookups,
	}
}

// toAppProduct
func toAppProduct(report *reportpb.Product, s3Domain string) *pb.Product {
	if report == nil {
		return nil
	}
	imgUrl := ""
	if report.GetImageUrl() != "" {
		link, _ := url.Parse(s3Domain)
		link.Path = path.Join(link.Path, report.GetImageUrl())
		imgUrl = link.String()
	}
	return &pb.Product{
		// 图片连接
		ImageUrl: imgUrl,
		// 名称
		Name: report.GetName(),
		// 简介
		Description: report.GetDescription(),
		ProductId:   report.GetProductId(),
		SymptomKey:  report.GetSymptomKey(),
	}
}

// toAppCInfos
func toAppCInfos(infos []*reportpb.CInfo) []*pb.CInfo {
	ins := make([]*pb.CInfo, len(infos))
	for k, v := range infos {
		ins[k] = &pb.CInfo{
			C0:         v.GetC0(),
			C1:         v.GetC1(),
			C2:         v.GetC2(),
			C3:         v.GetC3(),
			C4:         v.GetC4(),
			C5:         v.GetC5(),
			C6:         v.GetC6(),
			C7:         v.GetC7(),
			CreateTime: v.GetCreateTime(),
		}
	}
	return ins
}

// toAppStressStateCount
func toAppStressStateCount(s *reportpb.StressStateCount) *pb.StressStateCount {
	if s == nil {
		return nil
	}
	lookups := make([]*pb.Lookup, len(s.GetLookups()))
	for k, v := range s.GetLookups() {
		lookups[k] = toAppLookUp(v)
	}

	return &pb.StressStateCount{
		Lookups: lookups,
	}
}

// toAppRecommendations
func toAppRecommendations(report *reportpb.Recommendation, s3Domain string) *pb.Recommendation {
	if report == nil {
		return nil
	}
	products := make([]*pb.Product, len(report.GetProducts()))
	for k, v := range report.GetProducts() {
		products[k] = toAppProduct(v, s3Domain)
	}
	return &pb.Recommendation{
		// 中药调理
		Products: products,
		// 食疗建议
		DietaryAdvice: toAppDietaryAdviceModule(report.GetDietaryAdvice()),
		// 运动建议
		SportsAdvice: toAppSportsAdviceModule(report.GetSportsAdvice()),
		// 金姆和租户方案是否都是最新版
		AreLatestTreatments: report.GetAreLatestTreatments(),
		// 是否是应急态，如果是，只显示柱状图和应急态状态
		IsStressState: report.GetIsStressState(),
		// 中药调理
		ChineseMedicineAdvice: toAppChineseMedicineAdviceModule(report.GetChineseMedicineAdvice()),
		// 调理建议
		PhysicalTherapyAdvice: toAppPhysicalAdviceModule(report.GetPhysicalTherapyAdvice()),
	}
}

// toAppAddress
func toAppAddress(address *reportpb.Address) *pb.Address {
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

// toAppTenant
func toAppTenant(t *reportpb.Tenant, s3Domain string) *pb.Tenant {
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
		Phone: t.GetPhone(),
	}
}

// toAppReport
func toAppReport(report *reportpb.HealthReport, s3Domain string) *pb.HealthReport {
	if report == nil {
		return nil
	}
	// 判断测量异常模块
	if !judgeMeasurementModule(report.GetMeasurementJudgment()) {
		return &pb.HealthReport{
			// 报告ID
			ReportId: report.GetReportId(),
			// 用户档案
			UserProfile: toAppUserProfile(report.GetUserProfile()),
			// 测量备注(给测量者看)
			Remarks: report.GetRemarks(),
			// 员工备注
			StaffRemarks: report.GetStaffRemarks(),
			// 租户信息
			Tenant: toAppTenant(report.GetTenant(), s3Domain),
			// 创建时间
			CreateTime: report.GetCreateTime(),
			// 测量判断模块
			MeasurementJudgment: generateMeasurementJudgment(report.GetMeasurementJudgment()),
		}
	}
	// 返回报告
	return &pb.HealthReport{
		// 报告ID
		ReportId: report.GetReportId(),
		// 用户档案
		UserProfile: toAppUserProfile(report.GetUserProfile()),
		// 创建时间
		CreateTime: report.GetCreateTime(),
		// 西医报告
		WesternMedicineReport: toAppWesternMedicineReport(report.GetWesternMedicineReport()),
		// 中医报告
		TcmReport: toAppTCMReport(report.GetTcmReport()),
		// 理疗报告
		PhysicalTherapyReport: toAppPhysicalTherapyReport(report.GetPhysicalTherapyReport()),
		// 调理建议
		Recommendation: toAppRecommendations(report.GetRecommendation(), s3Domain),
		// 测量备注(给测量者看)
		Remarks: report.GetRemarks(),
		// 员工备注
		StaffRemarks: report.GetStaffRemarks(),
		// 租户信息
		Tenant: toAppTenant(report.GetTenant(), s3Domain),
		// 测量判断模块
		MeasurementJudgment: generateMeasurementJudgment(report.GetMeasurementJudgment()),
	}
}

// reportpb->app QuestionChoice
func toAppQuestionChoice(c *reportpb.QuestionChoice) *pb.QuestionChoice {
	if c == nil {
		return nil
	}
	return &pb.QuestionChoice{
		ChoiceKey:    c.GetChoiceKey(),
		Content:      c.GetContent(),
		ConflictKeys: c.GetConflictKeys(),
		Selected:     c.GetSelected(),
	}
}

// reportpb->app Question
func toAppQuestion(q *reportpb.Question) *pb.Question {
	if q == nil {
		return nil
	}
	choices := make([]*pb.QuestionChoice, len(q.GetChoices()))
	for k, v := range q.GetChoices() {
		choices[k] = toAppQuestionChoice(v)
	}
	return &pb.Question{
		QuestionKey: q.GetQuestionKey(),
		Content:     q.GetContent(),
		Type:        q.GetType(),
		Choices:     choices,
	}
}

// reportpb->app QuestionList
func toAppQuestionList(qs *reportpb.QuestionList) *pb.QuestionList {
	if qs == nil {
		return nil
	}
	s := make([]*pb.Question, len(qs.GetQuestions()))
	for k, v := range qs.GetQuestions() {
		s[k] = toAppQuestion(v)
	}
	return &pb.QuestionList{
		Questions: s,
	}
}

// 判断测量模块
func judgeMeasurementModule(m *reportpb.MeasurementJudgmentModule) bool {
	const (
		AbnormalKey = "C0001"
	)
	if m != nil && len(m.Lookups) != 0 {
		// 如果key是C0001，直接返回提示信息其他都不返回
		if m.Lookups[0].Key == AbnormalKey {
			return false
		}
	}
	return true
}

// generateMeasurementJudgment 生成测量判断模块
func generateMeasurementJudgment(measurementJudgment *reportpb.MeasurementJudgmentModule) *pb.MeasurementJudgmentModule {
	if measurementJudgment == nil {
		return nil
	}
	lookUps := make([]*pb.Lookup, len(measurementJudgment.Lookups))
	for k, v := range measurementJudgment.Lookups {
		lookUps[k] = &pb.Lookup{
			// 序号
			Key: v.GetKey(),
			// 显示内容
			Content: v.GetContent(),
			// 得分或排序值
			Score: v.GetScore(),
			// 词条链接 Key
			LinkKey: v.GetLinkKey(),
		}
	}
	return &pb.MeasurementJudgmentModule{
		Lookups: lookUps,
	}
}

// tAppProductType
func tAppProductType(productType productv1.ProductType) pb.ProductType {
	switch productType {
	case productv1.ProductType_PRODUCT_TYPE_UNSET:
		return pb.ProductType_PRODUCT_TYPE_UNSET
	case productv1.ProductType_PRODUCT_TYPE_SERVICE:
		return pb.ProductType_PRODUCT_TYPE_SERVICE
	case productv1.ProductType_PRODUCT_TYPE_CPD:
		return pb.ProductType_PRODUCT_TYPE_CPD
	case productv1.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT:
		return pb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT
	case productv1.ProductType_PRODUCT_TYPE_NUTRITION:
		return pb.ProductType_PRODUCT_TYPE_NUTRITION
	}
	return pb.ProductType_PRODUCT_TYPE_INVALID
}

// toReportProduct
func toReportProduct(p *productv1.Product, s3Domain string) *pb.ProductDetail {
	if p == nil {
		return nil
	}
	keys := p.GetSymptomKeys()
	keyName := make([]string, len(keys))
	for k, v := range keys {
		if summary.AEKeyMap[v].Label != "" {
			keyName[k] = summary.AEKeyMap[v].Label
		}
	}
	imgUrl := ""
	if p.GetImageUrl() != "" {
		link, _ := url.Parse(s3Domain)
		link.Path = path.Join(link.Path, p.GetImageUrl())
		imgUrl = link.String()
	}
	return &pb.ProductDetail{
		// 商品id
		ProductId: p.GetProductId(),
		// 商品类型
		ProductType: tAppProductType(p.GetProductType()),
		// 商品名称
		ProductName: p.GetProductName(),
		// 商品介绍
		Description: p.GetDescription(),
		// 商品图片
		ImageUrl: imgUrl,
		// 商品备注
		Remarks: p.GetRemarks(),
		// 药品名称(商品类型为理疗服务方案不用传)
		DrugName: p.GetDrugName(),
		// 药品是非处方药 (商品类型为理疗服务方案不用传)
		IsOtc: p.GetIsOtc(),
		// 准字号 (商品类型为理疗服务方案不用传)
		ApprovedNumber: p.GetApprovedNumber(),
		// 药品准效期 (商品类型为理疗服务方案不用传)
		DrugValidityPeriod: p.GetDrugValidityPeriod(),
		// 商品跳转链接
		Link: p.GetLink(),
		// 适应的症候
		Symptoms: keyName,
	}
}
