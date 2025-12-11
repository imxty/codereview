package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/suite"
)

type DeleteTenantTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *DeleteTenantTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *DeleteTenantTestSuite) TestDeleteTenant() {
	ctx := context.Background()
	req := &userpb.DeleteTenantRequest{
		UserId:         userId,
		OrganizationId: organizationId,
		TenantId:       tenantIdIsUnset,
	}
	resp := &userpb.DeleteTenantResponse{}
	err := suite.hdl.DeleteTenant(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *DeleteTenantTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.DeleteTenantRequest{
		UserId:         userId,
		OrganizationId: organizationId,
		TenantId:       tenantIdIsNull,
	}
	resp := &userpb.DeleteTenantResponse{}
	err := suite.hdl.DeleteTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *DeleteTenantTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.DeleteTenantRequest{
		OrganizationId: organizationId,
		UserId:         userId,
		TenantId:       tenantIdNotExist,
	}
	resp := &userpb.DeleteTenantResponse{}
	err := suite.hdl.DeleteTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}

// TestOrganizationIsNull
func (suite *DeleteTenantTestSuite) TestOrganizationIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.DeleteTenantRequest{
		UserId:         userId,
		OrganizationId: organizationIdIsNull,
		TenantId:       tenantId,
	}
	resp := &userpb.DeleteTenantResponse{}
	err := suite.hdl.DeleteTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOrganizationIsErr
func (suite *DeleteTenantTestSuite) TestOrganizationIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.DeleteTenantRequest{
		UserId:         userId,
		OrganizationId: organizationIdNotExist,
		TenantId:       tenantId,
	}
	resp := &userpb.DeleteTenantResponse{}
	err := suite.hdl.DeleteTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}

// TestUserIdIsNull
func (suite *DeleteTenantTestSuite) TestUserIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.DeleteTenantRequest{
		UserId:         userIdIsNull,
		OrganizationId: organizationId,
		TenantId:       tenantId,
	}
	resp := &userpb.DeleteTenantResponse{}
	err := suite.hdl.DeleteTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *DeleteTenantTestSuite) TearDownSuite() {

}

func TestDeleteTenantTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteTenantTestSuite))
}
