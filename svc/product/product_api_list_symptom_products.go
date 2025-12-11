package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取症候商品请求
func (s *ProductAPIHandler) ListSymptomProducts(ctx context.Context, req *pb.ListSymptomProductsRequest, rsp *pb.ListSymptomProductsResponse) error {
	err := validateListSymptomProductsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 获取组织下所有商品
	products, err := s.productStore.ListOrganizationProducts(ctx, req.GetOrganizationId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	resultProduct := getSymptomProducts(products, req.GetSymptomKey())
	rsp.Products = resultProduct

	return nil
}

func validateListSymptomProductsRequest(req *pb.ListSymptomProductsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetSymptomKey() == "" {
		return gerr.New("symptom key should not be empty")
	}
	return nil
}

// getSymptomProducts
func getSymptomProducts(products []domain.RecommendedProductIntf, symptom string) []*pb.Product {
	var results []*pb.Product
	for _, v := range products {
		psymtoms := toSymptomsKeys(v.GetSymptoms())
		for _, s := range psymtoms {
			if s == symptom {
				results = append(results, toProtoProductFromRecommendedProduct(v))
				break
			}
		}
	}
	return results
}
