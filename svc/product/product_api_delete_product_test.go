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

type DeleteProductTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *DeleteProductTestSuite) SetupSuite() {
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

// TestDeleteProduct
func (suite *DeleteProductTestSuite) TestDeleteProduct() {
	ctx := context.Background()
	req := &productpb.DeleteProductRequest{
		OrganizationId: organizationId1,
		ProductId:      productIdIsExist1,
	}
	resp := &productpb.DeleteProductResponse{}
	err := suite.hdl.DeleteProduct(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestOrganizationIdIsNull
func (suite *DeleteProductTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.DeleteProductRequest{
		OrganizationId: organizationIdIsNull,
		ProductId:      productId1,
	}
	resp := &productpb.DeleteProductResponse{}
	err := suite.hdl.DeleteProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductIdIsEmpty
func (suite *DeleteProductTestSuite) TestProductIdIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.DeleteProductRequest{
		OrganizationId: organizationId1,
		ProductId:      productIdIsEmpty,
	}
	resp := &productpb.DeleteProductResponse{}
	err := suite.hdl.DeleteProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductIdNotFound
func (suite *DeleteProductTestSuite) TestProductIdNotFound() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.DeleteProductRequest{
		OrganizationId: organizationId1,
		ProductId:      productIdNotExist,
	}
	resp := &productpb.DeleteProductResponse{}
	err := suite.hdl.DeleteProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrProductNotFound, err)
}

// TestProductIdIsErr
func (suite *DeleteProductTestSuite) TestProductIdIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.DeleteProductRequest{
		OrganizationId: organizationId1,
		ProductId:      productId1,
	}
	resp := &productpb.DeleteProductResponse{}
	err := suite.hdl.DeleteProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductIdNotBelongToOrganization
func (suite *DeleteProductTestSuite) TestProductIdNotBelongToOrganization() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.DeleteProductRequest{
		OrganizationId: organizationId2,
		ProductId:      productId1,
	}
	resp := &productpb.DeleteProductResponse{}
	err := suite.hdl.DeleteProduct(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrProductNotBelongToOrganization, err)
}

func (suite *DeleteProductTestSuite) TearDownSuite() {}

func TestDeleteProductTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteProductTestSuite))
}
