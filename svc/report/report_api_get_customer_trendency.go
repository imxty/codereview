package report

import (
	"context"
	"encoding/json"
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	ptime "github.com/jinmukeji/huimaibao-service/pkg/time"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 添加常客报告请求
func (s *ReportAPIHandler) GetCustomerTrendency(ctx context.Context, req *pb.GetCustomerTrendencyRequest, rsp *pb.GetCustomerTrendencyResponse) error {

	// 1.验证request
	err := validateGetCustomerTrendencyRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	endTime := time.Now().UTC()
	startTime := ptime.DayBeginUTCTime(endTime, ptime.LocBeijing).AddDate(0, 0, -int(req.GetDays())).UTC()
	stressCount := make(map[string]float64)
	var cinfos []*pb.CInfo
	// 查询报告
	reports, err := s.reportStore.ListCustomerReports(ctx, req.GetTenantId(), req.GetCustomerId(), startTime, endTime)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	for _, v := range reports {
		// 如果是应激态
		if v.GetIsStressState() {
			// 统计
			var stressState *domain.StressState
			err = json.Unmarshal([]byte(v.GetStressState()), &stressState)
			if err != nil {
				return errors.Errorf(codes.InvalidOperation, "failed to marshal stress state[%s]", err.Error())
			}
			// 如果是应激态统计次数
			// 是否运动过
			if stressState.HasDoneSports {
				stressCount["has_done_sports"]++
			}
			// 是否饮酒
			if stressState.HasDrinkedWine {
				stressCount["has_drunk_wine"]++
			}
			// 是否感冒
			if stressState.HasHadCold {
				stressCount["has_had_cold"]++
			}
			// 是否鼻炎发作
			if stressState.HasRhinitisEpisode {
				stressCount["has_rhinitis_episode"]++
			}
			// 是否腹痛腹泻
			if stressState.HasAbdominalPain {
				stressCount["has_abdominal_pain"]++
			}
			// 是否既往病毒感染
			if stressState.HasViralInfection {
				stressCount["has_viral_infection"]++
			}
			// 是否处于生理期
			if stressState.HasPhysiologicalPeriod {
				stressCount["has_physiological_period"]++
			}
			// 是否处于排卵期
			if stressState.HasOvulation {
				stressCount["has_ovulation"]++
			}
			// 是否怀孕
			if stressState.HasPregnant {
				stressCount["has_pregnant"]++
			}
			// 是否处于哺乳期
			if stressState.HasLactation {
				stressCount["has_lactation"]++
			}
			continue
		}
		// 其他情况则统计经络数据
		cinfos = append(cinfos, &pb.CInfo{
			// 心包经
			C0: v.GetC0(),
			// 肝经
			C1: v.GetC1(),
			// 肾经
			C2: v.GetC2(),
			// 脾经
			C3: v.GetC3(),
			// 肺经
			C4: v.GetC4(),
			// 胃经
			C5: v.GetC5(),
			// 胆经
			C6: v.GetC6(),
			// 膀胱经
			C7: v.GetC7(),
			// 创建时间
			CreateTime: timestamppb.New(v.GetCreatedAt()),
		})
	}

	// 构建应激态统计
	lookups := make([]*pb.Lookup, len(stressCount))
	index := 0
	for k, v := range stressCount {
		lookups[index] = &pb.Lookup{
			Key:     k,
			Content: StressKeyMap[k],
			Score:   v,
		}
		index++
	}

	rsp.CInfos = cinfos
	rsp.StressStateCount = &pb.StressStateCount{
		Lookups: lookups,
	}

	return nil
}

// 验证request
func validateGetCustomerTrendencyRequest(req *pb.GetCustomerTrendencyRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomerId() == "" {
		return gerr.New("customer id should not be empty")
	}
	return nil
}
