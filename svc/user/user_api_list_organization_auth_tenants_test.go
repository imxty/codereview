package user

import (
	"context"
	"testing"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type ListOrganizationAuthTenantsTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *ListOrganizationAuthTenantsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *ListOrganizationAuthTenantsTestSuite) TestListOrganizationAuthTenants() {
	ctx := context.Background()
	req := &userpb.ListOrganizationAuthTenantsRequest{
		OrganizationId: organizationId,
	}
	resp := &userpb.ListOrganizationAuthTenantsResponse{}
	err := suite.hdl.ListOrganizationAuthTenants(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *ListOrganizationAuthTenantsTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ListOrganizationAuthTenantsRequest{
		OrganizationId: organizationIdIsNull,
	}
	resp := &userpb.ListOrganizationAuthTenantsResponse{}
	err := suite.hdl.ListOrganizationAuthTenants(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListOrganizationAuthTenantsTestSuite) TearDownSuite() {
}

func TestListOrganizationAuthTenantsTestSuite(t *testing.T) {
	suite.Run(t, new(ListOrganizationAuthTenantsTestSuite))
}
