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

type ChangeProductStatusByTreatmentIdTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *ChangeProductStatusByTreatmentIdTestSuite) SetupSuite() {
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

// TestChangeProductStatusByTreatmentIdDelete
func (suite *ChangeProductStatusByTreatmentIdTestSuite) TestChangeProductStatusByTreatmentIdDelete() {
	ctx := context.Background()
	req := &productpb.ChangeProductStatusByTreatmentRevIdRequest{
		TreatmentRevId: treatmentId,
		ProductStatus:  productpb.ProductStatus_PRODUCT_STATUS_DELETED,
	}
	resp := &productpb.ChangeProductStatusByTreatmentRevIdResponse{}
	err := suite.hdl.ChangeProductStatusByTreatmentRevId(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestChangeProductStatusByTreatmentIdUsing
func (suite *ChangeProductStatusByTreatmentIdTestSuite) TestChangeProductStatusByTreatmentIdUsing() {
	ctx := context.Background()
	req := &productpb.ChangeProductStatusByTreatmentRevIdRequest{
		TreatmentRevId: treatmentId,
		ProductStatus:  productpb.ProductStatus_PRODUCT_STATUS_USING,
	}
	resp := &productpb.ChangeProductStatusByTreatmentRevIdResponse{}
	err := suite.hdl.ChangeProductStatusByTreatmentRevId(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestChangeProductStatusByTreatmentIdToBeUsed
func (suite *ChangeProductStatusByTreatmentIdTestSuite) TestChangeProductStatusByTreatmentIdToBeUsed() {
	ctx := context.Background()
	req := &productpb.ChangeProductStatusByTreatmentRevIdRequest{
		TreatmentRevId: treatmentId,
		ProductStatus:  productpb.ProductStatus_PRODUCT_STATUS_TO_BE_USED,
	}
	resp := &productpb.ChangeProductStatusByTreatmentRevIdResponse{}
	err := suite.hdl.ChangeProductStatusByTreatmentRevId(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestChangeProductStatusByTreatmentIdUnUsed
func (suite *ChangeProductStatusByTreatmentIdTestSuite) TestChangeProductStatusByTreatmentIdUnUsed() {
	ctx := context.Background()
	req := &productpb.ChangeProductStatusByTreatmentRevIdRequest{
		TreatmentRevId: treatmentId,
		ProductStatus:  productpb.ProductStatus_PRODUCT_STATUS_UNUSED,
	}
	resp := &productpb.ChangeProductStatusByTreatmentRevIdResponse{}
	err := suite.hdl.ChangeProductStatusByTreatmentRevId(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestTenantIdIsNull
func (suite *ChangeProductStatusByTreatmentIdTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ChangeProductStatusByTreatmentRevIdRequest{
		TreatmentRevId: treatmentIdIsNull,
		ProductStatus:  productpb.ProductStatus_PRODUCT_STATUS_DELETED,
	}
	resp := &productpb.ChangeProductStatusByTreatmentRevIdResponse{}
	err := suite.hdl.ChangeProductStatusByTreatmentRevId(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductStatusUnSet
func (suite *ChangeProductStatusByTreatmentIdTestSuite) TestProductStatusUnSet() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ChangeProductStatusByTreatmentRevIdRequest{
		TreatmentRevId: treatmentIdIsNull,
		ProductStatus:  productpb.ProductStatus_PRODUCT_STATUS_UNSET,
	}
	resp := &productpb.ChangeProductStatusByTreatmentRevIdResponse{}
	err := suite.hdl.ChangeProductStatusByTreatmentRevId(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductStatusInvalid
func (suite *ChangeProductStatusByTreatmentIdTestSuite) TestProductStatusInvalid() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ChangeProductStatusByTreatmentRevIdRequest{
		TreatmentRevId: treatmentIdIsNull,
		ProductStatus:  productpb.ProductStatus_PRODUCT_STATUS_INVALID,
	}
	resp := &productpb.ChangeProductStatusByTreatmentRevIdResponse{}
	err := suite.hdl.ChangeProductStatusByTreatmentRevId(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *ChangeProductStatusByTreatmentIdTestSuite) TearDownSuite() {}

func TestChangeProductStatusByTreatmentIdTestSuite(t *testing.T) {
	suite.Run(t, new(ChangeProductStatusByTreatmentIdTestSuite))
}
