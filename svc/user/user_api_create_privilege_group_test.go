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

type CreatePrivilegeGroupTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *CreatePrivilegeGroupTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *CreatePrivilegeGroupTestSuite) TestCreatePrivilegeGroup() {
	ctx := context.Background()
	req := &userpb.CreatePrivilegeGroupRequest{
		UserId:         userId,
		PrivilegeName:  privilegeName,
		Remark:         remark,
		PagePrivileges: []userpb.PagePrivilege{},
		DataPrivileges: []string{},
	}
	resp := &userpb.CreatePrivilegeGroupResponse{}
	err := suite.hdl.CreatePrivilegeGroup(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestUserIdIsNull
func (suite *CreatePrivilegeGroupTestSuite) TestUserIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.CreatePrivilegeGroupRequest{
		UserId:         userIdIsNull,
		PrivilegeName:  privilegeName,
		Remark:         remark,
		PagePrivileges: []userpb.PagePrivilege{},
		DataPrivileges: []string{},
	}
	resp := &userpb.CreatePrivilegeGroupResponse{}
	err := suite.hdl.CreatePrivilegeGroup(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestrivilegeNameIsNull
func (suite *CreatePrivilegeGroupTestSuite) TestrivilegeNameIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &userpb.CreatePrivilegeGroupRequest{
		UserId:         userId,
		PrivilegeName:  privilegeNameIsNull,
		Remark:         remark,
		PagePrivileges: []userpb.PagePrivilege{},
		DataPrivileges: []string{},
	}
	resp := &userpb.CreatePrivilegeGroupResponse{}
	err := suite.hdl.CreatePrivilegeGroup(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *CreatePrivilegeGroupTestSuite) TearDownSuite() {

}

func TestCreatePrivilegeGroupTestSuite(t *testing.T) {
	suite.Run(t, new(CreatePrivilegeGroupTestSuite))
}
