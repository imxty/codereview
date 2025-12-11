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

func (s *OrganizationAPIHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest, rsp *pb.GetProductResponse) error {
	err := validateGetProductRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 发送请求
	getRsp, err := s.productAPI.GetProduct(ctx, &productv1.GetProductRequest{
		ProductId: req.GetProductId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 转换
	pProduct := toAppProduct(getRsp.GetProduct(), s.s3Domain)

	rsp.Product = pProduct
	return nil
}

// 验证request
func validateGetProductRequest(req *pb.GetProductRequest) error {
	if req.GetProductId() == "" {
		return gerr.New("product id should not be empty")
	}
	return nil
}
