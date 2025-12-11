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

type UpdateTreatmentStatusTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *UpdateTreatmentStatusTestSuite) SetupSuite() {
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

// TestUpdateTreatmentStatus
func (suite *UpdateTreatmentStatusTestSuite) TestUpdateTreatmentStatus() {
	ctx := context.Background()
	req := &productpb.UpdateTreatmentStatusRequest{
		OrganizationId:  organizationId1,
		TreatmentRevId:  treatmentRevId,
		TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
	}
	resp := &productpb.UpdateTreatmentStatusResponse{}
	err := suite.hdl.UpdateTreatmentStatus(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestTreatmentIdIsEmpty
func (suite *UpdateTreatmentStatusTestSuite) TestTreatmentIdIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.UpdateTreatmentStatusRequest{
		OrganizationId:  organizationId1,
		TreatmentRevId:  treatmentIdIsNull,
		TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
	}
	resp := &productpb.UpdateTreatmentStatusResponse{}
	err := suite.hdl.UpdateTreatmentStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentStatusIsInvalid
func (suite *UpdateTreatmentStatusTestSuite) TestTreatmentStatusIsInvalid() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.UpdateTreatmentStatusRequest{
		OrganizationId:  organizationId1,
		TreatmentRevId:  treatmentRevId,
		TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_INVALID,
	}
	resp := &productpb.UpdateTreatmentStatusResponse{}
	err := suite.hdl.UpdateTreatmentStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentStatusIsUnset
func (suite *UpdateTreatmentStatusTestSuite) TestTreatmentStatusIsUnset() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.UpdateTreatmentStatusRequest{
		OrganizationId:  organizationId1,
		TreatmentRevId:  treatmentRevId,
		TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_UNSET,
	}
	resp := &productpb.UpdateTreatmentStatusResponse{}
	err := suite.hdl.UpdateTreatmentStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOrganizationIdIsNull
func (suite *UpdateTreatmentStatusTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.UpdateTreatmentStatusRequest{
		OrganizationId:  organizationIdIsNull,
		TreatmentRevId:  treatmentRevId,
		TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_UNAPPROVED,
	}
	resp := &productpb.UpdateTreatmentStatusResponse{}
	err := suite.hdl.UpdateTreatmentStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentIdNotExist
func (suite *UpdateTreatmentStatusTestSuite) TestTreatmentIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.UpdateTreatmentStatusRequest{
		OrganizationId:  organizationId1,
		TreatmentRevId:  treatmentIdForCancel,
		TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_UNAPPROVED,
	}
	resp := &productpb.UpdateTreatmentStatusResponse{}
	err := suite.hdl.UpdateTreatmentStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrGetTreatmentFailed, err)
}

func (suite *UpdateTreatmentStatusTestSuite) TearDownSuite() {}

func TestUpdateTreatmentStatusTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTreatmentStatusTestSuite))
}
