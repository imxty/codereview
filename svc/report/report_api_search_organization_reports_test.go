package report

import (
	"context"
	"testing"

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

type SearchOrganizationReportsTestSuite struct {
	suite.Suite
	hdl  *ReportAPIHandler
	user *usermock.UserAPIService
}

func (suite *SearchOrganizationReportsTestSuite) SetupSuite() {
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
func (suite *SearchOrganizationReportsTestSuite) SetupTest() {
	suite.user.On("BatchGetOrganizationNamesByIDs", mock.Anything, mock.Anything).Return(
		&userpb.BatchGetOrganizationNamesByIDsResponse{
			Organizations: make(map[string]string),
		}, nil).After(utils.RpcLatency())
	suite.user.On("BatchGetTenantNamesByIDs", mock.Anything, mock.Anything).Return(
		&userpb.BatchGetTenantNamesByIDsResponse{
			TenantNames: make(map[string]string),
		}, nil).After(utils.RpcLatency())
}

// TestSearchOrganizationReports
func (suite *SearchOrganizationReportsTestSuite) TestSearchOrganizationReports() {
	ctx := context.Background()
	req := &reportpb.SearchOrganizationReportsRequest{
		OrganizationId: organizationId,
		TenantName:     tenantName,
		StartTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		EndTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		OperatorId: staffId,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.SearchOrganizationReportsResponse{}
	err := suite.hdl.SearchOrganizationReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SearchOrganizationReportsTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.SearchOrganizationReportsRequest{
		OrganizationId: organizationIdIsNull,
		TenantName:     tenantName,
		StartTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		EndTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		OperatorId: staffId,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.SearchOrganizationReportsResponse{}
	err := suite.hdl.SearchOrganizationReports(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantNameIsNull
func (suite *SearchOrganizationReportsTestSuite) TestTenantNameIsNull() {
	ctx := context.Background()
	req := &reportpb.SearchOrganizationReportsRequest{
		OrganizationId: organizationId,
		TenantName:     tenantNameIsNull,
		StartTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		EndTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		OperatorId: staffId,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.SearchOrganizationReportsResponse{}
	err := suite.hdl.SearchOrganizationReports(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}
func (suite *SearchOrganizationReportsTestSuite) TearDownSuite() {
}

func TestSearchOrganizationReportsTestSuite(t *testing.T) {
	suite.Run(t, new(SearchOrganizationReportsTestSuite))
}
