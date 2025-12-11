package report

import (
	"context"
	"encoding/json"
	gerr "errors"
	"fmt"
	"strconv"
	"strings"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/huimaibao-service/svc/summary"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	platformpb "github.com/jinmukeji/platform-proto/gen/go/platform/report/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Address
type Address struct {
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
}

const (
	// ErrReportUpdateFailed
	ErrReportUpdateFailed = 3209
)

const (
	// LanguageSimpleChinese 简体中文
	LanguageSimpleChinese = "zh-Hans"
	// 中国国家编码
	CountryCodeChina = 0
)

const (
	// PulseTestPostureSetting 坐姿
	PulseTestPostureSetting int32 = 0
)

// GetReport 获取报告
func (u *ReportAPIHandler) GetReport(ctx context.Context, req *pb.GetReportRequest, rsp *pb.GetReportResponse) error {

	// 1.验证request
	err := validateGetReportRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 1. 获取报告
	report, err := u.reportStore.GetReportByID(ctx, req.GetTenantId(), req.GetReportId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if report == nil {
		return errors.Errorf(ErrReportNotFound, "report not found by report_id %s", req.GetReportId())
	}
	reportRev := report.GetRev()
	// 获取租户信息
	getTenantRsp, err := u.userAPI.GetTenant(ctx, &userpb.GetTenantRequest{
		OrganizationId: report.GetOrganizationID(),
		TenantId:       req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}
	tenant := toProtoTenant(getTenantRsp.GetTenant())
	// 报告的内容
	// 如果有应急态，那么dirtyDialect，diseases，factor，recommend都为nil
	var stressState *domain.StressState
	var dirtyDialect *domain.DirtyDialectic
	var physique *domain.PhysiqueDialectics
	var diseases *domain.DiseaseEstimate
	var factor *domain.FactorInterpretation
	var recommend *domain.Recommendations
	var measurementJudgment *domain.MeasurementJudgmentModule
	var meridianBarChart *domain.MeridianBarChart
	var heartModule *platformpb.HeartModuleOutput
	isFirstReported := false
	// 报告创建时间
	createTime := timestamppb.New(report.GetCreatedAt())
	// 获取报告是否完成
	if !report.GetIsReportComplete() {
		// 解析回答
		answers := getAnswers(req.GetModuleAnswers())
		ctx = buildPlatformContext(ctx)
		// 获取报告 相关地址
		ad, err := u.reportStore.GetReportAddress(ctx, req.GetReportId())
		if err != nil {
			return errors.Errorf(codes.DataAccessFailed, "failed to get report address [%s]", err.Error())
		}
		address := new(Address)
		// 只有中国的地理位置信息才统计
		if ad != nil && ad.GetCountryCode() == CountryCodeChina {
			address.Province = ad.GetProvince()
			address.City = ad.GetCity()
			address.District = ad.GetDistrict()
		}
		// 序列化
		res, err := json.Marshal(address)
		if err != nil {
			return errors.Errorf(codes.DataAccessFailed, "failed to marshall address [%s]", err.Error())
		}
		getReportRsp, err := u.platformClient.GetReport(ctx, &platformpb.GetReportRequest{
			// app_id
			AppId: u.appId,
			// 报告ID
			ReportId: report.GetReportID(),
			// 申请查询的模块内容清单
			// 返回的报告内容中仅包含申请的模块清单的数据
			// 必填
			RequestedModuleInputs: getRequestModules(),
			// 回答的问题，是模块名到回答的问题的映射关系
			// 有问题就要填
			ModuleAnswers: answers,
			// 默认使用moving average
			Options: &platformpb.ReportOptions{
				MovingAverage: true,
			},
			// 分析报告内容模版参数
			// language_code必填
			TemplateParams: &platformpb.ReportTemplateParams{
				// 默认汉语
				LanguageCode: LanguageSimpleChinese,
				// 额外参数
				ExtraParams: map[string]string{
					"assistant_name": "小慧",
					"device_name":    "检测仪",
				},
			},
			// 添加地址信息
			Extras: map[string]string{
				"address": string(res),
			},
		})
		if err != nil {
			return errors.Errorf(ErrPlatformError, "Platform Error:[%s]", err.Error())
		}

		// 如果有问题
		if getReportRsp.GetHasQuestions() {
			rsp.IsCompleteReport = false
			rsp.HasQuestions = true
			rsp.Report = &pb.HealthReport{
				ReportId: req.GetReportId(),
			}
			rsp.ModuleQuestions = parseOutputQuestions(getReportRsp.GetModuleQuestions())
			return nil
		} else {
			// 如果没有问题就去解析ctxData
			modulesResult := getReportRsp.GetReport().GetModuleResults()
			// 如果没有问题就去解析ctxData
			// 应激态模块
			stressState, err = domain.GetStressState(modulesResult)
			if err != nil {
				return errors.Error(codes.InvalidRequest, err.Error())
			}
			heartModule, err = domain.GetHeartModule(modulesResult)
			if err != nil {
				return errors.Error(codes.InvalidRequest, err.Error())
			}
			dirtyDialect, err = domain.GetDirtyDialectic(modulesResult, stressState.HasStressState)
			if err != nil {
				return errors.Error(codes.InvalidRequest, err.Error())
			}
			diseases, err = domain.GetDiseaseEstimate(modulesResult, stressState.HasStressState)
			if err != nil {
				return errors.Error(codes.InvalidRequest, err.Error())
			}
			factor, err = domain.GetFactorInterpretation(modulesResult)
			if err != nil {
				return errors.Error(codes.InvalidRequest, err.Error())
			}
			recommend, err = domain.GetRecommendations(modulesResult, stressState.HasStressState)
			if err != nil {
				return errors.Error(codes.InvalidRequest, err.Error())
			}
			// 获取测量判断模块
			measurementJudgment, err = domain.GetMeasurementModule(modulesResult)
			if err != nil {
				return errors.Error(codes.InvalidRequest, err.Error())
			}
			// 获取中医建议
			physique, err = domain.GetPhysiqueDialectics(modulesResult, stressState.HasStressState)
			if err != nil {
				return errors.Errorf(codes.InvalidRequest, "get physique failed[%s]", err.Error())
			}
			// 获取状图模块
			meridianBarChart, err = domain.GetMeridianBarChartModule(modulesResult)
			if err != nil {
				return errors.Error(codes.InvalidRequest, err.Error())
			}
			// 解析成json保存到数据库
			stressStateJson, _ := json.Marshal(stressState)
			dirtyDialectJson, _ := json.Marshal(dirtyDialect)
			diseasesJson, _ := json.Marshal(diseases)
			recommendJson, _ := json.Marshal(recommend)
			physiqueJson, _ := json.Marshal(physique)
			measurementJson, _ := json.Marshal(measurementJudgment)
			var dirtyDialectFlag int32
			// 获取脏腑辨证flag
			if dirtyDialect != nil && len(dirtyDialect.LookUps) != 0 {
				keys := make([]string, len(dirtyDialect.LookUps))
				for k, v := range dirtyDialect.LookUps {
					keys[k] = v.Key
				}
				dirtyDialectFlag, err = getDDKeyFlag(keys)
				if err != nil {
					return errors.Error(codes.InvalidRequest, err.Error())
				}
			}
			err = u.reportStore.UpdateReportContent(ctx, req.GetReportId(), req.GetTenantId(), heartModule.GetAverageHeartRate(), stressState.HasStressState, factor,
				string(stressStateJson), string(dirtyDialectJson), string(diseasesJson), string(recommendJson), string(physiqueJson), string(measurementJson), meridianBarChart, dirtyDialectFlag, reportRev)
			if err != nil {
				return errors.Error(ErrReportUpdateFailed, err.Error())
			}
			// 如果是常客测量更新最新状态
			if report.GetIsCustomer() {
				_, err = u.customerAPI.UpdateCustomerLastStatus(ctx, &customerpb.UpdateCustomerLastStatusRequest{
					CustomerId: report.GetCustomerID(),
					ReportId:   req.GetReportId(),
				})
				if err != nil {
					return err
				}
			}
			// 版本加1
			reportRev++
			isFirstReported = true
		}
	} else {
		// 如果报告已经完成直接json去解析数据即可
		// 老报告没有测量异常模块
		if report.GetMeasurementJudgment() != "" {
			err = json.Unmarshal([]byte(report.GetMeasurementJudgment()), &measurementJudgment)
			if err != nil {
				return errors.Errorf(codes.InvalidRequest, "unmarshal measurementJudgment failed [%s]", err.Error())
			}
		}
		meridianBarChart = &domain.MeridianBarChart{
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
		}
		heartModule = &platformpb.HeartModuleOutput{
			// 血氧
			Spo: &wrapperspb.Int32Value{
				Value: report.GetSpo(),
			},
			// 平均心率
			AverageHeartRate: report.GetHeartRate(),
		}
		// 如果没有异常才去继续分析
		if judgeMeasurementModule(measurementJudgment) {
			err = json.Unmarshal([]byte(report.GetStressState()), &stressState)
			if err != nil {
				return errors.Errorf(codes.InvalidRequest, "unmarshal stressState failed: [%s]", err.Error())
			}
			// 如果不是应急态才需要去解析
			if !stressState.HasStressState {
				err = json.Unmarshal([]byte(report.GetDirtyDialectic()), &dirtyDialect)
				if err != nil {
					return errors.Error(codes.InvalidOperation, err.Error())
				}
				err = json.Unmarshal([]byte(report.GetRisks()), &diseases)
				if err != nil {
					return errors.Error(codes.InvalidOperation, err.Error())
				}
				err = json.Unmarshal([]byte(report.GetDietaryAdvice()), &recommend)
				if err != nil {
					return errors.Error(codes.InvalidOperation, err.Error())
				}
				err = json.Unmarshal([]byte(report.GetPhysicalDialectics()), &physique)
				if err != nil {
					return errors.Error(codes.InvalidOperation, err.Error())
				}
			}
			factor = &domain.FactorInterpretation{
				// 阴
				F0: report.GetF0(),
				// 阳
				F1: report.GetF1(),
				// 湿气
				F2: report.GetF2(),
				// 血瘀
				F3: report.GetF3(),
				// 阴阳平衡
				F4: report.GetF4(),
			}
		}
	}
	// 构建用户档案
	userProfile := &pb.UserProfile{
		// 是否是常客
		IsCustomer: report.GetIsCustomer(),
		// 性别
		Gender: toProtoGender(report.GetGender()),
		// 手
		Hand: toProtoHand(report.GetHand()),
		// 年龄
		Age: report.GetAge(),
	}
	// 如果是常客的测量，去查询常客的档案
	if report.GetIsCustomer() {
		getCustomerRsp, err := u.customerAPI.GetCustomer(ctx, &customerpb.GetCustomerRequest{
			TenantId:   report.GetTenantID(),
			CustomerId: report.GetCustomerID(),
		})
		if err != nil {
			return err
		}
		customer := getCustomerRsp.GetCustomer()
		userProfile.Name = customer.GetNickname()
		userProfile.Pmh = customer.GetPmh()
	}

	rsp.IsCompleteReport = true
	// 如果hr_cv>10 直接只有测量异常模块
	if !judgeMeasurementModule(measurementJudgment) {
		// 返回异常报告
		rsp.Report = &pb.HealthReport{
			ReportId:    req.GetReportId(),
			UserProfile: userProfile,
			CreateTime:  createTime,
			Tenant:      tenant,
			// 备注
			Remarks:             report.GetRemarks(),
			StaffRemarks:        report.GetStaffRemarks(),
			MeasurementJudgment: generateMeasurementJudgment(measurementJudgment),
		}
		return nil
	}
	// 获取报告对应的商品
	var dirtyDialecticReportResult []*productpb.ReportResult
	var highRiskyDiseaseReportResult, midRiskyDiseaseReportResult []*productpb.ReportResult
	// 非应急态解析
	if !stressState.HasStressState {
		// 中医报告脏腑辩证结果
		dirtyDialecticReportResult = make([]*productpb.ReportResult, len(dirtyDialect.LookUps))
		for k, v := range dirtyDialect.LookUps {
			dirtyDialecticReportResult[k] = &productpb.ReportResult{
				SymptomKey: v.Key,
			}
		}
		// 西医报告高风险疾病结果
		// 西医报告中风险疾病结果
		highRiskyDiseaseReportResult, midRiskyDiseaseReportResult = generateHighAndMidRiskEstimate(diseases)
	}
	listProductsRsp, err := u.productAPI.ListRecommendedProducts(ctx, &productpb.ListRecommendedProductsRequest{
		// 租户id
		TenantId: req.GetTenantId(),
		// 租户方案版本ID
		TenantTreatmentRevId: report.GetTenantTreatmentRevID(),
		// 报告是否初次生成
		IsFirstReported: isFirstReported,
		// 中医报告脏腑辩证结果
		DirtyDialecticReportResult: dirtyDialecticReportResult,
		// 西医报告高风险疾病结果
		HighRiskyDiseaseReportResult: highRiskyDiseaseReportResult,
		// 西医报告中风险疾病结果
		MediumRiskyDiseaseReportResult: midRiskyDiseaseReportResult,
		// 理疗报告结果
		PhysicalTherapyReportResult: generatePhysicalTherapyReportResult(factor.F2, factor.F3),
		// 报告ID
		ReportId: req.GetReportId(),
		// 体质结果
		PhysiqueDialecticsReportResult: generatePhysique(physique),
	})
	if err != nil {
		return err
	}

	// 判断是否存在方案
	hasTreatment := false
	// 如果有id说明报告当前有方案
	hasTreatment = report.GetTenantTreatmentRevID() != ""

	// 如果是初次获取商品需要记录treatment_rev_id
	if isFirstReported {
		err = u.reportStore.UpdateReportTreatment(ctx, req.GetTenantId(), report.GetReportID(), listProductsRsp.GetTenantTreatmentRevId(), reportRev)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 如果是首次则通过返回值判断
		hasTreatment = listProductsRsp.GetTenantTreatmentRevId() != ""
	}

	// 如果没有问题，就完善报告的状态，并返回报告的内容
	// 获取所有报告状态并且json序列化到数据库
	// 如果是应激态
	wesReport, err := u.getWesternMedicineReport(ctx, stressState.HasStressState, req.GetTenantId(), report, heartModule, diseases)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 返回应激态报告
	rsp.Report = &pb.HealthReport{
		ReportId:              req.GetReportId(),
		UserProfile:           userProfile,
		CreateTime:            createTime,
		WesternMedicineReport: wesReport,
		TcmReport:             u.getTCMReport(stressState.HasStressState, meridianBarChart, dirtyDialect, physique, stressState),
		PhysicalTherapyReport: u.getPhysicalTherapyReport(factor),
		Recommendation:        u.getRecommendation(stressState.HasStressState, recommend, listProductsRsp, hasTreatment),
		Tenant:                tenant,
		// 备注
		Remarks:             report.GetRemarks(),
		StaffRemarks:        report.GetStaffRemarks(),
		MeasurementJudgment: generateMeasurementJudgment(measurementJudgment),
	}
	// 如果有舌面报告则返回
	if report.GetTongueFaceReportStatus() == domain.TongueFaceReportStatusSuccess {
		rsp.TongueFaceReport = &pb.TongueFaceReport{
			Success:        true,
			FaceImageUrl:   report.GetFaceImageURL(),
			TongueImageUrl: report.GetTongueImageURL(),
			Data:           report.GetTongueFaceReport(),
		}
	}

	return nil
}

// 判断测量模块
func judgeMeasurementModule(m *domain.MeasurementJudgmentModule) bool {
	const (
		AbnormalKey = "C0001"
	)
	if m != nil && len(m.MeasurementModule) != 0 {
		// 如果key是C0001，直接返回提示信息其他都不返回
		if m.MeasurementModule[0].Key == AbnormalKey {
			return false
		}
	}
	return true
}

// generateMeasurementJudgment 生成测量判断模块
func generateMeasurementJudgment(measurementJudgment *domain.MeasurementJudgmentModule) *pb.MeasurementJudgmentModule {
	if measurementJudgment == nil || len(measurementJudgment.MeasurementModule) == 0 {
		return &pb.MeasurementJudgmentModule{}
	}
	lookUps := make([]*pb.Lookup, len(measurementJudgment.MeasurementModule))
	for k, v := range measurementJudgment.MeasurementModule {
		lookUps[k] = toProtoLoopUp(v)
	}
	return &pb.MeasurementJudgmentModule{
		Lookups: lookUps,
	}
}

// 验证request
func validateGetReportRequest(req *pb.GetReportRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetReportId() == "" {
		return gerr.New("report_id should not be empty")
	}
	return nil
}

// getRequestModules 获取请求模块
func getRequestModules() []*platformpb.ReportModuleInput {
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
			Name: "heart",
		},
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
			Name: "measurement_judgment",
		},
		{
			Name: "stress_state_judgment",
		},
		{
			Name: "health_index",
		},
		{
			Name: "drunk_judgment",
		},
	}
}

