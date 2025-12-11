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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ListTopOrganizationReportCountTestSuite struct {
	suite.Suite
	hdl  *ReportAPIHandler
	user *usermock.UserAPIService
}

func (suite *ListTopOrganizationReportCountTestSuite) SetupSuite() {
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

// TestListTopOrganizationReportCount 获取组织商户测量top10
func (suite *ListTopOrganizationReportCountTestSuite) TestListTopOrganizationReportCount() {
	ctx := context.Background()
	suite.user.On("GetOrganizationTenants", mock.Anything, mock.Anything).Return(
		&userpb.GetOrganizationTenantsResponse{
			Tenants: []*userpb.TenantEntity{
				&userpb.TenantEntity{
					OrganizationId: organizationId,
				},
			},
		}, nil).After(utils.RpcLatency())
	req := &reportpb.ListTopOrganizationReportCountRequest{
		OrganizationId: organizationId,
	}
	resp := &reportpb.ListTopOrganizationReportCountResponse{}
	err := suite.hdl.ListTopOrganizationReportCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *ListTopOrganizationReportCountTestSuite) TestOrganizationIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &reportpb.ListTopOrganizationReportCountRequest{
		OrganizationId: organizationIdIsNull,
	}
	resp := &reportpb.ListTopOrganizationReportCountResponse{}
	err := suite.hdl.ListTopOrganizationReportCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListTopOrganizationReportCountTestSuite) TearDownSuite() {
}

func TestListTopOrganizationReportCountTestSuite(t *testing.T) {
	suite.Run(t, new(ListTopOrganizationReportCountTestSuite))
}
