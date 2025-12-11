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

type GetTenantTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetTenantTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetTenantTestSuite) TestGetTenant() {
	ctx := context.Background()
	req := &userpb.GetTenantRequest{
		OrganizationId: organizationId,
		TenantId:       tenantId,
	}
	resp := &userpb.GetTenantResponse{}
	err := suite.hdl.GetTenant(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(tenantId, resp.Tenant.TenantId)
}

// TestTenantIdIsNull
func (suite *GetTenantTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetTenantRequest{
		OrganizationId: organizationId,
		TenantId:       tenantIdIsNull,
	}
	resp := &userpb.GetTenantResponse{}
	err := suite.hdl.GetTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *GetTenantTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetTenantRequest{
		OrganizationId: organizationId,
		TenantId:       tenantIdNotExist,
	}
	resp := &userpb.GetTenantResponse{}
	err := suite.hdl.GetTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

// TestOrganizationIsNull
func (suite *GetTenantTestSuite) TestOrganizationIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetTenantRequest{
		OrganizationId: organizationIdIsNull,
		TenantId:       tenantId,
	}
	resp := &userpb.GetTenantResponse{}
	err := suite.hdl.GetTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetTenantTestSuite) TearDownSuite() {
}

func TestGetTenantTestSuite(t *testing.T) {
	suite.Run(t, new(GetTenantTestSuite))
}