// parseOutputQuestions 解析问题
func parseOutputQuestions(askQuestions map[string]*platformpb.QuestionList) map[string]*pb.QuestionList {
	questionsList := make(map[string]*pb.QuestionList)
	// 解析问题
	for module, v := range askQuestions {
		questions := make([]*pb.Question, len(v.GetQuestions()))
		// 构建问题
		for k, qs := range v.GetQuestions() {
			// 构建选项
			choices := make([]*pb.QuestionChoice, len(qs.GetChoices()))
			for index, choice := range qs.GetChoices() {
				choices[index] = &pb.QuestionChoice{
					// 选项的 key
					ChoiceKey: choice.GetChoiceKey(),
					// 内容
					Content: choice.GetContent(),
					// 冲突选项
					ConflictKeys: choice.GetConflictKeys(),
					// 是否默认选中
					Selected: choice.GetSelected(),
				}
			}
			questions[k] = &pb.Question{
				// 问题的 key
				QuestionKey: qs.GetQuestionKey(),
				// 内容
				Content: qs.GetContent(),
				// 类型
				Type: qs.GetType(),
				// 选项
				Choices: choices,
			}
		}
		questionsList[module] = &pb.QuestionList{
			Questions: questions,
		}
	}
	return questionsList
}

// generateHighAndMidRiskEstimate 生成中高风险报告
func generateHighAndMidRiskEstimate(disease *domain.DiseaseEstimate) ([]*productpb.ReportResult, []*productpb.ReportResult) {
	var highRisks, midRisks []*productpb.ReportResult
	for _, v := range disease.Diseases {
		switch v.Score {
		case 3:
			highRisks = append(highRisks, &productpb.ReportResult{
				SymptomKey: v.Key,
			})
		case 2:
			midRisks = append(midRisks, &productpb.ReportResult{
				SymptomKey: v.Key,
			})
		}
	}
	return highRisks, midRisks
}

