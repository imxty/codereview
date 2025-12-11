package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	physical "github.com/jinmukeji/huimaibao-service/svc/summary"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

func (s *ProductAPIHandler) UpdateTreatment(ctx context.Context, req *pb.UpdateTreatmentRequest, rsp *pb.UpdateTreatmentResponse) error {
	// 验证request
	err := validateUpdateTreatmentRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 根据方案id拿到最新方案版本信息
	latestTreatmentRev, err := s.productStore.GetTreatmentRev(ctx, req.GetTreatment().GetOrganizationId(), req.GetTreatment().GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if latestTreatmentRev == nil {
		return errors.Error(ErrGetTreatmentFailed, "get treatment failed")
	}
	revStatus := latestTreatmentRev.GetTreatmentStatus()
	// 判断当前方案的状态，删除，发布，审核中都无法进行修改
	if latestTreatmentRev.GetIsPublished() || revStatus == domain.TreatmentStatusUnderReview {
		return errors.Errorf(codes.InvalidRequest, "treatment_rev[%s] status can not update", latestTreatmentRev.GetTreatmentRevID())
	}
	tr := req.GetTreatment()
	tdd := tr.GetTreatmentItemsDirtyDialectic()
	trd := tr.GetTreatmentItemsRiskyDisease()
	tpt := tr.GetTreatmentItemsPhysicalTherapy()

	// 检验症候key是否合法 且是否是对应类别
	for _, v := range tdd {
		if _, ok := physical.DirtyDialecticsMap[v.GetSymptom()]; !ok {
			return errors.Error(codes.InvalidRequest, "invalid DirtyDialectics symptom key")
		}
	}
	for _, v := range trd {
		if _, ok := physical.RiskyDiseaseMap[v.GetSymptom()]; !ok {
			return errors.Error(codes.InvalidRequest, "invalid RiskyDisease symptom key")
		}
	}
	for _, v := range tpt {
		if _, ok := physical.PhysicalTherapyMap[v.GetSymptom()]; !ok {
			return errors.Error(codes.InvalidRequest, "invalid PhysicalTherapy symptom key")
		}
	}

	// 创建新的方案版本(更新方案版本id和方案名称)
	newTreatmentRev := &domain.TreatmentRev{
		TreatmentRevID:  xid.New().String(),
		TreatmentID:     latestTreatmentRev.GetTreatmentID(),
		TreatmentName:   req.GetTreatment().GetTreatmentName(),
		Remarks:         latestTreatmentRev.GetRemarks(),
		OrganizationID:  latestTreatmentRev.GetOrganizationID(),
		IsPublished:     false,
		TreatmentStatus: domain.TreatmentStatusDraft,
	}

	// 获取方案版本更新map(创建新的方案版本)
	err = s.productStore.CreateTreatmentRev(ctx, newTreatmentRev)
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "create sub treatment rev failed when updating treatment[%s]", err.Error())
	}
	// 转换方案配置类型
	dTreatmentItems, err := transferToDomainTreatmentItems(req.GetTreatment(), newTreatmentRev.TreatmentRevID)
	if err != nil {
		return err
	}

	// 获取老方案
	// 如果是已过审
	// 需要对商品状态进行修改
	// 对于老方案商品，需要判断是否有其他方案在使用
	if latestTreatmentRev.GetTreatmentStatus() == domain.TreatmentStatusApproved {
		// 获取新版本商品
		// newProductMap := make(map[string]bool)
		// newProductId := []string{}
		// for _, v := range dTreatmentItems {
		// 	if !newProductMap[v.GetRecommendedProductID()] {
		// 		newProductId = append(newProductId, v.GetRecommendedProductID())
		// 		newProductMap[v.GetRecommendedProductID()] = true
		// 	}
		// }
		// 获取老版本方案
		items, err := s.productStore.ListTreatmentRevItems(ctx, latestTreatmentRev.GetTreatmentRevID())
		if err != nil {
			return errors.Errorf(ErrGetTreatmentFailed, "get treatment items[%s] failed", latestTreatmentRev.GetTreatmentRevID())
		}
		// 获取老版本商品
		oldProductMap := make(map[string]bool)
		oldProductId := []string{}
		for _, v := range items {
			if !oldProductMap[v.GetRecommendedProductID()] {
				oldProductId = append(oldProductId, v.GetRecommendedProductID())
				oldProductMap[v.GetRecommendedProductID()] = true
			}
		}
		// 查询老版本商品是否在其他方案中使用
		oldItems, err := s.productStore.CheckProductUsage(ctx, latestTreatmentRev.GetOrganizationID())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 判断商品是否可以修改为
		canNotModifyProductId := make(map[string]bool)
		for _, v := range oldItems {
			// 不属于老方案且是老方案商品
			if v.GetTreatmentRevID() != latestTreatmentRev.GetTreatmentRevID() && oldProductMap[v.GetRecommendedProductID()] {
				// 判断状态,如果过审就无法修改
				if v.GetTreatmentStatus() == domain.TreatmentStatusApproved {
					canNotModifyProductId[v.GetRecommendedProductID()] = true
				}
			}
		}
		// 需要修改的商品ID
		draftProductId := []string{}
		for _, v := range oldProductId {
			if !canNotModifyProductId[v] {
				draftProductId = append(draftProductId, v)
			}
		}
		if len(draftProductId) != 0 {
			err = s.productStore.BatchChangeProductStatus(ctx, draftProductId, domain.ProductStatusUnused)
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
		}
	}

	// 创建新的方案配置
	err = s.productStore.CreateTreatmentItems(ctx, dTreatmentItems)
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "create sub treatment items failed when updating treatment[%s]", err.Error())
	}

	// 获取方案表当前版本号
	treatment, err := s.productStore.GetTreatment(ctx, req.GetTreatment().GetTreatmentId())
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "create treatment failed when updating treatment[%s]", err.Error())
	}
	if treatment == nil {
		return errors.Error(ErrGetTreatmentFailed, "get treatment failed")
	}
	rev := treatment.GetRev()
	// 更新方案表里的最新方案版本id
	err = s.productStore.UpdateLatestTreatmentRev(ctx, req.GetTreatment().GetTreatmentId(), newTreatmentRev.TreatmentRevID, rev)
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "update latest treatment rev failed[%s]", err.Error())
	}

	return nil
}

