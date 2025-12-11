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

type SearchTenantsByNameTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *SearchTenantsByNameTestSuite) SetupSuite() {
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

func (suite *SearchTenantsByNameTestSuite) TestSearchTenantsByName() {
	ctx := context.Background()
	req := &userpb.SearchTenantsByNameRequest{
		TenantName: tenantName,
	}
	resp := &userpb.SearchTenantsByNameResponse{}
	err := suite.hdl.SearchTenantsByName(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantNameIsNull
func (suite *SearchTenantsByNameTestSuite) TestTenantNameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SearchTenantsByNameRequest{
		TenantName: tenantNameIsNull,
	}
	resp := &userpb.SearchTenantsByNameResponse{}
	err := suite.hdl.SearchTenantsByName(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SearchTenantsByNameTestSuite) TearDownSuite() {
}

func TestSearchTenantsByNameTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTenantsByNameTestSuite))
}