// generateRiskEstimate 生成风险报告
func generateRiskEstimate(disease *domain.DiseaseEstimate) (map[string]*pb.RiskEstimate, []string) {
	res := make(map[string]*pb.RiskEstimate)
	highRiskCount, midRiskCount := 0, 0
	var highRisks, midRisks []string

	for _, v := range disease.Diseases {
		// 通过key获取名称
		label := summary.RiskMap[v.Key].Label
		risk := "正常"
		switch v.Score {
		case 3:
			risk = "高风险"
			highRiskCount++
			highRisks = append(highRisks, label)
		case 2:
			risk = "中风险"
			midRiskCount++
			midRisks = append(midRisks, label)
		case 1:
			risk = "轻风险"
		}
		res[label] = &pb.RiskEstimate{
			// 风险名称
			RiskName: label,
			// 风险预估预估
			RiskEstimate: risk,
			// key
			Key: v.Key,
		}
	}

	// 生成风险小结
	riskSummary := []string{}
	highRiskSummary := `<div class="summary">%s<span class="risk-high">高风险</span></div>`
	midRiskSummary := `<div class="summary">%s<span class="risk-medium">中风险</span></div>`
	if highRiskCount > 0 {
		str := strings.Join(highRisks, ",")
		riskSummary = append(riskSummary, fmt.Sprintf(highRiskSummary, str))
	}
	if midRiskCount > 0 {
		str := strings.Join(midRisks, ",")
		riskSummary = append(riskSummary, fmt.Sprintf(midRiskSummary, str))
	}
	return res, riskSummary
}

