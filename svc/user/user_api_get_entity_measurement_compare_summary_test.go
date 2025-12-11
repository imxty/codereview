package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	reportmock "github.com/jinmukeji/huimaibao-service/svc/report/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetEntityMeasurementCompareSummaryTestSuite struct {
	suite.Suite
	hdl  *UserAPIHandler
	calc *reportmock.ReportAPIService
}

func (suite *GetEntityMeasurementCompareSummaryTestSuite) SetupSuite() {
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

func (suite *GetEntityMeasurementCompareSummaryTestSuite) TestGetEntityMeasurementCompareSummary() {
	ctx := context.Background()
	suite.calc.On("GetTenantMonthlyReportCount", mock.Anything, mock.Anything).Return(
		&reportpb.GetTenantMonthlyReportCountResponse{
			MonthlyCustomerMeasurementCount: 1,
			MonthOnMonth:                    0.1,
			YearOnYear:                      0.1,
		}, nil).After(utils.RpcLatency())
	req := &userpb.GetEntityMeasurementCompareSummaryRequest{
		TenantId: tenantId,
		Date: &timestamppb.Timestamp{
			Seconds: 9,
			Nanos:   1,
		},
	}
	resp := &userpb.GetEntityMeasurementCompareSummaryResponse{}
	err := suite.hdl.GetEntityMeasurementCompareSummary(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *GetEntityMeasurementCompareSummaryTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetEntityMeasurementCompareSummaryRequest{
		TenantId: tenantIdIsNull,
		Date: &timestamppb.Timestamp{
			Seconds: 9,
			Nanos:   1,
		}}
	resp := &userpb.GetEntityMeasurementCompareSummaryResponse{}
	err := suite.hdl.GetEntityMeasurementCompareSummary(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetEntityMeasurementCompareSummaryTestSuite) TearDownSuite() {

}

func TestGetEntityMeasurementCompareSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(GetEntityMeasurementCompareSummaryTestSuite))
}
