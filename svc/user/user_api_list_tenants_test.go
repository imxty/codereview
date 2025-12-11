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

type ListTenantsTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *ListTenantsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *ListTenantsTestSuite) TestListTenants() {
	ctx := context.Background()
	req := &userpb.ListTenantsRequest{
		OrganizationId: organizationId,
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.ListTenantsResponse{}
	err := suite.hdl.ListTenants(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *ListTenantsTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ListTenantsRequest{
		OrganizationId: organizationIdIsNull,
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.ListTenantsResponse{}
	err := suite.hdl.ListTenants(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPaginationIsNull
func (suite *ListTenantsTestSuite) TestPaginationIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ListTenantsRequest{
		OrganizationId: organizationHasStencil,
		Pagination:     nil,
	}
	resp := &userpb.ListTenantsResponse{}
	err := suite.hdl.ListTenants(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListTenantsTestSuite) TearDownSuite() {
}

func TestListTenantsTestSuite(t *testing.T) {
	suite.Run(t, new(ListTenantsTestSuite))
}