// f2 湿气 f3 血瘀
func generatePhysicalTherapyReportResult(f2, f3 int32) []*productpb.ReportResult {
	const (
		// 血瘀Key
		xy = "XY"
		// 湿气Key
		sq = "SQ"
	)
	p := make([]*productpb.ReportResult, 2)
	p[0] = &productpb.ReportResult{
		// 症候key
		SymptomKey: sq,
		// 风险等级
		Score: f2,
	}
	p[0] = &productpb.ReportResult{
		// 症候key
		SymptomKey: xy,
		// 风险等级
		Score: f3,
	}
	return p
}

// generateDirtyDialectic 生成脏腑辩证
// 最多显示6个
func generateDirtyDialectic(dirtyDialect *domain.DirtyDialectic) *pb.DirtyDialecticModule {
	lookUps := make([]*pb.Lookup, len(dirtyDialect.LookUps))
	for k, v := range dirtyDialect.LookUps {
		lookUps[k] = toProtoLookUp(v)
	}
	// 最多显示6个
	if len(lookUps) > 6 {
		lookUps = lookUps[:6]
	}
	return &pb.DirtyDialecticModule{
		Lookups: lookUps,
	}
}

// generatePhysiqueDialectic 生成体质辩证
func generatePhysiqueDialectic(physique *domain.PhysiqueDialectics) *pb.PhysiqueDialecticsModule {
	lookUps := make([]*pb.Lookup, len(physique.LookUps))
	for k, v := range physique.LookUps {
		lookUps[k] = toPhysiqueProtoLookUp(v)
	}
	// 最多显示5个
	if len(lookUps) > 5 {
		lookUps = lookUps[:5]
	}
	return &pb.PhysiqueDialecticsModule{
		Lookups: lookUps,
	}
}

