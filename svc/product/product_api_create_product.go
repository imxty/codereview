package product

import (
	"bytes"
	"context"
	gerr "errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/image"
	"github.com/jinmukeji/huimaibao-service/pkg/time"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/huimaibao-service/svc/summary"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

const (
	imageSizeLimit = 2 * 1024 * 1024
)

// CreateProduct
func (s *ProductAPIHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest, rsp *pb.CreateProductResponse) error {
	err := validateCreateProductRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取图片url
	imageUrl, err := s.uploadImage(req.GetUploadingImage())
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}

	// 转换 proto->domain
	pProduct := req.GetProduct()
	pt, err := toDomainProductType(pProduct.GetProductType())
	if err != nil {
		return err
	}
	// 获取商品类型 判断是否需要传入药品准效期
	drugValidityPeriod := time.InitTime
	productType, _ := toDomainProductType(pProduct.GetProductType())
	if productType == domain.ProductTypeCpd {
		drugValidityPeriod = pProduct.GetDrugValidityPeriod().AsTime()
	}

	// 构建domain的product
	rp := &domain.RecommendedProduct{
		RecommendedProductID: xid.New().String(),
		OrganizationID:       req.GetOrganizationId(),
		ProductImageUrl:      imageUrl,
		ProductName:          pProduct.GetProductName(),
		ProductType:          pt,
		ProductStatus:        domain.ProductStatusUnused,
		ProductIntroduction:  pProduct.GetDescription(),
		Remarks:              pProduct.GetRemarks(),
		DrugName:             pProduct.GetDrugName(),
		DrugClassification:   toDomainDrugClassification[pProduct.GetIsOtc()],
		ApprovedNumber:       pProduct.GetApprovedNumber(),
		DrugValidityPeriod:   drugValidityPeriod,
		Symptoms:             getSymptoms(pProduct.GetSymptomKeys()),
		Rev:                  0,
	}

	// 添加商品
	err = s.productStore.AddProduct(ctx, req.GetOrganizationId(), rp)
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "add product failed[%s]", err.Error())
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
	product := req.GetProduct()
	if len(product.GetSymptomKeys()) == 0 {
		return gerr.New("symptom keys should not be empty")
	}
	// 验证keys
	for _, v := range product.GetSymptomKeys() {
		find := false
		// 验证risk
		for risk := range summary.RiskMap {
			if v == risk {
				find = true
				break
			}
		}
		// 验证脏腑
		for dd := range summary.DirtyDialecticsMap {
			if v == dd {
				find = true
				break
			}
		}
		for pt := range summary.PhysicalTherapyMap {
			if v == pt {
				find = true
				break
			}
		}
		for pt := range summary.PhysiqueMap {
			if v == pt {
				find = true
				break
			}
		}
		if !find {
			return gerr.New("invalid symptom key")
		}
	}
	pt, ptErr := toDomainProductType(product.GetProductType())
	if ptErr != nil {
		return gerr.New("invalid product type")
	}
	if pt == domain.ProductTypeCpd {
		if product.GetDrugName() == "" || utf8.RuneCountInString(product.GetDrugName()) > 50 {
			return gerr.New("invalid drug name")
		}
		if product.GetApprovedNumber() == "" {
			return gerr.New("approved number should not be empty")
		}
		if product.GetDrugValidityPeriod() == nil {
			return gerr.New("drug validity period should not be nil")
		}
	}

	if product.GetProductName() == "" || utf8.RuneCountInString(product.GetProductName()) > 50 {
		return gerr.New("invalid product name")
	}
	// 验证商品介绍和备注的长度
	if utf8.RuneCountInString(product.GetDescription()) > 100 {
		return gerr.New("description should not exceed 100 words")
	}
	if utf8.RuneCountInString(product.GetRemarks()) > 100 {
		return gerr.New("remarks should not exceed 100 words")
	}
	// 验证上传的图片
	if req.GetUploadingImage().GetImage() == nil {
		return gerr.New("image should not be nil")
	}
	if len(req.GetUploadingImage().GetImage()) >= imageSizeLimit {
		return gerr.New("image size should not exceed 2M")
	}
	if req.GetUploadingImage().GetFilename() == "" {
		return gerr.New("image file name should not be empty")
	}
	if req.GetUploadingImage().GetMime() == "" {
		return gerr.New("image mime should not be empty")
	}
	return nil
}

// 上传图片
func (u *ProductAPIHandler) uploadImage(logo *pb.UploadingImage) (string, error) {
	const (
		ImageSizeLimit = 2 * 1024 * 1024
	)
	logoUrl := ""
	if logo == nil {
		return "", nil
	}
	// 检查mime和大小
	if len(logo.GetImage()) > ImageSizeLimit {
		return "", gerr.New("image should not exceed 2M")
	}
	// 检测mime
	suffix, err := image.GetImageSuffix(logo.GetMime())
	if err != nil {
		return "", err
	}
	// 如果传了logo就去上传logo
	// 生成图片的ID
	logoUrl = fmt.Sprintf("%s.%s", xid.New().String(), suffix)
	path, err := u.s3Store.Save(logoUrl, logo.GetMime(), byteReader(logo.GetImage()))
	if err != nil {
		return "", err
	}
	return path, nil
}

// byteReader 将字节格式化为可读文本流
func byteReader(img []byte) io.ReadSeeker {
	return bytes.NewReader(img)
}

// getSymptoms
func getSymptoms(s []string) string {
	return strings.Join(s, ",")
}

// toSymptomsKeys
func toSymptomsKeys(s string) []string {
	return strings.Split(s, ",")
}
