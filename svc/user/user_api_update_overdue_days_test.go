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

type UpdateOverdueDaysTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *UpdateOverdueDaysTestSuite) SetupSuite() {
}
func (suite *UpdateOverdueDaysTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *UpdateOverdueDaysTestSuite) TestUpdateOverdueDays() {
	ctx := context.Background()
	req := &userpb.UpdateOverdueDaysRequest{
		TenantId: tenantId,
		Days:     YearDays,
	}
	resp := &userpb.UpdateOverdueDaysResponse{}
	err := suite.hdl.UpdateOverdueDays(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *UpdateOverdueDaysTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOverdueDaysRequest{
		TenantId: tenantIdIsNull,
		Days:     YearDays,
	}
	resp := &userpb.UpdateOverdueDaysResponse{}
	err := suite.hdl.UpdateOverdueDays(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestDayIsNull
func (suite *UpdateOverdueDaysTestSuite) TestDayIsNull() {

	ctx := context.Background()
	req := &userpb.UpdateOverdueDaysRequest{
		TenantId: tenantIdIsActivated,
		Days:     0,
	}
	resp := &userpb.UpdateOverdueDaysResponse{}
	err := suite.hdl.UpdateOverdueDays(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestDayIsErr
func (suite *UpdateOverdueDaysTestSuite) TestDayIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOverdueDaysRequest{
		TenantId: tenantIdIsActivated,
		Days:     -1,
	}
	resp := &userpb.UpdateOverdueDaysResponse{}
	err := suite.hdl.UpdateOverdueDays(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *UpdateOverdueDaysTestSuite) TearDownSuite() {

}

func TestUpdateOverdueDaysTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateOverdueDaysTestSuite))
}