// generatePhysique 生成体质报告
func generatePhysique(physique *domain.PhysiqueDialectics) *productpb.ReportResult {
	if physique == nil || len(physique.LookUps) == 0 {
		return nil
	}
	var ph *productpb.ReportResult
	ph = &productpb.ReportResult{
		SymptomKey: physique.LookUps[0].Key,
	}
	return ph
}

// generateRecommendation
func generateRecommendation(recommend *domain.Recommendations, hasTreatment bool) (*pb.SportsAdviceModule, *pb.DietaryAdviceModule, *pb.ChineseMedicineAdviceModule, *pb.PhysicalTherapyAdviceModule) {
	// 运功
	index := 0
	sportsLookups := make([]*pb.Lookup, len(recommend.SportsAdvice.LookUps)+len(recommend.SportsAdviceTips.LookUps))
	for _, v := range recommend.SportsAdvice.LookUps {
		sportsLookups[index] = toProtoLookUp(v)
		index++
	}
	for _, v := range recommend.SportsAdviceTips.LookUps {
		sportsLookups[index] = toProtoLookUp(v)
		index++
	}
	index = 0
	// 饮食
	dietaryLookups := make([]*pb.Lookup, len(recommend.DietaryAdvice.LookUps)+len(recommend.DietaryAdviceTips.LookUps))
	for _, v := range recommend.DietaryAdvice.LookUps {
		dietaryLookups[index] = toProtoLookUp(v)
		index++
	}
	for _, v := range recommend.DietaryAdviceTips.LookUps {
		dietaryLookups[index] = toProtoLookUp(v)
		index++
	}
	// 如果有方案则不需要返回中医调理和理疗调理
	// TODO: 当前版本不需要判断后续如果需要取消注释即可
	// if hasTreatment {
	// 	return &pb.SportsAdviceModule{
	// 			Lookups: sportsLookups,
	// 		}, &pb.DietaryAdviceModule{
	// 			Lookups: dietaryLookups,
	// 		}, nil, nil
	// }
	// 中医调理
	chineseMedicineLookups := make([]*pb.Lookup, len(recommend.ChineseMedicineAdvice.LookUps))
	for i, v := range recommend.ChineseMedicineAdvice.LookUps {
		chineseMedicineLookups[i] = toProtoLookUp(v)
	}
	// 理疗调理
	physicalAdviceLookups := make([]*pb.Lookup, len(recommend.PhysicalTherapyAdvice.LookUps))
	for i, v := range recommend.PhysicalTherapyAdvice.LookUps {
		physicalAdviceLookups[i] = toProtoLookUp(v)
	}
	// 返回数据
	return &pb.SportsAdviceModule{
			Lookups: sportsLookups,
		}, &pb.DietaryAdviceModule{
			Lookups: dietaryLookups,
		}, &pb.ChineseMedicineAdviceModule{
			Lookups: chineseMedicineLookups,
		}, &pb.PhysicalTherapyAdviceModule{
			Lookups: physicalAdviceLookups,
		}
}

