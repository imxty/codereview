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

type GetTenantEntityTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetTenantEntityTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetTenantEntityTestSuite) TestGetTenantEntity() {
	ctx := context.Background()
	req := &userpb.GetTenantEntityRequest{
		TenantId: tenantId,
	}
	resp := &userpb.GetTenantEntityResponse{}
	err := suite.hdl.GetTenantEntity(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(tenantId, resp.Entity.TenantId)
	suite.Assert().Equal(true, resp.ReportSharingSwitchStatus)
}

func (suite *GetTenantEntityTestSuite) TestReportNoSharing() {
	ctx := context.Background()
	req := &userpb.GetTenantEntityRequest{
		TenantId: tenantIdIsActivated,
	}
	resp := &userpb.GetTenantEntityResponse{}
	err := suite.hdl.GetTenantEntity(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(tenantIdIsActivated, resp.Entity.TenantId)
	suite.Assert().Equal(false, resp.ReportSharingSwitchStatus)
}

// TestTenantIdIsNull
func (suite *GetTenantEntityTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetTenantEntityRequest{
		TenantId: tenantIdIsNull,
	}
	resp := &userpb.GetTenantEntityResponse{}
	err := suite.hdl.GetTenantEntity(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *GetTenantEntityTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetTenantEntityRequest{
		TenantId: tenantIdNotExist,
	}
	resp := &userpb.GetTenantEntityResponse{}
	err := suite.hdl.GetTenantEntity(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

func (suite *GetTenantEntityTestSuite) TearDownSuite() {

}

func TestGetTenantEntityTestSuite(t *testing.T) {
	suite.Run(t, new(GetTenantEntityTestSuite))
}
