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

func (s *OrganizationAPIHandler) ListSymptomProducts(ctx context.Context, req *pb.ListSymptomProductsRequest, rsp *pb.ListSymptomProductsResponse) error {
	err := validateListSymptomProductsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 发送请求
	listRsp, err := s.productAPI.ListSymptomProducts(ctx, &productv1.ListSymptomProductsRequest{
		OrganizationId: req.GetOrganizationId(),
		SymptomKey:     req.GetSymptomKey(),
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

	rsp.Products = p
	return nil
}

// 验证request
func validateListSymptomProductsRequest(req *pb.ListSymptomProductsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetSymptomKey() == "" {
		return gerr.New("symptom key should not be empty")
	}
	return nil
}
