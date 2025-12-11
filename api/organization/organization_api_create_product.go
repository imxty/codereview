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

func (s *OrganizationAPIHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest, rsp *pb.CreateProductResponse) error {
	err := validateCreateProductRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.productAPI.CreateProduct(ctx, &productv1.CreateProductRequest{
		OrganizationId: req.GetOrganizationId(),
		Product:        toProductProduct(req.GetProduct()),
		UploadingImage: toProductUploadImage(req.GetUploadingImage()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateCreateProductRequest(req *pb.CreateProductRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetProduct() == nil {
		return gerr.New("product should not be nil")
	}
	if req.GetUploadingImage() == nil {
		return gerr.New("image should not be nil")
	}
	return nil
}
