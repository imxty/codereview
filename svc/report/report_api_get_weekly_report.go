package report

import (
	"context"
	gerr "errors"
	"strconv"
	"strings"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	platformpb "github.com/jinmukeji/platform-proto/gen/go/platform/report/v1"
)

const (
	// 平台报告数量不足
	PlatformErrReportNotEnough = "1107"
)

// 添加常客报告请求
func (s *ReportAPIHandler) GetWeeklyReport(ctx context.Context, req *pb.GetWeeklyReportRequest, rsp *pb.GetWeeklyReportResponse) error {

	// 1.验证request
	err := validateGetWeeklyReportRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 构建用户档案
	getCustomerRsp, err := s.customerAPI.GetCustomer(ctx, &customerpb.GetCustomerRequest{
		TenantId:   req.GetTenantId(),
		CustomerId: req.GetCustomerId(),
	})
	if err != nil {
		return err
	}
	customer := getCustomerRsp.GetCustomer()

	userProfile := &pb.UserProfile{
		// 是否是常客
		IsCustomer: true,
		// 性别
		Gender: toGenderFromCustomer(customer.GetGender()),
		// 名字
		Name: customer.GetNickname(),
		// 年龄
		Age: strconv.Itoa(int(customer.GetAge())),
		// 既往病史
		Pmh: customer.GetPmh(),
	}

	// 获取周报告
	ctx = buildPlatformContext(ctx)
	getRsp, err := s.platformClient.GetWeeklyReport(ctx, &platformpb.GetWeeklyReportRequest{
		AppId:                 s.appId,
		SubjectId:             req.GetCustomerId(),
		RequestedModuleInputs: getWeeklyRequestModules(),
		TemplateParams: &platformpb.ReportTemplateParams{
			LanguageCode: LanguageSimpleChinese,
		},
	})

	// 判断错误
	if err != nil {
		if strings.Contains(err.Error(), PlatformErrReportNotEnough) {
			rsp.Report = &pb.HealthReport{
				UserProfile: userProfile,
			}
			rsp.ReportEnough = false
			return nil
		}
		return err
	}

	rsp.ReportEnough = true

	// 如果没有问题就去解析ctxData
	modulesResult := getRsp.GetReport().GetModuleResults()
	// 如果没有问题就去解析ctxData
	dirtyDialect, err := domain.GetDirtyDialectic(modulesResult, false)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	diseases, err := domain.GetDiseaseEstimate(modulesResult, false)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	factor, err := domain.GetFactorInterpretation(modulesResult)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取状图模块
	meridianBarChart, err := domain.GetMeridianBarChartModule(modulesResult)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取周报
	weeklyReport, err := domain.GetWeeklyReport(modulesResult)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 返回应激态报告
	stressState := new(domain.StressState)
	rsp.Report = &pb.HealthReport{
		UserProfile:           userProfile,
		WesternMedicineReport: s.getWeeklyWesternMedicineReport(diseases),
		TcmReport:             s.getWeeklyTCMReport(stressState.HasStressState, meridianBarChart, dirtyDialect, weeklyReport, stressState),
		PhysicalTherapyReport: s.getPhysicalTherapyReport(factor),
	}

	return nil
}

// 验证request
func validateGetWeeklyReportRequest(req *pb.GetWeeklyReportRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomerId() == "" {
		return gerr.New("customer_id should not be empty")
	}
	return nil
}

// generateWeeklyPhysiqueDialectic 生成周报体质辩证
func generateWeeklyPhysiqueDialectic(physique *platformpb.WeeklyReport) *pb.PhysiqueDialecticsModule {
	lookUps := make([]*pb.Lookup, len(physique.GetPhysicalDialecticsStatSummary()))
	for k, v := range physique.GetPhysicalDialecticsStatSummary() {
		content := v.GetContent()
		res := DirtyDialecticReg.FindStringSubmatch(content)
		if len(res) > 1 {
			content = res[1]
		}
		lookUps[k] = &pb.Lookup{
			// 序号
			Key: v.GetKey(),
			// 显示内容
			Content: content,
			// 得分或排序值
			Score: v.GetScore(),
			// 词条链接 Key
			LinkKey: v.GetLinkKey(),
		}
	}
	return &pb.PhysiqueDialecticsModule{
		Lookups: lookUps,
	}
}

// getWeeklyRequestModules 获取周报请求模块
func getWeeklyRequestModules() []*platformpb.ReportModuleInput {
	// - heart 心率血氧模块
	// - risk_estimate 风险预估模块
	// - physical_dialectics 体质辩证模块
	// - dirty_dialectic 脏腑辩证模块
	// - dietary_advice 食疗建议模块
	// - sports_advice 运动方案模块
	// - chinese_medicine_advice 中药调理建议模块
	// - physical_therapy_advice 理疗建议模块
	// - meridian_bar_chart 经络柱状图模块
	// - factor_interpretation 阴阳模块
	// - measurement_judgment 测量异常判断模块
	// - stress_state_judgment 应激态模块
	// - health_index 健康指数模块
	// - drunk_judgment 饮酒判断模块
	return []*platformpb.ReportModuleInput{
		{
			Name: "risk_estimate",
		},
		{
			Name: "physical_dialectics",
		},
		{
			Name: "dirty_dialectic",
		},
		{
			Name: "dietary_advice",
		},
		{
			Name: "sports_advice",
		},
		{
			Name: "chinese_medicine_advice",
		},
		{
			Name: "physical_therapy_advice",
		},
		{
			Name: "meridian_bar_chart",
		},
		{
			Name: "factor_interpretation",
		},
		{
			Name: "weekly_report",
		},
	}
}

// getWeeklyWesternMedicineReport 获取西医报告
func (u *ReportAPIHandler) getWeeklyWesternMedicineReport(diseases *domain.DiseaseEstimate) *pb.WesternMedicineReport {
	wesReport := &pb.WesternMedicineReport{}
	riskEstimates, riskSummary := generateRiskEstimate(diseases)
	wesReport.RiskEstimate = riskEstimates
	wesReport.WesternMedicineSummaries = riskSummary
	return wesReport
}

// getWeeklyTCMReport 获取中医报告
func (u *ReportAPIHandler) getWeeklyTCMReport(isStressState bool, chart *domain.MeridianBarChart, dirtyDialect *domain.DirtyDialectic, physique *platformpb.WeeklyReport, stressState *domain.StressState) *pb.TCMReport {
	tcmReport := &pb.TCMReport{
		// 心包经
		C0: chart.C0,
		// 肝经
		C1: chart.C1,
		// 肾经
		C2: chart.C2,
		// 脾经
		C3: chart.C3,
		// 肺经
		C4: chart.C4,
		// 胃经
		C5: chart.C5,
		// 胆经
		C6: chart.C6,
		// 膀胱经
		C7: chart.C7,
		// 是否是应急态，如果是，只显示柱状图和应急态状态
		IsStressState: isStressState,
	}

	if !isStressState {
		dialect := generateDirtyDialectic(dirtyDialect)
		physiqueDialects := generateWeeklyPhysiqueDialectic(physique)
		tcmReport.DirtyDialecticModule = dialect
		tcmReport.PhysiqueDialecticsModule = physiqueDialects
		// 中医报告小结(富文本字段)
		tcmReport.TcmSummaries = u.generateTcmSummary(dialect.Lookups)
	} else {
		tcmReport.StressStateModule = u.genderStressStateModule(stressState)
	}

	return tcmReport
}
