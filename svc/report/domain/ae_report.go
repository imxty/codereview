package domain

import (
	platformpb "github.com/jinmukeji/platform-proto/gen/go/platform/report/v1"
)

const (
	// ae分析失败key
	AEAnalysisFailed = "C0001"
)

// StressState 应激态映射
type StressState struct {
	HasStressState         bool `json:"has_stress_state" mapstructure:"has_stress_state"`
	HasDoneSports          bool `json:"has_done_sports" mapstructure:"has_done_sports"`
	HasDrinkedWine         bool `json:"has_drunk_wine" mapstructure:"has_drunk_wine"`
	HasHadCold             bool `json:"has_had_cold" mapstructure:"has_had_cold"`
	HasRhinitisEpisode     bool `json:"has_rhinitis_episode" mapstructure:"has_rhinitis_episode"`
	HasAbdominalPain       bool `json:"has_abdominal_pain" mapstructure:"has_abdominal_pain"`
	HasViralInfection      bool `json:"has_viral_infection" mapstructure:"has_viral_infection"`
	HasPhysiologicalPeriod bool `json:"has_physiological_period" mapstructure:"has_physiological_period"`
	HasOvulation           bool `json:"has_ovulation" mapstructure:"has_ovulation"`
	HasPregnant            bool `json:"has_pregnant" mapstructure:"has_pregnant"`
	HasLactation           bool `json:"has_lactation" mapstructure:"has_lactation"`
	Hints                  []*LookUp
}

// MeasurementJudgmentModule 测量异常判断
type MeasurementJudgmentModule struct {
	// 测量异常模块
	MeasurementModule []*LookUp `json:"measurement_module"`
}

// DiseaseEstimate 疾病预估
type DiseaseEstimate struct {
	// 疾病预估
	Diseases []*LookUp `json:"diseases" mapstructure:"diseases"`
}

// Diseases
type Diseases struct {
	Diseases []*LookUp `json:"disease_estimate"  mapstructure:"disease_estimate"`
}

// FactorInterpretation 阴阳
type FactorInterpretation struct {
	// 阴
	F0 int32 `json:"f0" mapstructure:"f0"`
	// 阳
	F1 int32 `json:"f1" mapstructure:"f1"`
	// 湿气
	F2 int32 `json:"f2" mapstructure:"f2"`
	// 血瘀
	F3 int32 `json:"f3" mapstructure:"f3"`
	// 阴阳平衡
	F4 int32 `json:"f4" mapstructure:"f4"`
}

// 柱状图模块
type MeridianBarChart struct {
	// 心包经
	C0 int32
	// 肝经
	C1 int32
	// 肾经
	C2 int32
	// 脾经
	C3 int32
	// 肺经
	C4 int32
	// 胃经
	C5 int32
	// 胆经
	C6 int32
	// 膀胱经
	C7 int32
}

// PhysiqueDialectics 体质辩证
type PhysiqueDialectics struct {
	LookUps []*LookUp `json:"lookups"`
}

// Recommendations 调理建议
type Recommendations struct {
	// 饮食建议
	DietaryAdvice LookUps `json:"dietary_advice" mapstructure:"dietary_advice"`
	// 运动建议
	SportsAdvice LookUps `json:"sports_advice" mapstructure:"sports_advice"`
	// 饮食建议
	DietaryAdviceTips HealthTips `json:"dietary_advice_tips" mapstructure:"dietary_advice_tips"`
	// 运动建议
	SportsAdviceTips HealthTips `json:"sports_advice_tips" mapstructure:"sports_advice_tips"`
	// 中药调理
	ChineseMedicineAdvice LookUps `json:"chinese_medicine_advice" mapstructure:"chinese_medicine_advice"`
	// 理疗建议
	PhysicalTherapyAdvice LookUps `json:"physical_therapy_advice" mapstructure:"physical_therapy_advice"`
}

// HealthTips LookUp数组
type HealthTips struct {
	LookUps []*LookUp `json:"health_tips" mapstructure:"health_tips"`
}

