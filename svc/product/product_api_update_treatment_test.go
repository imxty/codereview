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

type UpdateTreatmentTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *UpdateTreatmentTestSuite) SetupSuite() {
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

// TestUpdateTreatment
func (suite *UpdateTreatmentTestSuite) TestUpdateTreatment() {
	ctx := context.Background()
	req := &productpb.UpdateTreatmentRequest{
		Treatment: &productpb.Treatment{
			OrganizationId:  organizationId1,
			TreatmentId:     treatmentId1,
			TreatmentName:   treatmentName,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     true,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsRiskyDisease: []*productpb.TreatmentItem{
				{
					Symptom: symptomKeyRisky,
					Products: []*productpb.Product{
						{ProductId: productId1},
						{ProductId: productId2},
					},
				},
			},
		},
	}
	resp := &productpb.UpdateTreatmentResponse{}
	err := suite.hdl.UpdateTreatment(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestOrganizationIdIsNull
func (suite *UpdateTreatmentTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.UpdateTreatmentRequest{
		Treatment: &productpb.Treatment{
			OrganizationId:  organizationIdIsNull,
			TreatmentId:     treatmentId,
			TreatmentName:   treatmentName,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     true,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsRiskyDisease: []*productpb.TreatmentItem{
				{
					Symptom: symptomKeyRisky,
					Products: []*productpb.Product{
						{ProductId: productId1},
						{ProductId: productId2},
					},
				},
			},
		},
	}
	resp := &productpb.UpdateTreatmentResponse{}
	err := suite.hdl.UpdateTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentIdIsEmpty
func (suite *UpdateTreatmentTestSuite) TestTreatmentIdIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.UpdateTreatmentRequest{
		Treatment: &productpb.Treatment{
			OrganizationId:  organizationId1,
			TreatmentId:     treatmentIdIsNull,
			TreatmentName:   treatmentName,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     true,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsRiskyDisease: []*productpb.TreatmentItem{
				{
					Symptom: symptomKeyRisky,
					Products: []*productpb.Product{
						{ProductId: productId1},
						{ProductId: productId2},
					},
				},
			},
		},
	}
	resp := &productpb.UpdateTreatmentResponse{}
	err := suite.hdl.UpdateTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSymptomIsEmpty
func (suite *UpdateTreatmentTestSuite) TestSymptomIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.UpdateTreatmentRequest{
		Treatment: &productpb.Treatment{
			TreatmentId:     treatmentId,
			TreatmentName:   treatmentName,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     true,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsRiskyDisease: []*productpb.TreatmentItem{
				{
					Symptom: symptomIsEmpty,
					Products: []*productpb.Product{
						{ProductId: productId1},
						{ProductId: productId2},
					},
				},
			},
		},
	}
	resp := &productpb.UpdateTreatmentResponse{}
	err := suite.hdl.UpdateTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentItemIsEmpty
func (suite *UpdateTreatmentTestSuite) TestTreatmentItemIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.UpdateTreatmentRequest{
		Treatment: &productpb.Treatment{
			TreatmentId:                treatmentId,
			TreatmentName:              treatmentName,
			ReviewPass:                 false,
			ReviewComment:              reviewComment,
			IsPublished:                true,
			TreatmentStatus:            productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsRiskyDisease: []*productpb.TreatmentItem{},
		},
	}
	resp := &productpb.UpdateTreatmentResponse{}
	err := suite.hdl.UpdateTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductIsEmpty
func (suite *UpdateTreatmentTestSuite) TestProductIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.UpdateTreatmentRequest{
		Treatment: &productpb.Treatment{
			TreatmentId:     treatmentId,
			TreatmentName:   treatmentName,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     true,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsRiskyDisease: []*productpb.TreatmentItem{
				{
					Symptom:  symptomKeyRisky,
					Products: []*productpb.Product{},
				},
			},
		},
	}
	resp := &productpb.UpdateTreatmentResponse{}
	err := suite.hdl.UpdateTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *UpdateTreatmentTestSuite) TearDownSuite() {}

func TestUpdateTreatmentTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTreatmentTestSuite))
}
