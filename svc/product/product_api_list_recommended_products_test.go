package product

import (
	"context"
	"testing"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/product"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	s3mock "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ListRecommendedProductsTestSuite struct {
	suite.Suite
	hdl  *ProductAPIHandler
	s3   *s3mock.FileStore
	user *s3mock.UserAPIService
}

func (suite *ListRecommendedProductsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, productFile)
	s := store.NewProductStore(conn)
	s3 := &s3mock.FileStore{}
	suite.s3 = s3
	user := &s3mock.UserAPIService{}
	suite.user = user
	suite.hdl = NewProductAPIHandler(s, s3, nil, user)
}

// TestReportIsFirstGenerated 报告初次生成.
func (suite *ListRecommendedProductsTestSuite) TestReportIsFirstGenerated() {
	ctx := context.Background()

	suite.user.On("GetTenantEntity", mock.Anything, mock.Anything).Return(
		&userpb.GetTenantEntityResponse{
			Entity: &userpb.Entity{
				OrganizationId: organizationId1,
			},
		}, nil).After(utils.RpcLatency())

	req := &productpb.ListRecommendedProductsRequest{
		TenantId:             tenantId1,
		ReportId:             reportId,
		TenantTreatmentRevId: treatmentRevId,
		IsFirstReported:      true,
		DirtyDialecticReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey1,
				Score:      4,
			},
			{
				SymptomKey: symptomKey2,
				Score:      4,
			},
		},
		HighRiskyDiseaseReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey3,
				Score:      5,
			},
		},
		MediumRiskyDiseaseReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey4,
				Score:      5,
			},
			{
				SymptomKey: symptomKey5,
				Score:      5,
			},
		},
		PhysicalTherapyReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey6,
				Score:      2,
			},
		},
	}
	resp := &productpb.ListRecommendedProductsResponse{}

	err := suite.hdl.ListRecommendedProducts(ctx, req, resp)
	suite.Assert().True(true, resp.GetAreLatestTreatments())
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestReportIsNotFirstGenerated 报告不是初次生成.
func (suite *ListRecommendedProductsTestSuite) TestReportIsNotFirstGenerated() {
	ctx := context.Background()
	req := &productpb.ListRecommendedProductsRequest{
		TenantId:             tenantId1,
		ReportId:             reportId,
		TenantTreatmentRevId: treatmentRevId,
		IsFirstReported:      false,
		DirtyDialecticReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey1,
				Score:      4,
			},
			{
				SymptomKey: symptomKey2,
				Score:      4,
			},
		},
		HighRiskyDiseaseReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey3,
				Score:      5,
			},
		},
		MediumRiskyDiseaseReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey4,
				Score:      5,
			},
			{
				SymptomKey: symptomKey5,
				Score:      5,
			},
		},
		PhysicalTherapyReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey6,
				Score:      2,
			},
		},
	}
	resp := &productpb.ListRecommendedProductsResponse{}

	err := suite.hdl.ListRecommendedProducts(ctx, req, resp)
	suite.T().Log(resp.GetTenantTreatmentRevId())
	suite.Assert().True(true, resp.GetAreLatestTreatments())
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestReportWithNoHighRiskyDiseaseReportResult.
func (suite *ListRecommendedProductsTestSuite) TestReportWithNoHighRiskyDiseaseReportResult() {
	ctx := context.Background()
	req := &productpb.ListRecommendedProductsRequest{
		TenantId:             tenantId1,
		ReportId:             reportId,
		TenantTreatmentRevId: treatmentRevId,
		IsFirstReported:      false,
		DirtyDialecticReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey1,
				Score:      4,
			},
			{
				SymptomKey: symptomKey2,
				Score:      4,
			},
		},
		MediumRiskyDiseaseReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey4,
				Score:      5,
			},
			{
				SymptomKey: symptomKey5,
				Score:      5,
			},
		},
		PhysicalTherapyReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey6,
				Score:      2,
			},
		},
	}
	resp := &productpb.ListRecommendedProductsResponse{}

	err := suite.hdl.ListRecommendedProducts(ctx, req, resp)
	suite.T().Log(resp.GetTenantTreatmentRevId())
	suite.Assert().True(true, resp.GetAreLatestTreatments())
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestReportIsFirstGeneratedWithProductRecommendationStatusIsFalse
func (suite *ListRecommendedProductsTestSuite) TestReportIsFirstGeneratedWithProductRecommendationStatusIsFalse() {
	ctx := context.Background()
	req := &productpb.ListRecommendedProductsRequest{
		TenantId:             tenantId1,
		ReportId:             reportId,
		TenantTreatmentRevId: "",
		IsFirstReported:      true,
		DirtyDialecticReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey1,
				Score:      4,
			},
			{
				SymptomKey: symptomKey2,
				Score:      4,
			},
		},
		HighRiskyDiseaseReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey3,
				Score:      5,
			},
		},
		MediumRiskyDiseaseReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey4,
				Score:      5,
			},
			{
				SymptomKey: symptomKey5,
				Score:      5,
			},
		},
		PhysicalTherapyReportResult: []*productpb.ReportResult{
			{
				SymptomKey: symptomKey6,
				Score:      2,
			},
		},
	}
	resp := &productpb.ListRecommendedProductsResponse{}

	err := suite.hdl.ListRecommendedProducts(ctx, req, resp)
	suite.T().Log(resp.GetTenantTreatmentRevId())
	suite.Assert().True(true, resp.GetAreLatestTreatments())
	suite.Assert().NoError(err)
	suite.Assert().Equal(1, len(resp.GetProducts()))
}

func (suite *ListRecommendedProductsTestSuite) TearDownSuite() {

}

func TestListRecommendedProductsTestSuite(t *testing.T) {
	suite.Run(t, new(ListRecommendedProductsTestSuite))
}
