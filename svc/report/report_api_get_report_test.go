package report

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
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

type GetReportTestSuite struct {
	suite.Suite
	hdl      *ReportAPIHandler
	user     *usermock.UserAPIService
	customer *customermock.CustomerAPIService
	calc     *reportmock.ReportAPIClient
	product  *productmock.ProductAPIService
}

func (suite *GetReportTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	user := &usermock.UserAPIService{}
	suite.user = user
	customer := &customermock.CustomerAPIService{}
	suite.customer = customer
	calc := &reportmock.ReportAPIClient{}
	suite.calc = calc
	product := &productmock.ProductAPIService{}
	suite.product = product
	suite.hdl = NewReportAPIHandler(s, calc, customer, nil, user, product, "", "")
}

// TestGetReport  获取报告
func (suite *GetReportTestSuite) TestGetReport() {
	ctx := context.Background()
	suite.calc.On("GetRawData", mock.Anything, mock.Anything).Return(
		&calcpb.GetRawDataResponse{}, nil).After(utils.RpcLatency())
	suite.user.On("GetTenant", mock.Anything, mock.Anything).Return(
		&userpb.GetTenantResponse{
			Tenant: &userpb.TenantEntity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())
	suite.customer.On("GetCustomer", mock.Anything, mock.Anything).Return(
		&customerpb.GetCustomerResponse{
			Customer: &customerpb.Customer{
				CustomerId: customerId,
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
	answers := make(map[string]*reportpb.AnswerList)
	answers["stress_state_judgment"] = &reportpb.AnswerList{
		Answers: []*reportpb.Answer{
			&reportpb.Answer{
				QuestionKey: "Q0021",
				AnswerKeys:  []string{"QC0082"},
			}, &reportpb.Answer{
				QuestionKey: "Q0033",
				AnswerKeys:  []string{"QC0121"},
			},
		},
	}
	req := &reportpb.GetReportRequest{
		TenantId:      tenantId,
		ReportId:      reportId,
		LanguageCode:  languageCode,
		ModuleAnswers: answers,
	}
	resp := &reportpb.GetReportResponse{}
	err := suite.hdl.GetReport(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().True(true, resp.IsCompleteReport)
}

// TestTenantIdIsNull
func (suite *GetReportTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	answers := make(map[string]*reportpb.AnswerList)
	answers["stress_state_judgment"] = &reportpb.AnswerList{
		Answers: []*reportpb.Answer{
			&reportpb.Answer{
				QuestionKey: "Q0021",
				AnswerKeys:  []string{"QC0082"},
			}, &reportpb.Answer{
				QuestionKey: "Q0033",
				AnswerKeys:  []string{"QC0121"},
			},
		},
	}
	req := &reportpb.GetReportRequest{
		TenantId:      tenantIdIsNull,
		ReportId:      reportId,
		LanguageCode:  languageCode,
		ModuleAnswers: answers,
	}
	resp := &reportpb.GetReportResponse{}
	err := suite.hdl.GetReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReportIdIsNull
func (suite *GetReportTestSuite) TestReportIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	answers := make(map[string]*reportpb.AnswerList)
	answers["stress_state_judgment"] = &reportpb.AnswerList{
		Answers: []*reportpb.Answer{
			&reportpb.Answer{
				QuestionKey: "Q0021",
				AnswerKeys:  []string{"QC0082"},
			}, &reportpb.Answer{
				QuestionKey: "Q0033",
				AnswerKeys:  []string{"QC0121"},
			},
		},
	}
	req := &reportpb.GetReportRequest{
		TenantId:      tenantId,
		ReportId:      reportIdIsNull,
		LanguageCode:  languageCode,
		ModuleAnswers: answers,
	}
	resp := &reportpb.GetReportResponse{}
	err := suite.hdl.GetReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReportIdIsNotExist
func (suite *GetReportTestSuite) TestReportIdIsNotExist() {
	t := suite.T()
	ctx := context.Background()
	answers := make(map[string]*reportpb.AnswerList)
	answers["stress_state_judgment"] = &reportpb.AnswerList{
		Answers: []*reportpb.Answer{
			&reportpb.Answer{
				QuestionKey: "Q0021",
				AnswerKeys:  []string{"QC0082"},
			}, &reportpb.Answer{
				QuestionKey: "Q0033",
				AnswerKeys:  []string{"QC0121"},
			},
		},
	}
	req := &reportpb.GetReportRequest{
		TenantId:      tenantId,
		ReportId:      reportIdIsNotExist,
		LanguageCode:  languageCode,
		ModuleAnswers: answers,
	}
	resp := &reportpb.GetReportResponse{}
	err := suite.hdl.GetReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrReportNotFound, err)
}

func (suite *GetReportTestSuite) TearDownSuite() {
}

func TestGetReportTestSuite(t *testing.T) {
	suite.Run(t, new(GetReportTestSuite))
}
