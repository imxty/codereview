package report

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	reportmock "github.com/jinmukeji/huimaibao-service/svc/report/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	calcpb "github.com/jinmukeji/platform-proto/gen/go/platform/report/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GetWeeklyReportTestSuite struct {
	suite.Suite
	hdl      *ReportAPIHandler
	customer *customermock.CustomerAPIService
	calc     *reportmock.ReportAPIClient
}

func (suite *GetWeeklyReportTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	customer := &customermock.CustomerAPIService{}
	suite.customer = customer
	calc := &reportmock.ReportAPIClient{}
	suite.calc = calc
	suite.hdl = NewReportAPIHandler(s, calc, customer, nil, nil, nil, "", "")
}

// TestGetWeeklyReport
func (suite *GetWeeklyReportTestSuite) TestGetWeeklyReport() {
	ctx := context.Background()
	suite.customer.On("GetCustomer", mock.Anything, mock.Anything).Return(
		&customerpb.GetCustomerResponse{
			Customer: &customerpb.Customer{
				CustomerId: customerId,
			},
		}, nil).After(utils.RpcLatency())
	suite.calc.On("GetWeeklyReport", mock.Anything, mock.Anything).Return(
		&calcpb.GetWeeklyReportResponse{
			Report: &calcpb.ReportContent{
				ReportId: reportId,
				ModuleResults: []*calcpb.ReportModuleResult{
					&calcpb.ReportModuleResult{
						ModuleName: "DirtyDialectic",
						Enabled:    true,
						Result:     &anypb.Any{},
					},
					&calcpb.ReportModuleResult{
						ModuleName: "MeasurementJudgment",
						Enabled:    true,
						Result:     &anypb.Any{},
					},
				},
			},
		}, nil).After(utils.RpcLatency())
	req := &reportpb.GetWeeklyReportRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
	}
	resp := &reportpb.GetWeeklyReportResponse{}
	err := suite.hdl.GetWeeklyReport(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *GetWeeklyReportTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetWeeklyReportRequest{
		TenantId:   tenantIdIsNull,
		CustomerId: customerId,
	}
	resp := &reportpb.GetWeeklyReportResponse{}
	err := suite.hdl.GetWeeklyReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdIsNull
func (suite *GetWeeklyReportTestSuite) TestCustomerIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetWeeklyReportRequest{
		TenantId:   tenantId,
		CustomerId: customerIdIsNull,
	}
	resp := &reportpb.GetWeeklyReportResponse{}
	err := suite.hdl.GetWeeklyReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetWeeklyReportTestSuite) TearDownSuite() {
}

func TestGetWeeklyReportTestSuite(t *testing.T) {
	suite.Run(t, new(GetWeeklyReportTestSuite))
}
