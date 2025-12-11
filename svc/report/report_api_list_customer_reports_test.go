package report

import (
	"context"
	"testing"
	"time"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	usermock "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ListCustomerReportsTestSuite struct {
	suite.Suite
	hdl      *ReportAPIHandler
	customer *customermock.CustomerAPIService
}

func (suite *ListCustomerReportsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	user := &usermock.UserAPIService{}
	customer := &customermock.CustomerAPIService{}
	suite.customer = customer
	suite.hdl = NewReportAPIHandler(s, nil, customer, nil, user, nil, "", "")
}

// TestListReports  查看常客报告
func (suite *ListCustomerReportsTestSuite) TestListCustomerReports() {
	ctx := context.Background()
	suite.customer.On("GetCustomer", mock.Anything, mock.Anything).Return(
		&customerpb.GetCustomerResponse{
			Customer: &customerpb.Customer{
				CustomerId: customerId,
			},
		}, nil).After(utils.RpcLatency())
	startTime := timestamppb.New(time.Date(2023, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2024, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListCustomerReportsRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
		ReportId:   reportId1,
		TimeRange: &reportpb.TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListCustomerReportsResponse{}
	err := suite.hdl.ListCustomerReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestReportIsNull
func (suite *ListCustomerReportsTestSuite) TestReportIsNull() {
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListCustomerReportsRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
		ReportId:   reportIdIsNull,
		TimeRange: &reportpb.TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListCustomerReportsResponse{}
	err := suite.hdl.ListCustomerReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *ListCustomerReportsTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListCustomerReportsRequest{
		TenantId:   tenantIdIsNull,
		CustomerId: customerId,
		ReportId:   reportIdIsNull,
		TimeRange: &reportpb.TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListCustomerReportsResponse{}
	err := suite.hdl.ListCustomerReports(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdIsNull
func (suite *ListCustomerReportsTestSuite) TestCustomerIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListCustomerReportsRequest{
		TenantId:   tenantId,
		CustomerId: customerIdIsNull,
		ReportId:   reportId,
		TimeRange: &reportpb.TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListCustomerReportsResponse{}
	err := suite.hdl.ListCustomerReports(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdNotExist
func (suite *ListCustomerReportsTestSuite) TestCustomerIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2020, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2021, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListCustomerReportsRequest{
		TenantId:   tenantIdIsNotExist,
		CustomerId: customerIdIsNotExist,
		ReportId:   reportId,
		TimeRange: &reportpb.TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListCustomerReportsResponse{}
	err := suite.hdl.ListCustomerReports(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTimeIsErr
func (suite *ListCustomerReportsTestSuite) TestTimeIsErr() {
	t := suite.T()
	ctx := context.Background()

	req := &reportpb.ListCustomerReportsRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
		ReportId:   reportIdIsNull,
		TimeRange:  nil,
		Pagination: nil,
	}
	resp := &reportpb.ListCustomerReportsResponse{}
	err := suite.hdl.ListCustomerReports(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListCustomerReportsTestSuite) TearDownSuite() {
}

func TestListCustomerReportsTestSuite(t *testing.T) {
	suite.Run(t, new(ListCustomerReportsTestSuite))
}