// getWesternMedicineReport 获取西医报告
func (u *ReportAPIHandler) getWesternMedicineReport(ctx context.Context, isStressState bool, tenantId string, report domain.ReportIntf, heartModule *platformpb.HeartModuleOutput, diseases *domain.DiseaseEstimate) (*pb.WesternMedicineReport, error) {
	const (
		ExceptWaveCount = 500
	)
	// 血氧
	hasSpo := true
	if heartModule.GetSpo().GetValue() == 0 {
		hasSpo = false
	}

	waveDate, err := u.getPlatformWaveData(ctx, report.GetReportID(), ExceptWaveCount)
	if err != nil {
		return nil, err
	}

	wesReport := &pb.WesternMedicineReport{
		// 心率(单位bpm)
		HeartRate: heartModule.GetAverageHeartRate(),
		// 是否有血氧
		HasSpo: hasSpo,
		// 血氧值
		Spo: report.GetSpo(),
		// 波形图数据
		WaveData: waveDate,
		// 是否是应急态，如果是，只显示血氧，心率和脉搏波数据
		IsStressState: isStressState,
	}
	if !isStressState {
		riskEstimates, riskSummary := generateRiskEstimate(diseases)
		wesReport.RiskEstimate = riskEstimates
		wesReport.WesternMedicineSummaries = riskSummary
	}
	return wesReport, nil
}

