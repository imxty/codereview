package report

import (
	"context"
	"testing"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	productmock "github.com/jinmukeji/huimaibao-service/svc/product/mock"
	reportmock "github.com/jinmukeji/huimaibao-service/svc/report/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	usermock "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	calcpb "github.com/jinmukeji/platform-proto/gen/go/platform/report/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GetPublicSharedReportTestSuite struct {
	suite.Suite
	hdl     *ReportAPIHandler
	user    *usermock.UserAPIService
	product *productmock.ProductAPIService
	calc    *reportmock.ReportAPIClient
}

func (suite *GetPublicSharedReportTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	user := &usermock.UserAPIService{}
	suite.user = user
	product := &productmock.ProductAPIService{}
	suite.product = product
	calc := &reportmock.ReportAPIClient{}
	suite.calc = calc
	suite.hdl = NewReportAPIHandler(s, calc, nil, nil, user, product, "", "")
}

// TestGetPublicSharedReport  获取公开分享的报告
func (suite *GetPublicSharedReportTestSuite) TestGetPublicSharedReport() {
	ctx := context.Background()
	suite.calc.On("GetRawData", mock.Anything, mock.Anything).Return(
		&calcpb.GetRawDataResponse{}, nil).After(utils.RpcLatency())

	suite.user.On("GetTenant", mock.Anything, mock.Anything).Return(
		&userpb.GetTenantResponse{
			Tenant: &userpb.TenantEntity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())
	suite.product.On("ListRecommendedProducts", mock.Anything, mock.Anything).Return(
		&productpb.ListRecommendedProductsResponse{
			Products: []*productpb.Product{
				&productpb.Product{
					ProductId: productId,
				},
			},
		}, nil).After(utils.RpcLatency())
	req := &reportpb.GetPublicSharedReportRequest{
		TenantId:     tenantId,
		Token:        token,
		LanguageCode: languageCode,
	}
	resp := &reportpb.GetPublicSharedReportResponse{}
	err := suite.hdl.GetPublicSharedReport(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(reportId, resp.Report.ReportId)
}

// TestTenantIdIsNull
func (suite *GetPublicSharedReportTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetPublicSharedReportRequest{
		TenantId:     tenantIdIsNull,
		Token:        token,
		LanguageCode: languageCode,
	}
	resp := &reportpb.GetPublicSharedReportResponse{}
	err := suite.hdl.GetPublicSharedReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTokenIsNull
func (suite *GetPublicSharedReportTestSuite) TestTokenIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetPublicSharedReportRequest{
		TenantId:     tenantId,
		Token:        tokenIsNull,
		LanguageCode: languageCode,
	}
	resp := &reportpb.GetPublicSharedReportResponse{}
	err := suite.hdl.GetPublicSharedReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetPublicSharedReportTestSuite) TearDownSuite() {
}

func TestGetPublicSharedReportTestSuite(t *testing.T) {
	suite.Run(t, new(GetPublicSharedReportTestSuite))
}
