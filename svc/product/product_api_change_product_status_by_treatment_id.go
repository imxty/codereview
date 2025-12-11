package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ProductAPIHandler) ChangeProductStatusByTreatmentRevId(ctx context.Context, req *pb.ChangeProductStatusByTreatmentRevIdRequest, rsp *pb.ChangeProductStatusByTreatmentRevIdResponse) error {
	// 验证request
	err := validateChangeProductStatusByTreatmentRevIdRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商品信息
	products, err := s.productStore.BatchGetRecommendedProductsByTreatmentRevId(ctx, req.GetTreatmentRevId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 转换商品状态
	status, _ := toDomainProductStatus(req.GetProductStatus())
	for _, v := range products {
		// 使用中无法修改
		if v.GetProductStatus() == domain.ProductStatusUsing {
			continue
		}
		err = s.productStore.ChangeProductStatus(ctx, v.GetRecommendedProductID(), v.GetRev(), status)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	return nil
}

// 验证request
func validateChangeProductStatusByTreatmentRevIdRequest(req *pb.ChangeProductStatusByTreatmentRevIdRequest) error {
	if req.GetTreatmentRevId() == "" {
		return gerr.New("treatment_rev_id should not be empty")
	}
	if _, err := toDomainProductStatus(req.GetProductStatus()); err != nil {
		return gerr.New("invalid product_status")
	}
	return nil
}
