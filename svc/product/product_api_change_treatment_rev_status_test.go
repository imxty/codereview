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

type ChangeTreatmentRevStatusTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *ChangeTreatmentRevStatusTestSuite) SetupSuite() {
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

// TestChangeTreatmentRevStatusDraft
func (suite *ChangeTreatmentRevStatusTestSuite) TestChangeTreatmentRevStatusDraft() {
	ctx := context.Background()
	req := &productpb.ChangeTreatmentRevStatusRequest{
		TreatmentRevId: treatmentRevId,
		Status:         productpb.TreatmentStatus_TREATMENT_STATUS_DRAFT,
	}
	resp := &productpb.ChangeTreatmentRevStatusResponse{}
	err := suite.hdl.ChangeTreatmentRevStatus(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestChangeTreatmentRevStatusReview
func (suite *ChangeTreatmentRevStatusTestSuite) TestChangeTreatmentRevStatusReview() {
	ctx := context.Background()
	req := &productpb.ChangeTreatmentRevStatusRequest{
		TreatmentRevId: treatmentRevId,
		Status:         productpb.TreatmentStatus_TREATMENT_STATUS_UNDER_REVIEW,
	}
	resp := &productpb.ChangeTreatmentRevStatusResponse{}
	err := suite.hdl.ChangeTreatmentRevStatus(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestChangeTreatmentRevStatusApproved
func (suite *ChangeTreatmentRevStatusTestSuite) TestChangeTreatmentRevStatusApproved() {
	ctx := context.Background()
	req := &productpb.ChangeTreatmentRevStatusRequest{
		TreatmentRevId: treatmentRevId,
		Status:         productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
	}
	resp := &productpb.ChangeTreatmentRevStatusResponse{}
	err := suite.hdl.ChangeTreatmentRevStatus(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestChangeTreatmentRevStatusUnApproved
func (suite *ChangeTreatmentRevStatusTestSuite) TestChangeTreatmentRevStatusUnApproved() {
	ctx := context.Background()
	req := &productpb.ChangeTreatmentRevStatusRequest{
		TreatmentRevId: treatmentRevId,
		Status:         productpb.TreatmentStatus_TREATMENT_STATUS_UNAPPROVED,
	}
	resp := &productpb.ChangeTreatmentRevStatusResponse{}
	err := suite.hdl.ChangeTreatmentRevStatus(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestTreatmentRevIdIsNull
func (suite *ChangeTreatmentRevStatusTestSuite) TestTreatmentRevIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ChangeTreatmentRevStatusRequest{
		TreatmentRevId: treatmentRevIdIsNull,
		Status:         productpb.TreatmentStatus_TREATMENT_STATUS_DRAFT,
	}
	resp := &productpb.ChangeTreatmentRevStatusResponse{}
	err := suite.hdl.ChangeTreatmentRevStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductStatusUnSet
func (suite *ChangeTreatmentRevStatusTestSuite) TestProductStatusUnSet() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ChangeTreatmentRevStatusRequest{
		TreatmentRevId: treatmentRevId,
		Status:         productpb.TreatmentStatus_TREATMENT_STATUS_UNSET,
	}
	resp := &productpb.ChangeTreatmentRevStatusResponse{}
	err := suite.hdl.ChangeTreatmentRevStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductStatusInvalid
func (suite *ChangeTreatmentRevStatusTestSuite) TestProductStatusInvalid() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ChangeTreatmentRevStatusRequest{
		TreatmentRevId: treatmentRevId,
		Status:         productpb.TreatmentStatus_TREATMENT_STATUS_INVALID,
	}
	resp := &productpb.ChangeTreatmentRevStatusResponse{}
	err := suite.hdl.ChangeTreatmentRevStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *ChangeTreatmentRevStatusTestSuite) TearDownSuite() {}

func TestChangeTreatmentRevStatusTestSuite(t *testing.T) {
	suite.Run(t, new(ChangeTreatmentRevStatusTestSuite))
}
