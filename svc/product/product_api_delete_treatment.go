package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrTreatmentNotBelongToOrganization
	ErrTreatmentNotBelongToOrganization = 5050
)

// 方法步骤
// 1. 删除方案表
// 2. 删除方案版本表
// 3. 删除方案配置表
// 4. 将方案状态改为已删除
// (草稿，已过审未发布)
func (s *ProductAPIHandler) DeleteTreatment(ctx context.Context, req *pb.DeleteTreatmentRequest, rsp *pb.DeleteTreatmentResponse) error {
	// 验证request
	err := validateDeleteTreatmentRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取方案信息
	treatment, err := s.productStore.GetTreatment(ctx, req.GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if treatment == nil {
		return errors.Errorf(ErrGetTreatmentFailed, "get treatment[%s] failed", req.GetTreatmentId())
	}
	if treatment.GetOrganizationID() != req.GetOrganizationId() {
		return errors.Errorf(ErrTreatmentNotBelongToOrganization, "treatment[%s] not belong to organization[%s]", req.GetTreatmentId(), req.GetOrganizationId())
	}

	// 获取treatment对应的最新方案rev
	tRev, err := s.productStore.GetTreatmentRev(ctx, req.GetOrganizationId(), req.GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if tRev == nil {
		return errors.Errorf(ErrGetTreatmentFailed, "organization[%s] treatment rev not found[%s]", req.GetOrganizationId(), treatment.GetLatestRev())
	}

	// 判断状态
	// 草稿或者已过审未发布
	if !(tRev.GetTreatmentStatus() == domain.TreatmentStatusDraft || (tRev.GetTreatmentStatus() == domain.TreatmentStatusApproved && !tRev.GetIsPublished())) {
		return errors.Errorf(codes.InvalidOperation, "treatment[%s] status can not delete", req.GetTreatmentId())
	}

	// 获取老方案
	// 如果是已过审
	// 需要对商品状态进行修改
	// 对于老方案商品，需要判断是否有其他方案在使用
	if tRev.GetTreatmentStatus() == domain.TreatmentStatusApproved {
		// 获取老版本方案
		items, err := s.productStore.ListTreatmentRevItems(ctx, tRev.GetTreatmentRevID())
		if err != nil {
			return errors.Errorf(ErrGetTreatmentFailed, "get treatment items[%s] failed", tRev.GetTreatmentRevID())
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
		oldItems, err := s.productStore.CheckProductUsage(ctx, tRev.GetOrganizationID())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 判断商品是否可以修改为未使用
		canNotModifyProductId := make(map[string]bool)
		for _, v := range oldItems {
			// 不属于老方案且是老方案商品
			if v.GetTreatmentRevID() != tRev.GetTreatmentRevID() && oldProductMap[v.GetRecommendedProductID()] {
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

	rev := treatment.GetRev()
	// 删除
	// 删除方案表
	err = s.productStore.DeleteTreatment(ctx, req.GetTreatmentId(), rev)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

// 验证request
func validateDeleteTreatmentRequest(req *pb.DeleteTreatmentRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment_id should not be empty")
	}
	return nil
}
