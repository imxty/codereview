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

type GetTenantReportCountTestSuite struct {
	suite.Suite
	hdl *ReportAPIHandler
}

func (suite *GetTenantReportCountTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, nil, nil, "", "")
}

// TestGetTenantReportCount  获取报告统计
func (suite *GetTenantReportCountTestSuite) TestGetTenantReportCount() {
	ctx := context.Background()
	req := &reportpb.GetTenantReportCountRequest{
		TenantId: tenantId,
		StaffIds: []string{staffId},
	}
	resp := &reportpb.GetTenantReportCountResponse{}
	err := suite.hdl.GetTenantReportCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(0), resp.TodayCustomerMeasurementCount)
	suite.Assert().Equal(int32(0), resp.TodayTempCustomerMeasurementCount)
	suite.Assert().Equal(int32(2), resp.CustomerMeasurementYearCount)
	suite.Assert().Equal(int32(0), resp.TempCustomerMeasurementYearCount)

}

// TestTenantIdIsNull
func (suite *GetTenantReportCountTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetTenantReportCountRequest{
		TenantId: tenantIdIsNull,
		StaffIds: []string{staffId},
	}
	resp := &reportpb.GetTenantReportCountResponse{}
	err := suite.hdl.GetTenantReportCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNoExist
func (suite *GetTenantReportCountTestSuite) TestTenantIdNoExist() {

	ctx := context.Background()
	req := &reportpb.GetTenantReportCountRequest{
		TenantId: tenantIdIsNotExist,
		StaffIds: []string{staffId},
	}
	resp := &reportpb.GetTenantReportCountResponse{}
	err := suite.hdl.GetTenantReportCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestStaffIdsIsNull
func (suite *GetTenantReportCountTestSuite) TestStaffIdsIsNull() {

	ctx := context.Background()
	req := &reportpb.GetTenantReportCountRequest{
		TenantId: tenantId,
		StaffIds: []string{staffIdIsNull},
	}
	resp := &reportpb.GetTenantReportCountResponse{}
	err := suite.hdl.GetTenantReportCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *GetTenantReportCountTestSuite) TearDownSuite() {
}

func TestGetTenantReportCountTestSuite(t *testing.T) {
	suite.Run(t, new(GetTenantReportCountTestSuite))
}
