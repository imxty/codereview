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

type DeleteTreatmentTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *DeleteTreatmentTestSuite) SetupSuite() {
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

// TestDeleteTreatment
func (suite *DeleteTreatmentTestSuite) TestDeleteTreatment() {
	ctx := context.Background()
	req := &productpb.DeleteTreatmentRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentIdIsDraft,
	}
	resp := &productpb.DeleteTreatmentResponse{}
	err := suite.hdl.DeleteTreatment(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestDeleteTreatmentErr
func (suite *DeleteTreatmentTestSuite) TestDeleteTreatmentErr() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.DeleteTreatmentRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentId,
	}
	resp := &productpb.DeleteTreatmentResponse{}
	err := suite.hdl.DeleteTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidOperation, err)
}

// TestOrganizationIdIsNull
func (suite *DeleteTreatmentTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.DeleteTreatmentRequest{
		OrganizationId: organizationIdIsNull,
		TreatmentId:    treatmentId,
	}
	resp := &productpb.DeleteTreatmentResponse{}
	err := suite.hdl.DeleteTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentIdIsEmpty
func (suite *DeleteTreatmentTestSuite) TestTreatmentIdIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.DeleteTreatmentRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentIdIsNull,
	}
	resp := &productpb.DeleteTreatmentResponse{}
	err := suite.hdl.DeleteTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *DeleteTreatmentTestSuite) TearDownSuite() {}

func TestDeleteTreatmentTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteTreatmentTestSuite))
}
