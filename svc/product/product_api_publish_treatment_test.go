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

type PublishTreatmentTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *PublishTreatmentTestSuite) SetupSuite() {
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

// TestPublishTreatment
func (suite *PublishTreatmentTestSuite) TestPublishTreatment() {
	ctx := context.Background()
	req := &productpb.PublishTreatmentRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentIdIsPassed,
	}
	resp := &productpb.PublishTreatmentResponse{}
	err := suite.hdl.PublishTreatment(ctx, req, resp)
	suite.Assert().NoError(err)
}

func (suite *PublishTreatmentTestSuite) TestTreatmentErr() {
	ctx := context.Background()
	t := suite.T()
	req := &productpb.PublishTreatmentRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentId,
	}
	resp := &productpb.PublishTreatmentResponse{}
	err := suite.hdl.PublishTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrPublishTreatmentFailed, err)
}

// TestOrganizationIdIsNull
func (suite *PublishTreatmentTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.PublishTreatmentRequest{
		OrganizationId: organizationIdIsNull,
		TreatmentId:    treatmentId,
	}
	resp := &productpb.PublishTreatmentResponse{}
	err := suite.hdl.PublishTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentIdIsNull
func (suite *PublishTreatmentTestSuite) TestTreatmentIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.PublishTreatmentRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentIdIsNull,
	}
	resp := &productpb.PublishTreatmentResponse{}
	err := suite.hdl.PublishTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *PublishTreatmentTestSuite) TearDownSuite() {}

func TestPublishTreatmentTestSuite(t *testing.T) {
	suite.Run(t, new(PublishTreatmentTestSuite))
}
