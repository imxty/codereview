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
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SearchOrganizationReportsCountTestSuite struct {
	suite.Suite
	hdl *ReportAPIHandler
}

func (suite *SearchOrganizationReportsCountTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, nil, nil, "", "")
}

// TestSearchOrganizationReportsCount
func (suite *SearchOrganizationReportsCountTestSuite) TestSearchOrganizationReportsCount() {
	ctx := context.Background()
	req := &reportpb.SearchOrganizationReportsCountRequest{
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
	}
	resp := &reportpb.SearchOrganizationReportsCountResponse{}
	err := suite.hdl.SearchOrganizationReportsCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SearchOrganizationReportsCountTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.SearchOrganizationReportsCountRequest{
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
	}
	resp := &reportpb.SearchOrganizationReportsCountResponse{}
	err := suite.hdl.SearchOrganizationReportsCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantNameIsNull
func (suite *SearchOrganizationReportsCountTestSuite) TestTenantNameIsNull() {
	ctx := context.Background()
	req := &reportpb.SearchOrganizationReportsCountRequest{
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
	}
	resp := &reportpb.SearchOrganizationReportsCountResponse{}
	err := suite.hdl.SearchOrganizationReportsCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *SearchOrganizationReportsCountTestSuite) TearDownSuite() {
}

func TestSearchOrganizationReportsCountTestSuite(t *testing.T) {
	suite.Run(t, new(SearchOrganizationReportsCountTestSuite))
}
