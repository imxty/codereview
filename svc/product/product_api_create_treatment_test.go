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

type CreateTreatmentTestSuite struct {
	suite.Suite
	hdl *ProductAPIHandler
	s3  *s3.FileStore
}

func (suite *CreateTreatmentTestSuite) SetupSuite() {
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

// TestCreateTreatment
func (suite *CreateTreatmentTestSuite) TestCreateTreatment() {
	ctx := context.Background()
	req := &productpb.CreateTreatmentRequest{
		OrganizationId: organizationId1,
		Treatment: &productpb.Treatment{
			TreatmentName:   treatmentName,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     false,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_DRAFT,
			TreatmentItemsRiskyDisease: []*productpb.TreatmentItem{
				{
					Symptom: symptomKeyRisky,
					Products: []*productpb.Product{
						{ProductId: productId1},
						{ProductId: productId2},
					},
				},
			},
			TreatmentItemsDirtyDialectic: []*productpb.TreatmentItem{
				{
					Symptom: symptomKeyDirtyDialectic,
					Products: []*productpb.Product{
						{ProductId: productId1},
						{ProductId: productId2},
					},
				},
			},
			TreatmentItemsPhysicalTherapy: []*productpb.TreatmentItem{
				{
					Symptom: symptomKeyPhysicalTherapy,
					Products: []*productpb.Product{
						{ProductId: productId1},
						{ProductId: productId2},
					},
				},
			},
		},
	}
	resp := &productpb.CreateTreatmentResponse{}
	err := suite.hdl.CreateTreatment(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp.Treatment)
	suite.Assert().Equal(treatmentName, resp.Treatment.TreatmentName)

}

// TestOrganizationIdIsNull
func (suite *CreateTreatmentTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.CreateTreatmentRequest{
		OrganizationId: organizationIdIsNull,
		Treatment: &productpb.Treatment{
			TreatmentName:   treatmentName,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     true,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsDirtyDialectic: []*productpb.TreatmentItem{
				{
					Symptom: symptomKeyDirtyDialectic,
					Products: []*productpb.Product{
						{ProductId: productId1},
						{ProductId: productId2},
					},
				},
			},
		},
	}
	resp := &productpb.CreateTreatmentResponse{}
	err := suite.hdl.CreateTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentNameIsEmpty
func (suite *CreateTreatmentTestSuite) TestTreatmentNameIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.CreateTreatmentRequest{
		OrganizationId: organizationId1,
		Treatment: &productpb.Treatment{
			TreatmentName:   treatmentNameIsEmpty,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     true,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsDirtyDialectic: []*productpb.TreatmentItem{
				{
					Symptom: symptomKeyDirtyDialectic,
					Products: []*productpb.Product{
						{ProductId: productId1},
						{ProductId: productId2},
					},
				},
			},
		},
	}
	resp := &productpb.CreateTreatmentResponse{}
	err := suite.hdl.CreateTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentItemIsEmpty
func (suite *CreateTreatmentTestSuite) TestTreatmentItemIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.CreateTreatmentRequest{
		OrganizationId: organizationId1,
		Treatment: &productpb.Treatment{
			TreatmentName:                treatmentName,
			ReviewPass:                   false,
			ReviewComment:                reviewComment,
			IsPublished:                  true,
			TreatmentStatus:              productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsDirtyDialectic: []*productpb.TreatmentItem{},
		},
	}
	resp := &productpb.CreateTreatmentResponse{}
	err := suite.hdl.CreateTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestProductsOutOfLimit
func (suite *CreateTreatmentTestSuite) TestProductsOutOfLimit() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.CreateTreatmentRequest{
		OrganizationId: organizationId1,
		Treatment: &productpb.Treatment{
			TreatmentName:   treatmentName,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     true,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsDirtyDialectic: []*productpb.TreatmentItem{{
				Symptom: symptomKeyDirtyDialectic,
				Products: []*productpb.Product{
					{ProductId: productId1},
					{ProductId: productId2},
					{ProductId: productId3},
					{ProductId: productId4},
					{ProductId: productId5},
					{ProductId: productId6},
					{ProductId: productId7},
					{ProductId: productId8},
					{ProductId: productId9},
					{ProductId: productId10},
					{ProductId: productId11},
				},
			}},
			TreatmentItemsPhysicalTherapy: []*productpb.TreatmentItem{{
				Symptom: symptomKeyDirtyDialectic,
				Products: []*productpb.Product{
					{ProductId: productId12},
					{ProductId: productId13},
					{ProductId: productId14},
					{ProductId: productId15},
					{ProductId: productId16},
					{ProductId: productId17},
					{ProductId: productId18},
					{ProductId: productId19},
					{ProductId: productId20},
					{ProductId: productId22},
					{ProductId: productId23},
					{ProductId: productId24},
					{ProductId: productId25},
					{ProductId: productId26},
					{ProductId: productId27},
					{ProductId: productId28},
					{ProductId: productId29},
					{ProductId: productId30},
					{ProductId: productId32},
					{ProductId: productId33},
					{ProductId: productId34},
					{ProductId: productId35},
					{ProductId: productId36},
					{ProductId: productId37},
					{ProductId: productId38},
					{ProductId: productId39},
					{ProductId: productId40},
				},
			}},
		},
	}
	resp := &productpb.CreateTreatmentResponse{}
	err := suite.hdl.CreateTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestManySymptoms
func (suite *CreateTreatmentTestSuite) TestManySymptoms() {
	ctx := context.Background()
	req := &productpb.CreateTreatmentRequest{
		OrganizationId: organizationId1,
		Treatment: &productpb.Treatment{
			TreatmentName:   treatmentName,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     true,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsDirtyDialectic: []*productpb.TreatmentItem{
				{
					Symptom: symptomKeyDirtyDialectic,
					Products: []*productpb.Product{
						{ProductId: productId1},
						{ProductId: productId2},
						{ProductId: productId3},
						{ProductId: productId4},
						{ProductId: productId5},
						{ProductId: productId6},
						{ProductId: productId7},
						{ProductId: productId8},
						{ProductId: productId9},
						{ProductId: productId10},
					},
				},
				{
					Symptom:  symptomKeyDirtyDialectic2,
					Products: []*productpb.Product{{ProductId: productId1}},
				},
				{
					Symptom:  symptomKeyDirtyDialectic3,
					Products: []*productpb.Product{{ProductId: productId2}},
				},
				{
					Symptom:  symptomKeyDirtyDialectic4,
					Products: []*productpb.Product{{ProductId: productId3}},
				},
				{
					Symptom:  symptomKeyDirtyDialectic5,
					Products: []*productpb.Product{},
				},
				{
					Symptom: symptomKeyDirtyDialectic6,
				},
			},
		},
	}
	resp := &productpb.CreateTreatmentResponse{}
	err := suite.hdl.CreateTreatment(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestNoProducts
func (suite *CreateTreatmentTestSuite) TestNoProducts() {
	ctx := context.Background()
	req := &productpb.CreateTreatmentRequest{
		OrganizationId: organizationId1,
		Treatment: &productpb.Treatment{
			TreatmentName:   treatmentName,
			ReviewPass:      false,
			ReviewComment:   reviewComment,
			IsPublished:     true,
			TreatmentStatus: productpb.TreatmentStatus_TREATMENT_STATUS_APPROVED,
			TreatmentItemsDirtyDialectic: []*productpb.TreatmentItem{
				{
					Symptom:  symptomKeyDirtyDialectic,
					Products: []*productpb.Product{{ProductId: productId3}},
				},
				{
					Symptom:  symptomKeyDirtyDialectic2,
					Products: []*productpb.Product{},
				},
				{
					Symptom: symptomKeyDirtyDialectic3,
				},
			},
		},
	}
	resp := &productpb.CreateTreatmentResponse{}
	err := suite.hdl.CreateTreatment(ctx, req, resp)
	suite.Assert().NoError(err)
}

func (suite *CreateTreatmentTestSuite) TearDownSuite() {}

func TestCreateTreatmentTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTreatmentTestSuite))
}
