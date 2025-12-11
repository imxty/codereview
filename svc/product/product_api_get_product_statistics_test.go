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

type GetProductStatisticsTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *GetProductStatisticsTestSuite) SetupSuite() {
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

// TestGetProductStatistics
func (suite *GetProductStatisticsTestSuite) TestGetProductStatistics() {
	ctx := context.Background()
	req := &productpb.GetProductStatisticsRequest{
		OrganizationId: organizationId1,
		ProductId:      productId1,
	}
	resp := &productpb.GetProductStatisticsResponse{}
	err := suite.hdl.GetProductStatistics(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(productId1, resp.Product.ProductId)
}

// TestOrganizationIdIsNull
func (suite *GetProductStatisticsTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.GetProductStatisticsRequest{
		OrganizationId: organizationIdIsNull,
		ProductId:      productId1,
	}
	resp := &productpb.GetProductStatisticsResponse{}
	err := suite.hdl.GetProductStatistics(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductIdIsEmpty
func (suite *GetProductStatisticsTestSuite) TestProductIdIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.GetProductStatisticsRequest{
		OrganizationId: organizationId1,
		ProductId:      productIdIsEmpty,
	}
	resp := &productpb.GetProductStatisticsResponse{}
	err := suite.hdl.GetProductStatistics(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductIdIsNotInDb
func (suite *GetProductStatisticsTestSuite) TestProductIdIsNotInDb() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.GetProductStatisticsRequest{
		OrganizationId: organizationId1,
		ProductId:      productIdNotExist,
	}
	resp := &productpb.GetProductStatisticsResponse{}
	err := suite.hdl.GetProductStatistics(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrProductNotFound, err)
}

func (suite *GetProductStatisticsTestSuite) TearDownSuite() {}

func TestGetProductStatisticsTestSuite(t *testing.T) {
	suite.Run(t, new(GetProductStatisticsTestSuite))
}
