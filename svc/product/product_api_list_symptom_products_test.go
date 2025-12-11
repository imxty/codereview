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

type ListSymptomProductsTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *ListSymptomProductsTestSuite) SetupSuite() {
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

// TestListSymptomProducts
func (suite *ListSymptomProductsTestSuite) TestListSymptomProducts() {
	ctx := context.Background()
	req := &productpb.ListSymptomProductsRequest{
		OrganizationId: organizationId1,
		SymptomKey:     symptomKeyRisky,
	}
	resp := &productpb.ListSymptomProductsResponse{}
	err := suite.hdl.ListSymptomProducts(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp.Products)
}

// TestOrganizationIdIsNull
func (suite *ListSymptomProductsTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ListSymptomProductsRequest{
		OrganizationId: organizationIdIsNull,
		SymptomKey:     symptomKeyRisky,
	}
	resp := &productpb.ListSymptomProductsResponse{}
	err := suite.hdl.ListSymptomProducts(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSymptomIsEmpty
func (suite *ListSymptomProductsTestSuite) TestSymptomIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ListSymptomProductsRequest{
		OrganizationId: organizationId1,
		SymptomKey:     symptomIsEmpty,
	}
	resp := &productpb.ListSymptomProductsResponse{}
	err := suite.hdl.ListSymptomProducts(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListSymptomProductsTestSuite) TearDownSuite() {}

func TestListSymptomProductsTestSuite(t *testing.T) {
	suite.Run(t, new(ListSymptomProductsTestSuite))
}