// DirtyDialectic 脏腑
type DirtyDialectic struct {
	LookUps []*LookUp `json:"lookups"`
}

// LookUps LookUp数组
type LookUps struct {
	LookUps []*LookUp `json:"lookups" mapstructure:"lookups"`
}

// LookUp
type LookUp struct {
	Content string `json:"content" mapstructure:"content"`
	Key     string `json:"key" mapstructure:"key"`
	Score   int    `json:"score" mapstructure:"score"`
}

// GetStressState 获取应急态
func GetStressState(data []*platformpb.ReportModuleResult) (*StressState, error) {
	for _, v := range data {
		if v.ModuleName == "stress_state_judgment" {
			stressModule := &platformpb.StressStateJudgmentModuleOutput{}
			if v.Result == nil {
				return &StressState{}, nil
			}
			// any反序列化
			err := v.Result.UnmarshalTo(stressModule)
			if err != nil {
				return nil, err
			}
			return &StressState{
				HasStressState:         stressModule.GetHasStressState(),
				HasDoneSports:          stressModule.GetHasDoneSports(),
				HasDrinkedWine:         stressModule.GetHasDrunkWine(),
				HasHadCold:             stressModule.GetHasHadCold(),
				HasRhinitisEpisode:     stressModule.GetHasRhinitisEpisode(),
				HasAbdominalPain:       stressModule.GetHasAbdominalPain(),
				HasViralInfection:      stressModule.GetHasViralInfection(),
				HasPhysiologicalPeriod: stressModule.GetHasPhysiologicalPeriod(),
				HasOvulation:           stressModule.GetHasOvulation(),
				HasPregnant:            stressModule.GetHasPregnant(),
				HasLactation:           stressModule.GetHasLactation(),
				Hints:                  mapPlatformLookupsToDomain(stressModule.GetLookups()),
			}, nil
		}
	}
	return nil, nil
}

// GetHeartModule 获取心率血氧
func GetHeartModule(data []*platformpb.ReportModuleResult) (*platformpb.HeartModuleOutput, error) {
	for _, v := range data {
		if v.ModuleName == "heart" {
			heartModule := &platformpb.HeartModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(heartModule)
			if err != nil {
				return nil, err
			}
			return heartModule, nil
		}
	}
	return nil, nil
}

// GetMeridianBarChartModule 获取柱状图模块
func GetMeridianBarChartModule(data []*platformpb.ReportModuleResult) (*MeridianBarChart, error) {
	for _, v := range data {
		if v.ModuleName == "meridian_bar_chart" {
			if v.Result == nil {
				return &MeridianBarChart{}, nil
			}
			meridianBarChartModule := &platformpb.MeridianBarChartModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(meridianBarChartModule)
			if err != nil {
				return nil, err
			}
			// 柱状图数据
			cInfo := meridianBarChartModule.GetMeridianValue()
			return &MeridianBarChart{
				// 心包经
				C0: cInfo.GetC0(),
				// 肝经
				C1: cInfo.GetC1(),
				// 肾经
				C2: cInfo.GetC2(),
				// 脾经
				C3: cInfo.GetC3(),
				// 肺经
				C4: cInfo.GetC4(),
				// 胃经
				C5: cInfo.GetC5(),
				// 胆经
				C6: cInfo.GetC6(),
				// 膀胱经
				C7: cInfo.GetC7(),
			}, nil
		}
	}
	return nil, nil
}

// GetHealthIndex 获取健康指数
func GetHealthIndex(data []*platformpb.ReportModuleResult) (int32, error) {
	for _, v := range data {
		if v.ModuleName == "health_index" {
			if v.Result == nil {
				return 0, nil
			}
			healthIndexModule := &platformpb.HealthIndexModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(healthIndexModule)
			if err != nil {
				return 0, err
			}
			return healthIndexModule.GetHealthIndex().GetValue(), nil
		}
	}
	return 0, nil
}

