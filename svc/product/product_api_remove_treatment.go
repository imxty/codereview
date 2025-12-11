package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrTreatmentUsing
	ErrTreatmentUsing = 5048
)

// 下架方案
// 1. 查询方案和方案对应的商品
// 2. 下架方案
// 3. 将未在其他方案中使用的商品改为待使用
func (s *ProductAPIHandler) RemoveTreatment(ctx context.Context, req *pb.RemoveTreatmentRequest, rsp *pb.RemoveTreatmentResponse) error {
	err := validateRemoveTreatmentRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 查询方案和方案对应的商品
	treatmentRev, err := s.productStore.GetTreatmentRev(ctx, req.GetOrganizationId(), req.GetTreatmentId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if treatmentRev == nil {
		return errors.Errorf(ErrGetTreatmentFailed, "get treatment[%s,%s] failed", req.GetOrganizationId(), req.GetTreatmentId())
	}

	// 查询方案是否被使用，如果被使用无法下架
	checkRsp, err := s.userAPI.CheckTenantTreatmentUsing(ctx, &userv1.CheckTenantTreatmentUsingRequest{
		TreatmentId: req.GetTreatmentId(),
	})
	if err != nil {
		return err
	}
	// 如果正在使用报错
	if checkRsp.GetUsing() {
		return errors.Errorf(ErrTreatmentUsing, "treatment[%s] is using", req.GetTreatmentId())
	}

	// 批量获取即将发布的方案的所有商品信息
	products, err := s.productStore.BatchGetRecommendedProductsByTreatmentRevId(ctx, treatmentRev.GetTreatmentRevID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 获取所有商品ID
	productIds := make([]string, len(products))
	for k, v := range products {
		productIds[k] = v.GetRecommendedProductID()
	}

	// 查询在其他已发布的方案中的商品ID
	pids, err := s.productStore.ListTreatmentProductsNotUsingWithoutItself(ctx, productIds, treatmentRev.GetTreatmentRevID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 将productIds-pids的剩余商品ID状态变为待使用
	var notUsingPids []string
	usingPid := make(map[string]bool)
	for _, v := range pids {
		usingPid[v] = true
	}
	for _, v := range productIds {
		if !usingPid[v] {
			notUsingPids = append(notUsingPids, v)
		}
	}

	// 将商品状态变为待使用
	if len(notUsingPids) != 0 {
		err = s.productStore.BatchChangeProductStatus(ctx, notUsingPids, domain.ProductStatusToBeUsed)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}
	// 将方案下架
	err = s.productStore.RemoveTreatmentByTreatmentRevId(ctx, treatmentRev.GetTreatmentRevID(), treatmentRev.GetRev())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	return nil
}

func validateRemoveTreatmentRequest(req *pb.RemoveTreatmentRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetTreatmentId() == "" {
		return gerr.New("treatment id should not be empty")
	}
	return nil
}
