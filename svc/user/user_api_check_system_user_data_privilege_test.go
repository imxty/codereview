package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/suite"
)

type CheckSystemUserDataPrivilegeTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *CheckSystemUserDataPrivilegeTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *CheckSystemUserDataPrivilegeTestSuite) TestCheckSystemUserDataPrivilege() {
	ctx := context.Background()
	req := &userpb.CheckSystemUserDataPrivilegeRequest{
		UserId:         userId,
		OrganizationId: organizationId,
	}
	resp := &userpb.CheckSystemUserDataPrivilegeResponse{}
	err := suite.hdl.CheckSystemUserDataPrivilege(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestUserIdIsNull
func (suite *CheckSystemUserDataPrivilegeTestSuite) TestUserIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.CheckSystemUserDataPrivilegeRequest{
		UserId:         userIdIsNull,
		OrganizationId: organizationId,
	}
	resp := &userpb.CheckSystemUserDataPrivilegeResponse{}
	err := suite.hdl.CheckSystemUserDataPrivilege(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOrganizationIsNull
func (suite *CheckSystemUserDataPrivilegeTestSuite) TestOrganizationIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &userpb.CheckSystemUserDataPrivilegeRequest{
		UserId:         userId,
		OrganizationId: organizationIdIsNull,
	}
	resp := &userpb.CheckSystemUserDataPrivilegeResponse{}
	err := suite.hdl.CheckSystemUserDataPrivilege(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *CheckSystemUserDataPrivilegeTestSuite) TearDownSuite() {

}

func TestCheckSystemUserDataPrivilegeTestSuite(t *testing.T) {
	suite.Run(t, new(CheckSystemUserDataPrivilegeTestSuite))
}
