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

type ModifySharingReportSwitchTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *ModifySharingReportSwitchTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *ModifySharingReportSwitchTestSuite) TestModifySharingReportSwitch() {
	ctx := context.Background()
	req := &userpb.ModifySharingReportSwitchRequest{
		TenantId: tenantId,
	}
	resp := &userpb.ModifySharingReportSwitchResponse{}
	err := suite.hdl.ModifySharingReportSwitch(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *ModifySharingReportSwitchTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ModifySharingReportSwitchRequest{
		TenantId: tenantIdIsNull,
	}
	resp := &userpb.ModifySharingReportSwitchResponse{}
	err := suite.hdl.ModifySharingReportSwitch(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *ModifySharingReportSwitchTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ModifySharingReportSwitchRequest{
		TenantId: tenantIdNotExist,
	}
	resp := &userpb.ModifySharingReportSwitchResponse{}
	err := suite.hdl.ModifySharingReportSwitch(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

func (suite *ModifySharingReportSwitchTestSuite) TearDownSuite() {
}

func TestModifySharingReportSwitchTestSuite(t *testing.T) {
	suite.Run(t, new(ModifySharingReportSwitchTestSuite))
}
