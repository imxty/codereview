package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	productv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *OrganizationAPIHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest, rsp *pb.ListProductsResponse) error {
	err := validateListProductsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	listRsp, err := s.productAPI.ListProducts(ctx, &productv1.ListProductsRequest{
		// 组织id
		OrganizationId: req.GetOrganizationId(),
		// 分页
		Pagination: toProductPagination(req.GetPagination()),
		// 是否获取所有商品
		GetAll: req.GetGetAll(),
		// 请求商品状态(先判断是否获取所有商品)
		ProductStatus: toSvcProductStatus(req.GetProductStatus()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 商品转化
	pProducts := listRsp.GetProducts()
	p := make([]*pb.Product, len(pProducts))
	for k, v := range pProducts {
		p[k] = toAppProduct(v, s.s3Domain)
	}

	rsp.ProductUsingCount = listRsp.GetProductUsingCount()
	rsp.Products = p
	rsp.TotalCount = listRsp.GetTotalCount()
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
