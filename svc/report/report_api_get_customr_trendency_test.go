package report

import (
	"context"
	"testing"

	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GetCustomerTrendencyTestSuite struct {
	suite.Suite
	hdl *ReportAPIHandler
}

func (suite *GetCustomerTrendencyTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, nil, nil, "", "")
}

// TestGetCustomerTrendency
func (suite *GetCustomerTrendencyTestSuite) TestGetCustomerTrendency() {
	ctx := context.Background()
	req := &reportpb.GetCustomerTrendencyRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
		Days:       2,
	}
	resp := &reportpb.GetCustomerTrendencyResponse{}
	err := suite.hdl.GetCustomerTrendency(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)

}

// TestTenantIdIsNull
func (suite *GetCustomerTrendencyTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetCustomerTrendencyRequest{
		TenantId:   tenantIdIsNull,
		CustomerId: customerId,
		Days:       2,
	}
	resp := &reportpb.GetCustomerTrendencyResponse{}
	err := suite.hdl.GetCustomerTrendency(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdIsNull
func (suite *GetCustomerTrendencyTestSuite) TestCustomerIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetCustomerTrendencyRequest{
		TenantId:   tenantId,
		CustomerId: customerIdIsNull,
		Days:       2,
	}
	resp := &reportpb.GetCustomerTrendencyResponse{}
	err := suite.hdl.GetCustomerTrendency(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetCustomerTrendencyTestSuite) TearDownSuite() {
}

func TestGetCustomerTrendencyTestSuite(t *testing.T) {
	suite.Run(t, new(GetCustomerTrendencyTestSuite))
}
