package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	reportmock "github.com/jinmukeji/huimaibao-service/svc/report/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetEntityMeasurementSummaryTestSuite struct {
	suite.Suite
	hdl  *UserAPIHandler
	calc *reportmock.ReportAPIService
}

func (suite *GetEntityMeasurementSummaryTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	calc := &reportmock.ReportAPIService{}
	suite.calc = calc
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, calc, nil, nil)
}

func (suite *GetEntityMeasurementSummaryTestSuite) TestGetEntityMeasurementSummary() {
	ctx := context.Background()
	suite.calc.On("GetTenantReportCount", mock.Anything, mock.Anything).Return(
		&reportpb.GetTenantReportCountResponse{
			TodayCustomerMeasurementCount:     1,
			TodayTempCustomerMeasurementCount: 2,
			TempCustomerMeasurementYearCount:  10,
			CustomerMeasurementYearCount:      10,
		}, nil).After(utils.RpcLatency())
	req := &userpb.GetEntityMeasurementSummaryRequest{
		TenantId: tenantId,
	}
	resp := &userpb.GetEntityMeasurementSummaryResponse{}
	err := suite.hdl.GetEntityMeasurementSummary(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *GetEntityMeasurementSummaryTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetEntityMeasurementSummaryRequest{
		TenantId: tenantIdIsNull}
	resp := &userpb.GetEntityMeasurementSummaryResponse{}
	err := suite.hdl.GetEntityMeasurementSummary(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetEntityMeasurementSummaryTestSuite) TearDownSuite() {

}

func TestGetEntityMeasurementSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(GetEntityMeasurementSummaryTestSuite))
}
