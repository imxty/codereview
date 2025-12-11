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

type ListSystemUserDataPrivilegesTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *ListSystemUserDataPrivilegesTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *ListSystemUserDataPrivilegesTestSuite) TestListSystemUserDataPrivileges() {
	ctx := context.Background()
	req := &userpb.ListSystemUserDataPrivilegesRequest{
		UserId: userId1,
	}
	resp := &userpb.ListSystemUserDataPrivilegesResponse{}
	err := suite.hdl.ListSystemUserDataPrivileges(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestUserIdIsNull
func (suite *ListSystemUserDataPrivilegesTestSuite) TestUserIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ListSystemUserDataPrivilegesRequest{
		UserId: userIdIsNull,
	}
	resp := &userpb.ListSystemUserDataPrivilegesResponse{}
	err := suite.hdl.ListSystemUserDataPrivileges(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestUserIdNotExist
func (suite *ListSystemUserDataPrivilegesTestSuite) TestUserIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ListSystemUserDataPrivilegesRequest{
		UserId: userIdNotExist,
	}
	resp := &userpb.ListSystemUserDataPrivilegesResponse{}
	err := suite.hdl.ListSystemUserDataPrivileges(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrSystemUserNotFound, err)
}

func (suite *ListSystemUserDataPrivilegesTestSuite) TearDownSuite() {
}

func TestListSystemUserDataPrivilegesTestSuite(t *testing.T) {
	suite.Run(t, new(ListSystemUserDataPrivilegesTestSuite))
}
