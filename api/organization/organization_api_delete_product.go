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

// DeleteProduct
func (s *OrganizationAPIHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest, rsp *pb.DeleteProductResponse) error {
	err := validateDeleteProductRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 删除商品
	_, err = s.productAPI.DeleteProduct(ctx, &productv1.DeleteProductRequest{
		OrganizationId: req.GetOrganizationId(),
		ProductId:      req.GetProductId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
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
