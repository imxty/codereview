package user

import (
	"context"
	"testing"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	notificationmock "github.com/jinmukeji/huimaibao-service/svc/notification/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type SearchTenantsTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *SearchTenantsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	notification := &notificationmock.NotificationAPIService{}
	suite.notification = notification
	suite.hdl = NewUserAPIHandler(s, notification, nil, nil, nil, nil, nil)
}

func (suite *SearchTenantsTestSuite) TestSearchTenants() {
	ctx := context.Background()
	req := &userpb.SearchTenantsRequest{
		OrganizationId: organizationId,
		Status:         userpb.TenantStatus_TENANT_STATUS_USING,
		Pagination: &userpb.Pagination{
			Size:   size,
			Offset: offset,
		},
	}
	resp := &userpb.SearchTenantsResponse{}
	err := suite.hdl.SearchTenants(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SearchTenantsTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SearchTenantsRequest{
		OrganizationId: organizationIdIsNull,
		Status:         userpb.TenantStatus_TENANT_STATUS_INVALID,
		Pagination: &userpb.Pagination{
			Size:   size,
			Offset: offset,
		},
	}
	resp := &userpb.SearchTenantsResponse{}
	err := suite.hdl.SearchTenants(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SearchTenantsTestSuite) TearDownSuite() {
}

func TestSearchTenantsTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTenantsTestSuite))
}
