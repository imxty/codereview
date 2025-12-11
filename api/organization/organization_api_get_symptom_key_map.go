package organization

import (
	"context"

	physical "github.com/jinmukeji/huimaibao-service/svc/summary"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
)

func (s *OrganizationAPIHandler) GetSymptomKeyMap(ctx context.Context, req *pb.GetSymptomKeyMapRequest, rsp *pb.GetSymptomKeyMapResponse) error {
	r := make(map[string]string)
	for k, v := range physical.RiskyDiseaseMap {
		r[k] = v.Label
	}
	d := make(map[string]string)
	for k, v := range physical.DirtyDialecticsMap {
		d[k] = v.Label
	}
	p := make(map[string]string)
	for k, v := range physical.PhysicalTherapyMap {
		p[k] = v.Label
	}
	pd := make(map[string]string)
	for k, v := range physical.PhysiqueMap {
		// 平和质不需要返回
		if k == "T0012" {
			continue
		}
		pd[k] = v.Label
	}
	rsp.RiskDiseaseKeyMap = r
	rsp.DirtyDialecticsKeyMap = d
	rsp.PhysicalTherapyKeyMap = p
	rsp.PhysiqueDialecticsKeyMap = pd
	return nil
}
