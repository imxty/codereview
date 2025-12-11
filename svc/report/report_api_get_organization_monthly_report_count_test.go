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

type GetOrganizationMonthlyReportCountTestSuite struct {
	suite.Suite
	hdl *ReportAPIHandler
}

func (suite *GetOrganizationMonthlyReportCountTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, nil, nil, "", "")
}

// TestGetOrganizationMonthlyReportCount  (当月6，上月1，去年同月1条)
func (suite *GetOrganizationMonthlyReportCountTestSuite) TestGetOrganizationMonthlyReportCount() {
	ctx := context.Background()
	req := &reportpb.GetOrganizationMonthlyReportCountRequest{
		OrganizationId: organizationId,
		MeasurementTime: &reportpb.Date{
			Year:  newYear,
			Month: month,
			Day:   day,
		},
	}
	resp := &reportpb.GetOrganizationMonthlyReportCountResponse{}
	err := suite.hdl.GetOrganizationMonthlyReportCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(float64(5), resp.MonthOnMonth)
	suite.Assert().Equal(float64(5), resp.YearOnYear)
	suite.Assert().Equal(int32(6), resp.MonthlyCustomerMeasurementCount)

}

// TestGetOrganizationMonthlyReportCountMonthIsNull
func (suite *GetOrganizationMonthlyReportCountTestSuite) TestGetOrganizationMonthlyReportCountMonthIsNull() {
	ctx := context.Background()
	req := &reportpb.GetOrganizationMonthlyReportCountRequest{
		OrganizationId: organizationId,
		MeasurementTime: &reportpb.Date{
			Year:  newYear,
			Month: month2,
			Day:   day,
		},
	}
	resp := &reportpb.GetOrganizationMonthlyReportCountResponse{}
	err := suite.hdl.GetOrganizationMonthlyReportCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(float64(0), resp.MonthOnMonth)
	suite.Assert().Equal(float64(0), resp.YearOnYear)
	suite.Assert().Equal(int32(1), resp.MonthlyCustomerMeasurementCount)

}
func (suite *GetOrganizationMonthlyReportCountTestSuite) TestGetOrganizationMonthlyReportCountErr() {
	ctx := context.Background()
	req := &reportpb.GetOrganizationMonthlyReportCountRequest{
		OrganizationId: organizationId,
		MeasurementTime: &reportpb.Date{
			Year:  newYear,
			Month: month1,
			Day:   day,
		},
	}
	resp := &reportpb.GetOrganizationMonthlyReportCountResponse{}
	err := suite.hdl.GetOrganizationMonthlyReportCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(float64(2.147483647e+09), resp.MonthOnMonth)
	suite.Assert().Equal(float64(2.147483647e+09), resp.YearOnYear)
	suite.Assert().Equal(int32(0), resp.MonthlyCustomerMeasurementCount)

}

// TestOrganizationIdIsNull
func (suite *GetOrganizationMonthlyReportCountTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetOrganizationMonthlyReportCountRequest{
		OrganizationId: organizationIdIsNull,
		MeasurementTime: &reportpb.Date{
			Year:  year,
			Month: month,
			Day:   day,
		},
	}
	resp := &reportpb.GetOrganizationMonthlyReportCountResponse{}
	err := suite.hdl.GetOrganizationMonthlyReportCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestMeasurementTimeIsNull
func (suite *GetOrganizationMonthlyReportCountTestSuite) TestMeasurementTimeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetOrganizationMonthlyReportCountRequest{
		OrganizationId:  organizationId,
		MeasurementTime: nil,
	}
	resp := &reportpb.GetOrganizationMonthlyReportCountResponse{}
	err := suite.hdl.GetOrganizationMonthlyReportCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestMeasurementTimeIsErr
func (suite *GetOrganizationMonthlyReportCountTestSuite) TestMeasurementTimeIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetOrganizationMonthlyReportCountRequest{
		OrganizationId:  organizationId,
		MeasurementTime: &reportpb.Date{},
	}
	resp := &reportpb.GetOrganizationMonthlyReportCountResponse{}
	err := suite.hdl.GetOrganizationMonthlyReportCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationMonthlyReportCountTestSuite) TearDownSuite() {
}

func TestGetOrganizationMonthlyReportCountTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationMonthlyReportCountTestSuite))
}
