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
	// ErrProductNotFound
	ErrProductNotFound = 5029
	// ErrProductNotBelongToOrganization
	ErrProductNotBelongToOrganization = 5030
)

func (s *ProductAPIHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest, rsp *pb.DeleteProductResponse) error {
	// 验证request
	err := validateDeleteProductRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商品版本号
	dProduct, err := s.productStore.GetProduct(ctx, req.GetProductId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if dProduct == nil {
		return errors.Errorf(ErrProductNotFound, "product[%s] not found", req.GetProductId())
	}
	rev := dProduct.GetRev()

	// 不属于该组织商品无法删除
	if dProduct.GetOrganizationID() != req.GetOrganizationId() {
		return errors.Errorf(ErrProductNotBelongToOrganization, "product[%s] not belong organization[%s]", req.GetProductId(), req.GetOrganizationId())
	}
	// 如果商品使用中无法删除
	if dProduct.GetProductStatus() != domain.ProductStatusUnused {
		return errors.Errorf(codes.InvalidRequest, "product[%s] status is not unused", req.GetProductId())
	}

	// 删除单个商品
	err = s.productStore.DeleteProduct(ctx, req.GetProductId(), rev)
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "delete Product failed[%s]", err.Error())
	}
	return nil
}

// 验证request
func validateDeleteProductRequest(req *pb.DeleteProductRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetProductId() == "" {
		return gerr.New("product id should not be empty")
	}
	return nil
}
