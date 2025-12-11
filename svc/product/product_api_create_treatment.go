package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/huimaibao-service/svc/summary"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

const (
	// ErrTreatmentExceedLimit
	ErrTreatmentExceedLimit = 5302
)

const (
	MaxTreatmentNumber = 20
)

// CreateTreatment
func (s *ProductAPIHandler) CreateTreatment(ctx context.Context, req *pb.CreateTreatmentRequest, rsp *pb.CreateTreatmentResponse) error {
	// 验证request
	err := validateCreateTreatmentRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 方案类型转换
	tr := req.GetTreatment()
	treatmentId := xid.New().String()
	treatmentRevId := xid.New().String()

	// 查询方案数量
	ts, err := s.productStore.ListLatestTreatmentRevs(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if len(ts) >= MaxTreatmentNumber {
		return errors.Errorf(ErrTreatmentExceedLimit, "organization[%s] treatments exceed limit 20", req.GetOrganizationId())
	}

	dTreatment := &domain.Treatment{
		TreatmentID:    treatmentId,
		OrganizationID: req.GetOrganizationId(),
		LatestRev:      treatmentRevId,
	}

	// 方案版本类型转换
	dTreatmentRev := &domain.TreatmentRev{
		TreatmentRevID:  treatmentRevId,
		TreatmentID:     treatmentId,
		TreatmentName:   tr.GetTreatmentName(),
		OrganizationID:  req.GetOrganizationId(),
		IsPublished:     domain.TreatmentRevIsNotPublished,
		TreatmentStatus: domain.TreatmentStatusDraft,
	}
	// 方案配置类型转换
	dTreatmentItems, err := transferToDomainTreatmentItems(req.GetTreatment(), treatmentRevId)
	if err != nil {
		return err
	}

	ctx = s.productStore.BeginTx(ctx)
	// 创建方案表
	err = s.productStore.CreateTreatment(ctx, dTreatment)
	if err != nil {
		// 事务回滚
		s.productStore.RollbackTx(ctx)
		return errors.Errorf(codes.DataAccessFailed, "create SubTreatment failed [%s]", err.Error())
	}

	// 创建方案版本表
	err = s.productStore.CreateTreatmentRev(ctx, dTreatmentRev)
	if err != nil {
		// 事务回滚
		s.productStore.RollbackTx(ctx)
		return errors.Errorf(codes.DataAccessFailed, "create SubTreatmentRev failed[%s]", err.Error())
	}
	// 创建多个方案配置表
	err = s.productStore.CreateTreatmentItems(ctx, dTreatmentItems)
	if err != nil {
		// 事务回滚
		s.productStore.RollbackTx(ctx)
		return errors.Errorf(codes.DataAccessFailed, "create SubTreatmentItems failed[%s]", err.Error())
	}

	// 通过方案id获取proto类型方案 返回
	protoTreatment, _, err := s.getProtoTreatment(ctx, req.GetOrganizationId(), treatmentId)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	s.productStore.CommitTx(ctx)

	rsp.Treatment = protoTreatment
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

// 方案配置类型转换
func transferToDomainTreatmentItems(tr *pb.Treatment, treatmentRevId string) ([]domain.TreatmentItemIntf, error) {
	// 分别初始化三个方案配置
	tdd := tr.GetTreatmentItemsDirtyDialectic()
	trd := tr.GetTreatmentItemsRiskyDisease()
	tpt := tr.GetTreatmentItemsPhysicalTherapy()
	tpd := tr.GetTreatmentItemsPhysicalDialectics()
	// 设置匿名函数 分别转换
	f := func(treatmentItems []*pb.TreatmentItem, treatmentType int32) ([]domain.TreatmentItemIntf, error) {
		if len(treatmentItems) == 0 {
			return nil, nil
		}
		var res []domain.TreatmentItemIntf
		for _, v := range treatmentItems {
			if len(v.GetProducts()) > 0 {
				items, err := toDomainTreatmentItem(v, treatmentRevId, treatmentType)
				if err != nil {
					return nil, err
				}
				if len(items) > 0 {
					res = append(res, items...)
				}
			}
		}
		if len(res) == 0 {
			return nil, nil
		}
		return res, nil
	}
	tddItems, err := f(tdd, domain.TreatmentItemTypeDirtyDialectic)
	if err != nil {
		return nil, err
	}
	trdItems, err := f(trd, domain.TreatmentItemTypeRiskyDisease)
	if err != nil {
		return nil, err
	}
	tptItems, err := f(tpt, domain.TreatmentItemTypePhysicalTherapy)
	if err != nil {
		return nil, err
	}
	tpdItems, err := f(tpd, domain.TreatmentItemTypePhysical)
	if err != nil {
		return nil, err
	}

	// 结果合并为一个数组
	res := []domain.TreatmentItemIntf{}
	for _, v := range [][]domain.TreatmentItemIntf{tddItems, trdItems, tptItems, tpdItems} {
		if v != nil {
			res = append(res, v...)
		}
	}
	return res, nil
}

// 转换 toDomainTreatmentItem
func toDomainTreatmentItem(items *pb.TreatmentItem, treatmentRevId string, treatmentItemType int32) ([]domain.TreatmentItemIntf, error) {
	if len(items.GetProducts()) == 0 {
		return nil, nil
	}
	var symptomProducts []domain.TreatmentItemIntf
	for _, v := range items.GetProducts() {
		if v.GetProductId() != "" {
			tt := &domain.TreatmentItem{
				TreatmentItemID:      xid.New().String(),
				TreatmentRevID:       treatmentRevId,
				TreatmentItemType:    treatmentItemType,
				Symptom:              items.GetSymptom(),
				RecommendedProductID: v.GetProductId(),
			}
			symptomProducts = append(symptomProducts, tt)
		}
	}
	return symptomProducts, nil
}

func (p *ProductAPIHandler) getProtoTreatment(ctx context.Context, organizationID string, treatmentId string) (*pb.Treatment, int32, error) {

	// 获取方案信息
	dTreatmentRev, err := p.productStore.GetTreatmentRev(ctx, organizationID, treatmentId)
	if err != nil {
		return nil, 0, errors.Error(codes.DataAccessFailed, err.Error())
	}
	if dTreatmentRev == nil {
		return nil, 0, errors.Error(ErrGetTreatmentFailed, "get treatment_rev failed")
	}
	// 获取方案配置
	dAssociations, err := p.productStore.ListTreatmentItems(ctx, treatmentId)
	if err != nil {
		return nil, 0, errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获得proto treatment_item
	itemsRiskyDisease, itemsDirtyDialectic, itemsPhysicalTherapy, itemsPhysical := transferToProtoTreatmentItems(dAssociations)

	// 获得proto treatment
	pbTreatment := transferToProtoTreatment(dTreatmentRev, itemsRiskyDisease, itemsDirtyDialectic, itemsPhysicalTherapy, itemsPhysical)

	// 获得方案不同商品的数量(去重)
	pMap := make(map[string]bool)
	var count int32
	for _, v := range dAssociations {
		pId := v.GetRecommendedProductID()
		if !pMap[pId] {
			count++
		}
		pMap[pId] = true
	}
	return pbTreatment, count, nil
}

func transferToProtoTreatmentItems(dAssociations []domain.TreatmentProductAssociationIntf) (rd, dd, pt, pd []*pb.TreatmentItem) {
	// 症候与商品map
	rdm := make(map[string][]*pb.Product)
	ddm := make(map[string][]*pb.Product)
	ptm := make(map[string][]*pb.Product)
	ptd := make(map[string][]*pb.Product)
	// 构建匿名函数 将associations中的症候商品转换成map[string][]*pb.Product的形式
	tempFunc := func(treatmentItemType int32, res map[string][]*pb.Product) {
		for _, v := range dAssociations {
			if v.GetTreatmentItemType() == treatmentItemType {
				res[v.GetSymptom()] = append(res[v.GetSymptom()], toProtoProductFromTreatmentProductAssociation(v))
			}
		}
	}
	// 分别转换三种类型方案配置
	tempFunc(domain.TreatmentItemTypeRiskyDisease, rdm)
	tempFunc(domain.TreatmentItemTypeDirtyDialectic, ddm)
	tempFunc(domain.TreatmentItemTypePhysicalTherapy, ptm)
	tempFunc(domain.TreatmentItemTypePhysical, ptd)

	// 将症候商品的map转换为proto []*pb.TreatmentItem类型
	convert := func(m map[string][]*pb.Product) []*pb.TreatmentItem {
		res := []*pb.TreatmentItem{}
		for k, v := range m {
			item := &pb.TreatmentItem{
				Symptom:  k,
				Products: v,
			}
			res = append(res, item)
		}
		return res
	}
	return convert(rdm), convert(ddm), convert(ptm), convert(ptd)
}
