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

type GetUnAuthTenantRevisionTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetUnAuthTenantRevisionTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetUnAuthTenantRevisionTestSuite) TestGetUnAuthTenantRevision() {
	ctx := context.Background()
	req := &userpb.GetUnAuthTenantRevisionRequest{
		TenantId: tenantId,
	}
	resp := &userpb.GetUnAuthTenantRevisionResponse{}
	err := suite.hdl.GetUnAuthTenantRevision(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *GetUnAuthTenantRevisionTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetUnAuthTenantRevisionRequest{
		TenantId: tenantIdIsNull,
	}
	resp := &userpb.GetUnAuthTenantRevisionResponse{}
	err := suite.hdl.GetUnAuthTenantRevision(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *GetUnAuthTenantRevisionTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetUnAuthTenantRevisionRequest{
		TenantId: tenantIdNotExist,
	}
	resp := &userpb.GetUnAuthTenantRevisionResponse{}
	err := suite.hdl.GetUnAuthTenantRevision(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

func (suite *GetUnAuthTenantRevisionTestSuite) TearDownSuite() {
}

func TestGetUnAuthTenantRevisionTestSuite(t *testing.T) {
	suite.Run(t, new(GetUnAuthTenantRevisionTestSuite))
}
