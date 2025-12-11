package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ProductAPIHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest, rsp *pb.ListProductsResponse) error {
	// 验证request
	err := validateListProductsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	var products []domain.RecommendedProductIntf
	var totalCount int64
	// 是否获取组织下所有商品
	if req.GetGetAll() {
		// 查询所有商品
		products, err = s.productStore.ListAllProducts(ctx, req.GetOrganizationId(), int(req.GetPagination().GetOffset()), int(req.GetPagination().GetSize()))
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		totalCount, err = s.productStore.GetProductsTotalCount(ctx, req.GetOrganizationId())
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	} else {
		status, err := toDomainProductStatus(req.GetProductStatus())
		if err != nil {
			return errors.Errorf(codes.InvalidRequest, "invalid product status[%s]", req.GetProductStatus().String())
		}
		// 分页查询商品
		products, err = s.productStore.ListProducts(ctx, req.GetOrganizationId(), status, int(req.GetPagination().GetOffset()), int(req.GetPagination().GetSize()))
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		// 获取相关状态商品数量
		totalCount, err = s.productStore.GetProductsTotalCountWithCertainStatus(ctx, req.GetOrganizationId(), status)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
	}

	// 统计使用中状态的商品总数
	countInUse, err := s.productStore.GetProductsTotalCountWithCertainStatus(ctx, req.GetOrganizationId(), domain.ProductStatusUsing)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 转换
	pProducts := make([]*pb.Product, len(products))
	for i, v := range products {
		pProducts[i] = toProtoProductFromRecommendedProduct(v)
	}

	// 返回结果
	rsp.Products = pProducts
	rsp.TotalCount = int32(totalCount)
	rsp.ProductUsingCount = int32(countInUse)
	return nil
}

// 验证request
func validateListProductsRequest(req *pb.ListProductsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be empty")
	}
	if !req.GetGetAll() && req.GetProductStatus() == pb.ProductStatus_PRODUCT_STATUS_UNSET {
		return gerr.New("invalid search condition")
	}
	return nil
}
