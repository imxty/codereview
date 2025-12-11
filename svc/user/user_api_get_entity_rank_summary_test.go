package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	reportmock "github.com/jinmukeji/huimaibao-service/svc/report/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetEntityRankSummaryTestSuite struct {
	suite.Suite
	hdl      *UserAPIHandler
	customer *customermock.CustomerAPIService
	calc     *reportmock.ReportAPIService
}

func (suite *GetEntityRankSummaryTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	customer := &customermock.CustomerAPIService{}
	suite.customer = customer
	calc := &reportmock.ReportAPIService{}
	suite.calc = calc
	suite.hdl = NewUserAPIHandler(s, nil, nil, customer, calc, nil, nil)
}

func (suite *GetEntityRankSummaryTestSuite) TestGetEntityRankSummary() {
	ctx := context.Background()
	suite.customer.On("GetCustomerCount", mock.Anything, mock.Anything).Return(
		&customerpb.GetCustomerCountResponse{
			TodayAddedCustomerCount: 2,
			CustomerTotalCount:      3,
		}, nil).After(utils.RpcLatency())
	suite.calc.On("GetTenantReportCount", mock.Anything, mock.Anything).Return(
		&reportpb.GetTenantReportCountResponse{
			TodayCustomerMeasurementCount:     2,
			TempCustomerMeasurementYearCount:  3,
			TodayTempCustomerMeasurementCount: 1,
			CustomerMeasurementYearCount:      3,
		}, nil).After(utils.RpcLatency())
	req := &userpb.GetEntityRankSummaryRequest{
		TenantId: tenantId,
	}
	resp := &userpb.GetEntityRankSummaryResponse{}
	err := suite.hdl.GetEntityRankSummary(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *GetEntityRankSummaryTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetEntityRankSummaryRequest{
		TenantId: tenantIdIsNull}
	resp := &userpb.GetEntityRankSummaryResponse{}
	err := suite.hdl.GetEntityRankSummary(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetEntityRankSummaryTestSuite) TearDownSuite() {

}

func TestGetEntityRankSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(GetEntityRankSummaryTestSuite))
}
