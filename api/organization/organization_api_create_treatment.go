package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	productv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/huimaibao-service/svc/summary"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *OrganizationAPIHandler) CreateTreatment(ctx context.Context, req *pb.CreateTreatmentRequest, rsp *pb.CreateTreatmentResponse) error {
	err := validateCreateTreatmentRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	createRsp, err := s.productAPI.CreateTreatment(ctx, &productv1.CreateTreatmentRequest{
		OrganizationId: req.GetOrganizationId(),
		Treatment:      toProductTreatment(req.GetTreatment()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	rsp.Treatment = toAppTreatment(createRsp.GetTreatment(), s.s3Domain)

	return nil
}

// 验证request
func validateCreateTreatmentRequest(req *pb.CreateTreatmentRequest) error {
	tr := req.GetTreatment()
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}

	if tr.GetTreatmentName() == "" {
		return gerr.New("treatment_name should not be empty")
	}

	// 获取三中方案类型配置 合并成一个数组进行检验
	tdd := tr.GetTreatmentItemsDirtyDialectic()
	trd := tr.GetTreatmentItemsRiskyDisease()
	tpt := tr.GetTreatmentItemsPhysicalTherapy()
	tpd := tr.GetTreatmentItemsPhysicalDialectics()
	// 检验症候key是否合法 且是否是对应类别
	for _, v := range tdd {
		if _, ok := summary.DirtyDialecticsMap[v.GetSymptom()]; !ok {
			return gerr.New("invalid DirtyDialectics symptom key")
		}
	}
	for _, v := range trd {
		if _, ok := summary.RiskyDiseaseMap[v.GetSymptom()]; !ok {
			return gerr.New("invalid RiskyDisease symptom key")
		}
	}
	for _, v := range tpt {
		if _, ok := summary.PhysicalTherapyMap[v.GetSymptom()]; !ok {
			return gerr.New("invalid PhysicalTherapy symptom key")
		}
	}
	for _, v := range tpd {
		if _, ok := summary.PhysiqueMap[v.GetSymptom()]; !ok {
			return gerr.New("invalid Physical symptom key")
		}
	}

	if len(tdd) == 0 && len(trd) == 0 && len(tpt) == 0 && len(tpd) == 0 {
		return gerr.New("treatment_items should not be empty")
	}
	return nil
}
