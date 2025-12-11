package report

import (
	"context"
	"testing"
	"time"

	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	usermock "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ListTempCustomerReportsTestSuite struct {
	suite.Suite
	hdl *ReportAPIHandler
}

func (suite *ListTempCustomerReportsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	user := &usermock.UserAPIService{}
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, user, nil, "", "")
}

// TestListTempCustomerReports  查看散客报告
func (suite *ListTempCustomerReportsTestSuite) TestListTempCustomerReports() {
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListTempCustomerReportsRequest{
		TenantId: tenantId,
		ReportId: reportId,
		TimeRange: &reportpb.TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListTempCustomerReportsResponse{}
	err := suite.hdl.ListTempCustomerReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(1), resp.TotalCount)
}

// TestReportIsNull
func (suite *ListTempCustomerReportsTestSuite) TestReportIsNull() {
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListTempCustomerReportsRequest{
		TenantId: tenantId,
		ReportId: reportIdIsNull,
		TimeRange: &reportpb.TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListTempCustomerReportsResponse{}
	err := suite.hdl.ListTempCustomerReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(0), resp.TotalCount)
}

// TestTenantIdIsNull
func (suite *ListTempCustomerReportsTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListTempCustomerReportsRequest{
		TenantId: tenantIdIsNull,
		ReportId: reportIdIsNull,
		TimeRange: &reportpb.TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListTempCustomerReportsResponse{}
	err := suite.hdl.ListTempCustomerReports(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTimeIsErr
func (suite *ListTempCustomerReportsTestSuite) TestTimeIsErr() {
	ctx := context.Background()

	req := &reportpb.ListTempCustomerReportsRequest{
		TenantId: tenantId,
		ReportId: reportId,
		TimeRange: &reportpb.TimeRange{
			StartTime: nil,
			EndTime:   nil,
		},
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListTempCustomerReportsResponse{}
	err := suite.hdl.ListTempCustomerReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPaginationSizeIsNegative 分页size小于0

func (suite *ListTempCustomerReportsTestSuite) TestPaginationSizeIsNegative() {
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListTempCustomerReportsRequest{
		TenantId: tenantId,
		ReportId: reportId,
		TimeRange: &reportpb.TimeRange{
			StartTime: endTime,
			EndTime:   startTime,
		},
		Pagination: &reportpb.Pagination{
			Offset: offsetIsNegative,
			Size:   sizeIsNegative,
		},
	}
	resp := &reportpb.ListTempCustomerReportsResponse{}
	err := suite.hdl.ListTempCustomerReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(1), resp.TotalCount)
}

func (suite *ListTempCustomerReportsTestSuite) TearDownSuite() {
}

func TestListTempCustomerReportsTestSuite(t *testing.T) {
	suite.Run(t, new(ListTempCustomerReportsTestSuite))
}
