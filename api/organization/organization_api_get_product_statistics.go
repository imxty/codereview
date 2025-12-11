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

// GetProductStatistics
func (s *OrganizationAPIHandler) GetProductStatistics(ctx context.Context, req *pb.GetProductStatisticsRequest, rsp *pb.GetProductStatisticsResponse) error {
	err := validateGetProductStatisticsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	getRsp, err := s.productAPI.GetProductStatistics(ctx, &productv1.GetProductStatisticsRequest{
		OrganizationId: req.GetOrganizationId(),
		ProductId:      req.GetProductId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 转换
	pProduct := toAppProduct(getRsp.GetProduct(), s.s3Domain)
	// 转换
	se := make([]*pb.SymptomExposure, len(getRsp.GetSymptomExposure()))
	for i, v := range getRsp.GetSymptomExposure() {
		se[i] = toAppSymptomExposure(v)
	}
	// 转换
	pe := make([]*pb.ProductExposure, len(getRsp.GetProductExposure()))
	for i, v := range getRsp.GetProductExposure() {
		pe[i] = toAppProductExposure(v)
	}

	// 返回结果
	rsp.Product = pProduct
	rsp.ProductInstantTotalCount = getRsp.GetProductInstantTotalCount()
	rsp.SymptomExposure = se
	rsp.ProductExposure = pe
	return nil
}

// 验证request
func validateGetProductStatisticsRequest(req *pb.GetProductStatisticsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetProductId() == "" {
		return gerr.New("product id should not be empty")
	}
	return nil
}
