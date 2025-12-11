package report

import (
	"context"
	"testing"
	"time"

	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	usermock "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ListReportsTestSuite struct {
	suite.Suite
	hdl  *ReportAPIHandler
	user *usermock.UserAPIService
}

func (suite *ListReportsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	user := &usermock.UserAPIService{}
	suite.user = user
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, user, nil, "", "")
}

// TestListReports  查看常客报告
func (suite *ListReportsTestSuite) TestListReports() {
	ctx := context.Background()
	staffs := make(map[string]*userpb.Staff)
	staffs[staffId] = &userpb.Staff{
		StaffId: staffId,
	}

	suite.user.On("GetSharingReportSwitchStatus", mock.Anything, mock.Anything).Return(
		&userpb.GetSharingReportSwitchStatusResponse{
			SharingReportSwitchStatus: true,
		}, nil).After(utils.RpcLatency())
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListReportsRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
		StartTime:  startTime,
		EndTime:    endTime,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListReportsResponse{}
	err := suite.hdl.ListReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTempCustomerListReports  查看散客报告
func (suite *ListReportsTestSuite) TestTempCustomerListReports() {
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListReportsRequest{
		TenantId:   tenantId,
		CustomerId: customerIdIsNull,
		StartTime:  startTime,
		EndTime:    endTime,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListReportsResponse{}
	err := suite.hdl.ListReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *ListReportsTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListReportsRequest{
		TenantId:   tenantIdIsNull,
		CustomerId: customerId,
		StartTime:  startTime,
		EndTime:    endTime,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListReportsResponse{}
	err := suite.hdl.ListReports(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTimeIsErr
func (suite *ListReportsTestSuite) TestTimeIsErr() {
	t := suite.T()
	ctx := context.Background()

	req := &reportpb.ListReportsRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
		StartTime:  nil,
		EndTime:    nil,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListReportsResponse{}
	err := suite.hdl.ListReports(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPaginationSizeIsNegative 分页size小于0

func (suite *ListReportsTestSuite) TestPaginationSizeIsNegative() {
	t := suite.T()
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListReportsRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
		StartTime:  endTime,
		EndTime:    startTime,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   sizeIsNegative,
		},
	}
	resp := &reportpb.ListReportsResponse{}
	err := suite.hdl.ListReports(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.DataAccessFailed, err)
}

func (suite *ListReportsTestSuite) TearDownSuite() {
}

func TestListReportsTestSuite(t *testing.T) {
	suite.Run(t, new(ListReportsTestSuite))
}