// getTCMReport 获取中医报告
func (u *ReportAPIHandler) getTCMReport(isStressState bool, chart *domain.MeridianBarChart, dirtyDialect *domain.DirtyDialectic, physique *domain.PhysiqueDialectics, stressState *domain.StressState) *pb.TCMReport {
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
		physiqueDialects := generatePhysiqueDialectic(physique)
		tcmReport.DirtyDialecticModule = dialect
		tcmReport.PhysiqueDialecticsModule = physiqueDialects
		// 中医报告小结(富文本字段)
		tcmReport.TcmSummaries = u.generateTcmSummary(dialect.Lookups)
	} else {
		tcmReport.StressStateModule = u.genderStressStateModule(stressState)
	}

	return tcmReport
}

// getPhysicalTherapyReport 获取理疗报告
func (u *ReportAPIHandler) getPhysicalTherapyReport(factor *domain.FactorInterpretation) *pb.PhysicalTherapyReport {
	return &pb.PhysicalTherapyReport{
		// 阴虚指数
		F0: factor.F0,
		// 阳虚指数
		F1: factor.F1,
		// 湿气指数
		F2: factor.F2,
		// 血瘀指数
		F3: factor.F3,
		// 阴阳平衡指数
		F4: factor.F4,
	}
}

// getRecommendation 获取商品推荐
func (u *ReportAPIHandler) getRecommendation(isStressState bool, recommend *domain.Recommendations, listProductsRsp *productpb.ListRecommendedProductsResponse, hasTreatment bool) *pb.Recommendation {
	// 如果是应激态，该模块删除
	if isStressState {
		return &pb.Recommendation{
			// 是否是应急态，如果是，只显示柱状图和应急态状态
			IsStressState: true,
		}
	}
	sports, dietary, chineseAdvice, physicalAdvice := generateRecommendation(recommend, hasTreatment)
	// 构造report返回的商品
	products := make([]*pb.Product, len(listProductsRsp.GetProducts()))
	for k, v := range listProductsRsp.GetProducts() {
		products[k] = toProtoProduct(v)
	}
	return &pb.Recommendation{
		// 推荐商品
		Products: products,
		// 食疗建议
		DietaryAdvice: dietary,
		// 运动建议
		SportsAdvice: sports,
		// 是否是应急态，如果是，只显示柱状图和应急态状态
		IsStressState: false,
		// 是否是最新的商品
		AreLatestTreatments: listProductsRsp.GetAreLatestTreatments(),
		// 中药调理
		ChineseMedicineAdvice: chineseAdvice,
		// 理疗建议
		PhysicalTherapyAdvice: physicalAdvice,
	}
}

// getPlatformWaveData 获取脉搏波(平台)
func (u *ReportAPIHandler) getPlatformWaveData(ctx context.Context, reportId string, exceptCount int32) ([]uint32, error) {
	// 获取脉搏波数据
	// 构建OutgoingContext
	ctx = buildPlatformContext(ctx)
	getRawResult, err := u.platformClient.GetRawData(ctx, &platformpb.GetRawDataRequest{
		// App ID
		// 必填
		AppId: u.appId,
		// 报告 ID
		// 必填
		ReportId: reportId,
		// 期望获得的数据长度
		// expected_count>=0 时，至多返回指定数量的数据点数，数量不超过1000；
		// expected_count<0 时，返回全部数据点，仅对高级用户开放.
		// 必填
		ExpectedCount: exceptCount,
	})
	if err != nil {
		return nil, err
	}
	return getRawResult.GetRawData(), nil
}

