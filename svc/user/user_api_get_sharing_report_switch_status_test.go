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

type GetSharingReportSwitchStatusTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetSharingReportSwitchStatusTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetSharingReportSwitchStatusTestSuite) TestGetSharingReportSwitchStatus() {
	ctx := context.Background()
	req := &userpb.GetSharingReportSwitchStatusRequest{
		TenantId: tenantId,
	}
	resp := &userpb.GetSharingReportSwitchStatusResponse{}
	err := suite.hdl.GetSharingReportSwitchStatus(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(true, resp.SharingReportSwitchStatus)
}

func (suite *GetSharingReportSwitchStatusTestSuite) TestGetSharingReportStatusFalse() {
	ctx := context.Background()
	req := &userpb.GetSharingReportSwitchStatusRequest{
		TenantId: tenantIdIsActivated,
	}
	resp := &userpb.GetSharingReportSwitchStatusResponse{}
	err := suite.hdl.GetSharingReportSwitchStatus(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(false, resp.SharingReportSwitchStatus)
}

// TestTenantIdIsNull
func (suite *GetSharingReportSwitchStatusTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetSharingReportSwitchStatusRequest{
		TenantId: tenantIdIsNull,
	}
	resp := &userpb.GetSharingReportSwitchStatusResponse{}
	err := suite.hdl.GetSharingReportSwitchStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *GetSharingReportSwitchStatusTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetSharingReportSwitchStatusRequest{
		TenantId: tenantIdNotExist,
	}
	resp := &userpb.GetSharingReportSwitchStatusResponse{}
	err := suite.hdl.GetSharingReportSwitchStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

func (suite *GetSharingReportSwitchStatusTestSuite) TearDownSuite() {

}

func TestGetSharingReportSwitchStatusTestSuite(t *testing.T) {
	suite.Run(t, new(GetSharingReportSwitchStatusTestSuite))
}
