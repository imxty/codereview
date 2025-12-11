package report

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SearchReportsTestSuite struct {
	suite.Suite
	hdl      *ReportAPIHandler
	customer *customermock.CustomerAPIService
}

func (suite *SearchReportsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	customer := &customermock.CustomerAPIService{}
	suite.customer = customer
	suite.hdl = NewReportAPIHandler(s, nil, customer, nil, nil, nil, "", "")
}

func (suite *SearchReportsTestSuite) SetupTest() {
	suite.customer.On("SearchCustomerIDByCustomerInfo", mock.Anything, mock.Anything).Return(
		&customerpb.SearchCustomerIDByCustomerInfoResponse{
			CustomerId: customerId,
		}, nil).After(utils.RpcLatency())
}

// TestSearchReports
func (suite *SearchReportsTestSuite) TestSearchReports() {
	ctx := context.Background()

	req := &reportpb.SearchReportsRequest{
		OrganizationId: organizationId,
		CustomerType:   reportpb.CustomerType_CUSTOMER_TYPE_CUSTOMER,
		CustomerName:   customerName,
		CustomerPhone:  phone,
		StartTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		EndTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		DirtyDialectics: []string{},
		TenantName:      tenantName,
		TenantId:        tenantId,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.SearchReportsResponse{}
	err := suite.hdl.SearchReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SearchReportsTestSuite) TestOrganizationIdIsNull() {
	ctx := context.Background()
	req := &reportpb.SearchReportsRequest{
		OrganizationId: organizationIdIsNull,
		CustomerType:   reportpb.CustomerType_CUSTOMER_TYPE_CUSTOMER,
		CustomerName:   customerName,
		CustomerPhone:  phone,
		StartTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		EndTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		DirtyDialectics: []string{},
		TenantName:      tenantName,
		TenantId:        tenantId,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.SearchReportsResponse{}
	err := suite.hdl.SearchReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantNameIsNull
func (suite *SearchReportsTestSuite) TestTenantNameIsNull() {
	ctx := context.Background()
	req := &reportpb.SearchReportsRequest{
		OrganizationId: organizationId,
		CustomerType:   reportpb.CustomerType_CUSTOMER_TYPE_CUSTOMER,
		CustomerName:   customerName,
		CustomerPhone:  phone,
		StartTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		EndTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		DirtyDialectics: []string{},
		TenantName:      tenantNameIsNull,
		TenantId:        tenantId,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.SearchReportsResponse{}
	err := suite.hdl.SearchReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *SearchReportsTestSuite) TearDownSuite() {
}

func TestSearchReportsTestSuite(t *testing.T) {
	suite.Run(t, new(SearchReportsTestSuite))
}
