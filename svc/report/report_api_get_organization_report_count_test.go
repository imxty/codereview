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

type GetOrganizationReportCountTestSuite struct {
	suite.Suite
	hdl *ReportAPIHandler
}

func (suite *GetOrganizationReportCountTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, nil, nil, "", "")
}

// TestGetOrganizationReportCount  获取公开分享的报告
func (suite *GetOrganizationReportCountTestSuite) TestGetOrganizationReportCount() {
	ctx := context.Background()
	req := &reportpb.GetOrganizationReportCountRequest{
		OrganizationId: organizationId,
	}
	resp := &reportpb.GetOrganizationReportCountResponse{}
	err := suite.hdl.GetOrganizationReportCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(0), resp.TodayCustomerMeasurementCount)
	suite.Assert().Equal(int32(0), resp.TodayTempCustomerMeasurementCount)
	suite.Assert().Equal(int32(2), resp.CustomerMeasurementYearCount)
	suite.Assert().Equal(int32(0), resp.TempCustomerMeasurementYearCount)

}

// TestOrganizationIdIsNull
func (suite *GetOrganizationReportCountTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetOrganizationReportCountRequest{
		OrganizationId: organizationIdIsNull,
	}
	resp := &reportpb.GetOrganizationReportCountResponse{}
	err := suite.hdl.GetOrganizationReportCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationReportCountTestSuite) TearDownSuite() {
}

func TestGetOrganizationReportCountTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationReportCountTestSuite))
}