// GetDiseaseEstimate 获取16个风险预估
func GetDiseaseEstimate(data []*platformpb.ReportModuleResult, stressState bool) (*DiseaseEstimate, error) {
	if stressState {
		return &DiseaseEstimate{}, nil
	}
	for _, v := range data {
		if v.ModuleName == "risk_estimate" {
			if v.Result == nil {
				return &DiseaseEstimate{}, nil
			}
			riskEstimateModule := &platformpb.RiskEstimateModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(riskEstimateModule)
			if err != nil {
				return nil, err
			}
			return &DiseaseEstimate{
				Diseases: mapPlatformLookupsToDomain(riskEstimateModule.GetLookups()),
			}, nil
		}
	}
	return &DiseaseEstimate{}, nil
}

// GetFactorInterpretation 获取阴阳
func GetFactorInterpretation(data []*platformpb.ReportModuleResult) (*FactorInterpretation, error) {
	for _, v := range data {
		if v.ModuleName == "factor_interpretation" {
			if v.Result == nil {
				return &FactorInterpretation{}, nil
			}
			factorInterpretationModule := &platformpb.FactorInterpretationModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(factorInterpretationModule)
			if err != nil {
				return nil, err
			}
			return &FactorInterpretation{
				// 阴
				F0: int32(factorInterpretationModule.F0),
				// 阳
				F1: int32(factorInterpretationModule.F1),
				// 湿气
				F2: int32(factorInterpretationModule.F2),
				// 血瘀
				F3: int32(factorInterpretationModule.F3),
				// 阴阳平衡
				F4: int32(factorInterpretationModule.F4),
			}, nil
		}
	}
	return &FactorInterpretation{}, nil
}

// GetRecommendations 获取调理建议
func GetRecommendations(data []*platformpb.ReportModuleResult, stressState bool) (*Recommendations, error) {
	if stressState {
		return &Recommendations{}, nil
	}
	var dls LookUps
	var sls LookUps
	var mls HealthTips
	var tls HealthTips
	var cma LookUps
	var pta LookUps
	for _, v := range data {
		// 食疗建议模块
		if v.ModuleName == "dietary_advice" {
			if v.Result == nil {
				return &Recommendations{}, nil
			}
			dietaryAdviceModule := &platformpb.DietaryAdviceModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(dietaryAdviceModule)
			if err != nil {
				return nil, err
			}
			dls.LookUps = mapPlatformLookupsToDomain(dietaryAdviceModule.GetLookups())
			mls.LookUps = mapPlatformLookupsToDomain(dietaryAdviceModule.GetHealthTips())
		}
		// 运动方案模块
		if v.ModuleName == "sports_advice" {
			if v.Result == nil {
				return &Recommendations{}, nil
			}
			sportsAdviceModule := &platformpb.SportsAdviceModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(sportsAdviceModule)
			if err != nil {
				return nil, err
			}
			sls.LookUps = mapPlatformLookupsToDomain(sportsAdviceModule.GetLookups())
			tls.LookUps = mapPlatformLookupsToDomain(sportsAdviceModule.GetHealthTips())
		}
		// 中医调理模块
		if v.ModuleName == "chinese_medicine_advice" {
			if v.Result == nil {
				return &Recommendations{}, nil
			}
			chineseAdvice := &platformpb.ChineseMedicineAdviceModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(chineseAdvice)
			if err != nil {
				return nil, err
			}
			cma.LookUps = mapPlatformLookupsToDomain(chineseAdvice.GetLookups())
		}
		// 理疗调理模块
		if v.ModuleName == "physical_therapy_advice" {
			if v.Result == nil {
				return &Recommendations{}, nil
			}
			physicalAdvice := &platformpb.PhysicalTherapyAdviceModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(physicalAdvice)
			if err != nil {
				return nil, err
			}
			lks := mapPlatformLookupsToDomain(physicalAdvice.GetLookups())
			tips := mapPlatformLookupsToDomain(physicalAdvice.GetHealthTips())
			index := 0
			physicalLookups := make([]*LookUp, len(lks)+len(tips))
			for _, v := range lks {
				physicalLookups[index] = v
				index++
			}
			for _, v := range tips {
				physicalLookups[index] = v
				index++
			}
			pta.LookUps = physicalLookups
		}
	}
	return &Recommendations{
		// 饮食建议
		DietaryAdvice: dls,
		// 运动建议
		SportsAdvice: sls,
		// 饮食建议
		DietaryAdviceTips: mls,
		// 运动建议
		SportsAdviceTips: tls,
		// 中药调理
		ChineseMedicineAdvice: cma,
		// 理疗建议
		PhysicalTherapyAdvice: pta,
	}, nil
}

