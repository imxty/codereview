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

type GetProductTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *GetProductTestSuite) SetupSuite() {
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

// TestGetProduct
func (suite *GetProductTestSuite) TestGetProduct() {
	ctx := context.Background()
	req := &productpb.GetProductRequest{
		ProductId: productId1,
	}
	resp := &productpb.GetProductResponse{}
	err := suite.hdl.GetProduct(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(productId1, resp.Product.ProductId)
}

// TestProductIdIsEmpty
func (suite *GetProductTestSuite) TestProductIdIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.GetProductRequest{
		ProductId: productIdIsEmpty,
	}
	resp := &productpb.GetProductResponse{}
	err := suite.hdl.GetProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductIdIsNotInDb
func (suite *GetProductTestSuite) TestProductIdIsNotInDb() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.GetProductRequest{
		ProductId: productIdNotExist,
	}
	resp := &productpb.GetProductResponse{}
	err := suite.hdl.GetProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrProductNotFound, err)
}

func (suite *GetProductTestSuite) TearDownSuite() {}

func TestGetProductTestSuite(t *testing.T) {
	suite.Run(t, new(GetProductTestSuite))
}
