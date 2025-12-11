package product

import (
	"context"
	gerr "errors"
	"unicode/utf8"

	"github.com/fatih/structs"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	productImageUrl = "ProductImageUrl"
)

func (s *ProductAPIHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest, rsp *pb.UpdateProductResponse) error {
	// 验证request
	err := validateUpdateProductRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取更新前的商品信息
	bProduct, err := s.productStore.GetProduct(ctx, req.GetProduct().GetProductId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 判断当前商品的状态 只有未使用状态下的商品才能被更新
	if bProduct.GetProductStatus() != domain.ProductStatusUnused {
		return errors.Errorf(codes.InvalidRequest, "update product failed for invalid product_status product:[%s]", req.GetProduct().GetProductId())
	}
	// 获取更新前版本号
	rev := bProduct.GetRev()

	imageUrl := ""
	if req.GetUploadingImage() != nil {
		// 获取需要更新的图片url
		imageUrl, err = s.uploadImage(req.GetUploadingImage())
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
	} else {
		imageUrl = bProduct.GetProductImageUrl()
	}

	// 转换 proto->domain
	pProduct := req.GetProduct()
	pt, err := toDomainProductType(pProduct.GetProductType())
	if err != nil {
		return err
	}
	drugValidityPeriod := pProduct.GetDrugValidityPeriod().AsTime()

	// 构建domain的product
	rp := &domain.RecommendedProduct{
		RecommendedProductID: req.GetProduct().GetProductId(),
		OrganizationID:       req.GetOrganizationId(),
		ProductImageUrl:      imageUrl,
		ProductName:          pProduct.GetProductName(),
		ProductType:          pt,
		ProductStatus:        bProduct.GetProductStatus(),
		ProductIntroduction:  pProduct.GetDescription(),
		Remarks:              pProduct.GetRemarks(),
		DrugName:             pProduct.GetDrugName(),
		DrugClassification:   toDomainDrugClassification[pProduct.GetIsOtc()],
		ApprovedNumber:       pProduct.GetApprovedNumber(),
		DrugValidityPeriod:   drugValidityPeriod,
		Symptoms:             getSymptoms(pProduct.GetSymptomKeys()),
		Rev:                  rev,
	}
	// struct->map获得更新的map
	updates := structs.Map(rp)
	// 判断图片地址是否需要更新
	// 图片地址为空不需要更新
	if imageUrl == "" {
		delete(updates, productImageUrl)
	}

	// 更新单个商品
	err = s.productStore.UpdateProduct(ctx, req.GetProduct().GetProductId(), updates, rev)
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "update product failed:[%s]", err.Error())
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
	// 验证商品类型
	pt, ptErr := toDomainProductType(p.GetProductType())
	if ptErr != nil {
		return gerr.New("invalid product type")
	}
	// 验证不同商品类型的必填项
	if pt == domain.ProductTypeHealthCareProduct || pt == domain.ProductTypeNutrition || pt == domain.ProductTypeCpd {
		if p.GetApprovedNumber() == "" {
			return gerr.New("approved number should not be empty")
		}
	}
	if pt == domain.ProductTypeCpd {
		if p.GetDrugName() == "" || utf8.RuneCountInString(p.GetDrugName()) > 50 {
			return gerr.New("invalid drug name")
		}
		if p.GetDrugValidityPeriod() == nil {
			return gerr.New("drug validity period should not be nil")
		}
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