// generateTcmSummary 生成中医小结
func (u *ReportAPIHandler) generateTcmSummary(lookups []*pb.Lookup) []string {
	res := make([]string, len(lookups))
	summaryTemplate := `
        <div class="summary">
            <div class="summary-title">%s</div>
            <div class="summary-content">
                <p>%s</p >
            </div>
        </div>
	`
	// 新报告
	for k, v := range lookups {
		res[k] = fmt.Sprintf(summaryTemplate, summary.DirtyDialectSummary[v.GetKey()].Label, summary.DirtyDialectSummary[v.GetKey()].Summary)
	}
	return res
}

// genderStressStateModule 生成应急态模块
func (u *ReportAPIHandler) genderStressStateModule(s *domain.StressState) *pb.StressStateModule {
	if !s.HasStressState {
		return nil
	}
	lookups := []*pb.Lookup{}
	if s.HasDoneSports {
		lookups = append(lookups, &pb.Lookup{
			// 序号
			Key: "done_sports",
			// 显示内容
			Content: "运动",
		})
	}
	if s.HasDrinkedWine {
		lookups = append(lookups, &pb.Lookup{
			// 序号
			Key: "has_drinked_wine",
			// 显示内容
			Content: "饮酒",
		})
	}
	if s.HasHadCold {
		lookups = append(lookups, &pb.Lookup{
			// 序号
			Key: "has_had_cold",
			// 显示内容
			Content: "感冒",
		})
	}
	if s.HasRhinitisEpisode {
		lookups = append(lookups, &pb.Lookup{
			// 序号
			Key: "has_rhinitis_episode",
			// 显示内容
			Content: "鼻炎",
		})
	}
	if s.HasAbdominalPain {
		lookups = append(lookups, &pb.Lookup{
			// 序号
			Key: "has_abdominal_pain",
			// 显示内容
			Content: "腹痛",
		})
	}
	if s.HasViralInfection {
		lookups = append(lookups, &pb.Lookup{
			// 序号
			Key: "has_viral_infection",
			// 显示内容
			Content: "病毒",
		})
	}
	if s.HasPhysiologicalPeriod {
		lookups = append(lookups, &pb.Lookup{
			// 序号
			Key: "has_physiological_period",
			// 显示内容
			Content: "生理期",
		})
	}
	if s.HasOvulation {
		lookups = append(lookups, &pb.Lookup{
			// 序号
			Key: "has_ovulation",
			// 显示内容
			Content: "排卵期",
		})
	}
	if s.HasPregnant {
		lookups = append(lookups, &pb.Lookup{
			// 序号
			Key: "has_pregnant",
			// 显示内容
			Content: "怀孕",
		})
	}
	if s.HasLactation {
		lookups = append(lookups, &pb.Lookup{
			// 序号
			Key: "has_lactation",
			// 显示内容
			Content: "哺乳期",
		})
	}
	return &pb.StressStateModule{
		Lookups: lookups,
	}
}

// getAnswers 建立模块名到回答到对应关系
func getAnswers(answers map[string]*pb.AnswerList) map[string]*platformpb.AnswerList {
	questionAnswers := make(map[string]*platformpb.AnswerList)
	for module, qa := range answers {
		als := new(platformpb.AnswerList)
		anws := make([]*platformpb.Answer, len(qa.Answers))
		for k, answer := range qa.Answers {
			anws[k] = &platformpb.Answer{
				// 问题
				QuestionKey: answer.QuestionKey,
				// 答案
				AnswerKeys: answer.AnswerKeys,
			}
		}
		als.Answers = anws
		questionAnswers[module] = als
	}
	return questionAnswers
}

// 获取脏腑辨证的key的falg
// 比如包含Z0003和Z0004
// 返回1<<(3-1)+1<<(4-1)
func getDDKeyFlag(keys []string) (int32, error) {
	res := 0
	for _, v := range keys {
		// 移除字母Z，保留数字部分
		numberStr := strings.TrimLeft(v, "Z")
		// 将字符串转换为整数
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return 0, err
		} else {
			res += 1 << (number - 1)
		}
	}
	return int32(res), nil
}
