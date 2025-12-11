package organization

import (
	"context"
	gerr "errors"
	"unicode/utf8"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	productv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// 图片大小
	imageSizeLimit = 2 * 1024 * 1024
)

func (s *OrganizationAPIHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest, rsp *pb.UpdateProductResponse) error {
	err := validateUpdateProductRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	_, err = s.productAPI.UpdateProduct(ctx, &productv1.UpdateProductRequest{
		// 组织id
		OrganizationId: req.GetOrganizationId(),
		// 商品
		Product: toProductProduct(req.GetProduct()),
		// 商品上传图片
		UploadingImage: toProductUploadImage(req.GetUploadingImage()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	return nil
}

// 验证request
func validateUpdateProductRequest(req *pb.UpdateProductRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	p := req.GetProduct()
	if p.GetProductId() == "" {
		return gerr.New("product_id should not be empty")
	}
	if p.GetProductName() == "" || utf8.RuneCountInString(p.GetProductName()) > 50 {
		return gerr.New("invalid product name")
	}
	// 验证商品介绍和备注的长度
	if utf8.RuneCountInString(p.GetDescription()) > 100 {
		return gerr.New("description should not exceed 100 words")
	}
	if utf8.RuneCountInString(p.GetRemarks()) > 100 {
		return gerr.New("remarks should not exceed 100 words")
	}
	// 验证上传的图片
	if req.GetUploadingImage().GetImage() != nil {
		if len(req.GetUploadingImage().GetImage()) >= imageSizeLimit {
			return gerr.New("image size should not exceed 2M")
		}
		if req.GetUploadingImage().GetFilename() == "" {
			return gerr.New("image file name should not be empty")
		}
		if req.GetUploadingImage().GetMime() == "" {
			return gerr.New("image mime should not be empty")
		}
	}
	return nil
}
