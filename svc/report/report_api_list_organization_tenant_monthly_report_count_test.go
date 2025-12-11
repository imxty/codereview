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

type ListOrganizationTenantMonthlyReportCountTestSuite struct {
	suite.Suite
	hdl  *ReportAPIHandler
	user *usermock.UserAPIService
}

func (suite *ListOrganizationTenantMonthlyReportCountTestSuite) SetupSuite() {
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

// TestListOrganizationTenantMonthlyReportCount  查看常客报告
func (suite *ListOrganizationTenantMonthlyReportCountTestSuite) TestListOrganizationTenantMonthlyReportCount() {
	ctx := context.Background()

	suite.user.On("SearchTenantByNameAndOrganizationID", mock.Anything, mock.Anything).Return(
		&userpb.SearchTenantByNameAndOrganizationIDResponse{
			TenantIds: []string{tenantId},
		}, nil).After(utils.RpcLatency())
	suite.user.On("BatchGetTenantNamesByIDs", mock.Anything, mock.Anything).Return(
		&userpb.BatchGetTenantNamesByIDsResponse{
			TenantNames: make(map[string]string),
		}, nil).After(utils.RpcLatency())
	startTime := timestamppb.New(time.Date(2023, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2024, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListOrganizationTenantMonthlyReportCountRequest{
		OrganizationId: organizationId,
		TenantName:     tenantName,
		StartTime:      startTime,
		EndTime:        endTime,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListOrganizationTenantMonthlyReportCountResponse{}
	err := suite.hdl.ListOrganizationTenantMonthlyReportCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *ListOrganizationTenantMonthlyReportCountTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2023, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2024, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListOrganizationTenantMonthlyReportCountRequest{
		OrganizationId: organizationIdIsNull,
		TenantName:     tenantName,
		StartTime:      startTime,
		EndTime:        endTime,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListOrganizationTenantMonthlyReportCountResponse{}
	err := suite.hdl.ListOrganizationTenantMonthlyReportCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantNameIsNull
func (suite *ListOrganizationTenantMonthlyReportCountTestSuite) TestTenantNameIsNull() {
	t := suite.T()
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2023, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2024, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListOrganizationTenantMonthlyReportCountRequest{
		OrganizationId: organizationId,
		TenantName:     tenantNameIsNull,
		StartTime:      startTime,
		EndTime:        endTime,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListOrganizationTenantMonthlyReportCountResponse{}
	err := suite.hdl.ListOrganizationTenantMonthlyReportCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTimeIsErr
func (suite *ListOrganizationTenantMonthlyReportCountTestSuite) TestTimeIsErr() {
	t := suite.T()
	ctx := context.Background()

	req := &reportpb.ListOrganizationTenantMonthlyReportCountRequest{
		OrganizationId: organizationId,
		TenantName:     tenantName,
		StartTime:      nil,
		EndTime:        nil,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &reportpb.ListOrganizationTenantMonthlyReportCountResponse{}
	err := suite.hdl.ListOrganizationTenantMonthlyReportCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPaginationSizeIsNegative 分页size小于0

func (suite *ListOrganizationTenantMonthlyReportCountTestSuite) TestPaginationSizeIsNegative() {
	t := suite.T()
	ctx := context.Background()
	startTime := timestamppb.New(time.Date(2023, 5, 11, 1, 1, 1, 0, time.UTC))
	endTime := timestamppb.New(time.Date(2024, 11, 11, 1, 1, 1, 0, time.UTC))

	req := &reportpb.ListOrganizationTenantMonthlyReportCountRequest{
		OrganizationId: organizationId,
		TenantName:     tenantName,
		StartTime:      endTime,
		EndTime:        startTime,
		Pagination: &reportpb.Pagination{
			Offset: offset,
			Size:   sizeIsNegative,
		},
	}
	resp := &reportpb.ListOrganizationTenantMonthlyReportCountResponse{}
	err := suite.hdl.ListOrganizationTenantMonthlyReportCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.DataAccessFailed, err)
}

func (suite *ListOrganizationTenantMonthlyReportCountTestSuite) TearDownSuite() {
}

func TestListOrganizationTenantMonthlyReportCountTestSuite(t *testing.T) {
	suite.Run(t, new(ListOrganizationTenantMonthlyReportCountTestSuite))
}
