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
	// ErrPublishTreatmentFailed
	ErrPublishTreatmentFailed = 5031
)

// 发布方案
// 操作步骤
// 1. 判断方案是否过审 审核通过才能发布
// 2. 将目标方案改为已发布 并把改方案相关的所有商品状态改为使用中
func (s *ProductAPIHandler) PublishTreatment(ctx context.Context, req *pb.PublishTreatmentRequest, rsp *pb.PublishTreatmentResponse) error {
	// 验证request
	err := validatePublishTreatmentRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取方案id
	treatmentId := req.GetTreatmentId()

	// 获取方案版本号 判断方案是否已过审
	treatmentRev, err := s.productStore.GetTreatmentRev(ctx, req.GetOrganizationId(), treatmentId)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if treatmentRev == nil {
		return errors.Errorf(ErrGetTreatmentFailed, "get treatment[%s,%s] failed", req.GetOrganizationId(), treatmentId)
	}
	// 方案没过审不能发布
	if treatmentRev.GetTreatmentStatus() != domain.TreatmentStatusApproved {
		return errors.Errorf(ErrPublishTreatmentFailed, "publish treatment[%s] failed for unapproved treatment status", req.GetTreatmentId())
	}
	rev := treatmentRev.GetRev()

	// 批量获取即将发布的方案的所有商品信息,都变为使用中即可
	productsToBePublished, err := s.productStore.BatchGetRecommendedProductsByTreatmentRevId(ctx, treatmentRev.GetTreatmentRevID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	productIds := make([]string, len(productsToBePublished))
	for k, v := range productsToBePublished {
		productIds[k] = v.GetRecommendedProductID()
	}

	// 开启事务
	ctx = s.productStore.BeginTx(ctx)

	// 将目标方案的发布状态改为已发布
	err = s.productStore.PublishTreatment(ctx, treatmentId, rev)
	if err != nil {
		// 回滚
		s.productStore.RollbackTx(ctx)
		return errors.Errorf(ErrPublishTreatmentFailed, "publish failed [%s]", err.Error())
	}

	// 发布方案后 修改发布方案中所有商品的状态为 使用中状态
	err = s.productStore.BatchChangeProductStatus(ctx, productIds, domain.ProductStatusUsing)
	if err != nil {
		// 回滚
		s.productStore.RollbackTx(ctx)
		return errors.Errorf(ErrPublishTreatmentFailed, "change products status failed [%s]", err.Error())
	}

	// 提交事物
	s.productStore.CommitTx(ctx)

	return nil
}

// 验证request
func validatePublishTreatmentRequest(req *pb.PublishTreatmentRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment id should not be empty")
	}
	return nil
}
