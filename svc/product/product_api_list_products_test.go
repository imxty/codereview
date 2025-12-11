package product

import (
	"context"
	"testing"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/product"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	s3 "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type ListProductsTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *ListProductsTestSuite) SetupSuite() {
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

// TestListProducts
func (suite *ListProductsTestSuite) TestListProducts() {
	ctx := context.Background()
	req := &productpb.ListProductsRequest{
		OrganizationId: organizationId1,
		Pagination: &productpb.Pagination{
			Offset: offset,
			Size:   size,
		},
		GetAll:        true,
		ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_USING,
	}
	resp := &productpb.ListProductsResponse{}
	err := suite.hdl.ListProducts(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp.Products)
	suite.Assert().Equal(int32(8), resp.TotalCount)
	suite.Assert().Equal(int32(3), resp.ProductUsingCount)
}

// TestTenantIdIsNull
func (suite *ListProductsTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ListProductsRequest{
		OrganizationId: organizationId1,
		Pagination: &productpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &productpb.ListProductsResponse{}
	err := suite.hdl.ListProducts(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOffsetIsNegative
func (suite *ListProductsTestSuite) TestOffsetIsNegative() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ListProductsRequest{
		OrganizationId: organizationId1,
		Pagination: &productpb.Pagination{
			Offset: offsetIsNegative,
			Size:   size,
		},
	}
	resp := &productpb.ListProductsResponse{}
	err := suite.hdl.ListProducts(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSizeIsNegative
func (suite *ListProductsTestSuite) TestSizeIsNegative() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ListProductsRequest{
		OrganizationId: organizationId1,
		Pagination: &productpb.Pagination{
			Offset: offset,
			Size:   sizeIsNegative,
		},
	}
	resp := &productpb.ListProductsResponse{}
	err := suite.hdl.ListProducts(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestGetAllIsTrue
func (suite *ListProductsTestSuite) TestGetAllIsTrue() {
	ctx := context.Background()
	req := &productpb.ListProductsRequest{
		OrganizationId: organizationId1,
		Pagination: &productpb.Pagination{
			Offset: offset,
			Size:   sizeIsNegative,
		},
		GetAll: true,
	}
	resp := &productpb.ListProductsResponse{}
	err := suite.hdl.ListProducts(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(8, len(resp.GetProducts()))
}

func (suite *ListProductsTestSuite) TearDownSuite() {}

func TestListProductsTestSuite(t *testing.T) {
	suite.Run(t, new(ListProductsTestSuite))
}