// 验证request
func validateUpdateTreatmentRequest(req *pb.UpdateTreatmentRequest) error {
	tr := req.GetTreatment()
	if tr == nil {
		return gerr.New("treatment id should not be empty")
	}
	if req.GetTreatment().GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}

	if tr.GetTreatmentId() == "" {
		return gerr.New("treatment_id should not be empty")
	}
	tdd := tr.GetTreatmentItemsDirtyDialectic()
	trd := tr.GetTreatmentItemsRiskyDisease()
	tpt := tr.GetTreatmentItemsPhysicalTherapy()
	// 检验症候key是否合法 且是否是对应类别
	for _, v := range tdd {
		if _, ok := physical.DirtyDialecticsMap[v.GetSymptom()]; !ok {
			return gerr.New("invalid DirtyDialectics symptom key")
		}
	}
	for _, v := range trd {
		if _, ok := physical.RiskyDiseaseMap[v.GetSymptom()]; !ok {
			return gerr.New("invalid RiskyDisease symptom key")
		}
	}
	for _, v := range tpt {
		if _, ok := physical.PhysicalTherapyMap[v.GetSymptom()]; !ok {
			return gerr.New("invalid PhysicalTherapy symptom key")
		}
	}

	if len(tdd) == 0 && len(trd) == 0 && len(tpt) == 0 {
		return gerr.New("treatment_items should not be empty")
	}
	return nil
}
