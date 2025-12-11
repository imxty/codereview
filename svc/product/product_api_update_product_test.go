package product

import (
	"context"
	"log"
	"testing"
	"time"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/product"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	s3 "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UpdateProductTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *UpdateProductTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, productFile)
	s := store.NewProductStore(conn)
	s3 := &s3.FileStore{}
	suite.s3 = s3
	suite.hdl = NewProductAPIHandler(s, s3, nil, nil)
}

// TestUpdateProductCPD  中药
func (suite *UpdateProductTestSuite) TestUpdateProductCPD() {
	ctx := context.Background()
	drugValidityPeriod := timestamppb.New(time.Now())
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())

	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			ProductId: productIdIsExist1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_CPD,
			// 商品名称
			ProductName: productName1,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: description1,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName1,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumber1,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestUpdateProductHealth  保健品
func (suite *UpdateProductTestSuite) TestUpdateProductHealth() {
	ctx := context.Background()
	drugValidityPeriod := timestamppb.New(time.Now())
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())

	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productIdIsExist2,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT,
			// 商品名称
			ProductName: productName2,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName2,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumber1,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{"JB0001", "JB0002"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestUpdateProductNutrition 营养食品
func (suite *UpdateProductTestSuite) TestUpdateProductNutrition() {
	ctx := context.Background()
	drugValidityPeriod := timestamppb.New(time.Now())
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())

	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productIdIsExist3,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_NUTRITION,
			// 商品名称
			ProductName: productName3,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName1,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumber1,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestUpdateProductService  理疗
func (suite *UpdateProductTestSuite) TestUpdateProductService() {
	ctx := context.Background()
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())

	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productIdIsExist4,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_SERVICE,
			// 商品名称
			ProductName: productName4,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks:     remarksIsEmpty,
			Link:        link,
			SymptomKeys: []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestOrganizationIdIsNull
func (suite *UpdateProductTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())

	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationIdIsNull,
		Product: &productpb.Product{
			// 商品id
			ProductId: productId1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_SERVICE,
			// 商品名称
			ProductName: productName4,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks:     remarksIsEmpty,
			Link:        link,
			SymptomKeys: []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductTypeIsUnset
func (suite *UpdateProductTestSuite) TestProductTypeIsUnset() {
	t := suite.T()
	ctx := context.Background()
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())
	drugValidityPeriod := timestamppb.New(time.Now())
	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productId1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_UNSET,
			// 商品名称
			ProductName: productName1,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: description1,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName3,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumber1,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductNameIsNull
func (suite *UpdateProductTestSuite) TestProductNameIsNull() {
	t := suite.T()
	ctx := context.Background()
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())
	drugValidityPeriod := timestamppb.New(time.Now())
	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productId1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_CPD,
			// 商品名称
			ProductName: productNameIsNull,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName1,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumberIsEmpty,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestApprovedNumberIsEmpty
func (suite *UpdateProductTestSuite) TestApprovedNumberIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())
	drugValidityPeriod := timestamppb.New(time.Now())
	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productId1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_CPD,
			// 商品名称
			ProductName: productName1,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName1,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumberIsEmpty,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestDrugNameIsEmpty
func (suite *UpdateProductTestSuite) TestDrugNameIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())
	drugValidityPeriod := timestamppb.New(time.Now())
	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productId1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_CPD,
			// 商品名称
			ProductName: productName1,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugNameIsEmpty,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumber1,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestDrugValidityPeriodIsNull
func (suite *UpdateProductTestSuite) TestDrugValidityPeriodIsNull() {
	t := suite.T()
	ctx := context.Background()
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())
	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productId1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_CPD,
			// 商品名称
			ProductName: productName1,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName1,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumber1,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: nil,
			Link:               link,
			SymptomKeys:        []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSymptomKeysIsNull
func (suite *UpdateProductTestSuite) TestSymptomKeysIsNull() {
	t := suite.T()
	ctx := context.Background()
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	drugValidityPeriod := timestamppb.New(time.Now())
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())
	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productId1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_CPD,
			// 商品名称
			ProductName: productName1,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName1,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumber1,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestImageMimeIsEmpty
func (suite *UpdateProductTestSuite) TestImageMimeIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())
	drugValidityPeriod := timestamppb.New(time.Now())
	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productId1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_CPD,
			// 商品名称
			ProductName: productName1,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_USING,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName1,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumber1,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mimeIsEmpty,
			Image:    image,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestImageFilenameIsEmpty
func (suite *UpdateProductTestSuite) TestImageFilenameIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	image, err := utils.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())
	drugValidityPeriod := timestamppb.New(time.Now())
	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productId1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_CPD,
			// 商品名称
			ProductName: productName1,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_USING,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName1,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumber1,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    image,
			Filename: filenameIsEmpty,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err = suite.hdl.UpdateProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestImageIsNil
func (suite *UpdateProductTestSuite) TestImageIsNil() {
	t := suite.T()
	ctx := context.Background()
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(imageStorePath, nil).After(utils.RpcLatency())
	drugValidityPeriod := timestamppb.New(time.Now())
	req := &productpb.UpdateProductRequest{
		OrganizationId: organizationId1,
		Product: &productpb.Product{
			// 商品id
			ProductId: productId1,
			// 商品类型
			ProductType: productpb.ProductType_PRODUCT_TYPE_CPD,
			// 商品名称
			ProductName: productName1,
			// 商品状态
			ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_USING,
			// 商品介绍
			Description: descriptionIsEmpty,
			// 商品图片
			ImageUrl: imageUrl,
			// 商品备注
			Remarks: remarksIsEmpty,
			// 药品名称(商品类型为理疗服务方案不用传)
			DrugName: drugName1,
			// 药品是非处方药 (商品类型为理疗服务方案不用传)
			IsOtc: false,
			// 准字号 (商品类型为理疗服务方案不用传)
			ApprovedNumber: approvedNumber1,
			// 药品准效期 (商品类型为理疗服务方案不用传)
			DrugValidityPeriod: drugValidityPeriod,
			Link:               link,
			SymptomKeys:        []string{"JB0001"},
		},
		UploadingImage: &productpb.UploadingImage{
			Mime:     mime,
			Image:    nil,
			Filename: filename,
		},
	}
	resp := &productpb.UpdateProductResponse{}
	err := suite.hdl.UpdateProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *UpdateProductTestSuite) TearDownSuite() {}

func TestUpdateProductTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateProductTestSuite))
}