// GetDirtyDialectic 获取脏腑信息
func GetDirtyDialectic(data []*platformpb.ReportModuleResult, stressState bool) (*DirtyDialectic, error) {
	if stressState {
		return &DirtyDialectic{}, nil
	}
	for _, v := range data {
		if v.ModuleName == "dirty_dialectic" {
			if v.Result == nil {
				return &DirtyDialectic{}, nil
			}
			dirtyDialecticModule := &platformpb.DirtyDialecticModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(dirtyDialecticModule)
			if err != nil {
				return nil, err
			}
			return &DirtyDialectic{
				LookUps: mapPlatformLookupsToDomain(dirtyDialecticModule.GetLookups()),
			}, nil
		}
	}
	return &DirtyDialectic{}, nil
}

// GetMeasurementModule 获取测量异常模块
func GetMeasurementModule(data []*platformpb.ReportModuleResult) (*MeasurementJudgmentModule, error) {
	for _, v := range data {
		if v.ModuleName == "measurement_judgment" {
			if v.Result == nil {
				return &MeasurementJudgmentModule{}, nil
			}
			measurementModule := &platformpb.MeasurementJudgmentModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(measurementModule)
			if err != nil {
				return nil, err
			}
			return &MeasurementJudgmentModule{
				MeasurementModule: mapPlatformLookupsToDomain(measurementModule.GetLookups()),
			}, nil
		}
	}
	return &MeasurementJudgmentModule{}, nil
}

// lookups转化
func mapPlatformLookupsToDomain(lookups []*platformpb.Lookup) []*LookUp {
	lks := make([]*LookUp, len(lookups))
	for k, v := range lookups {
		lks[k] = &LookUp{
			Content: v.GetContent(),
			Key:     v.GetKey(),
			Score:   int(v.GetScore()),
		}
	}
	return lks
}

// GetPhysiqueDialectics 获取体质辩证
func GetPhysiqueDialectics(data []*platformpb.ReportModuleResult, stressState bool) (*PhysiqueDialectics, error) {
	if stressState {
		return &PhysiqueDialectics{}, nil
	}
	for _, v := range data {
		if v.ModuleName == "physical_dialectics" {
			if v.Result == nil {
				return &PhysiqueDialectics{}, nil
			}
			physicalDialecticModule := &platformpb.PhysicalDialecticsModuleOutput{}
			// any反序列化
			err := v.Result.UnmarshalTo(physicalDialecticModule)
			if err != nil {
				return nil, err
			}
			return &PhysiqueDialectics{
				LookUps: mapPlatformLookupsToDomain(physicalDialecticModule.GetLookups()),
			}, nil
		}
	}
	return &PhysiqueDialectics{}, nil
}

// GetWeeklyReport 获取周报
func GetWeeklyReport(data []*platformpb.ReportModuleResult) (*platformpb.WeeklyReport, error) {
	for _, v := range data {
		if v.ModuleName == "weekly_report" {
			if v.Result == nil {
				return &platformpb.WeeklyReport{}, nil
			}
			physicalDialecticModule := &platformpb.WeeklyReport{}
			// any反序列化
			err := v.Result.UnmarshalTo(physicalDialecticModule)
			if err != nil {
				return nil, err
			}
			return physicalDialecticModule, nil
		}
	}
	return &platformpb.WeeklyReport{}, nil
}
