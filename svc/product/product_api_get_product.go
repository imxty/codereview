package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ProductAPIHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest, rsp *pb.GetProductResponse) error {
	// 验证request
	err := validateGetProductRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取单个商品
	product, err := s.productStore.GetProduct(ctx, req.GetProductId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if product == nil {
		return errors.Error(ErrProductNotFound, "product from Db is nil")
	}

	// 转换
	rsp.Product = toProtoProductFromRecommendedProduct(product)
	return nil
}

// 验证request
func validateGetProductRequest(req *pb.GetProductRequest) error {
	if req.GetProductId() == "" {
		return gerr.New("product id should not be empty")
	}
	return nil
}
