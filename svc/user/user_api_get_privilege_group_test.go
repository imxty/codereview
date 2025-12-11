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

type GetPrivilegeGroupTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetPrivilegeGroupTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetPrivilegeGroupTestSuite) TestGetPrivilegeGroup() {
	ctx := context.Background()
	req := &userpb.GetPrivilegeGroupRequest{
		PrivilegeId: privilegeId,
	}
	resp := &userpb.GetPrivilegeGroupResponse{}
	err := suite.hdl.GetPrivilegeGroup(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *GetPrivilegeGroupTestSuite) TestPrivilegeIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &userpb.GetPrivilegeGroupRequest{
		PrivilegeId: privilegeIdIsNull,
	}
	resp := &userpb.GetPrivilegeGroupResponse{}
	err := suite.hdl.GetPrivilegeGroup(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.DataAccessFailed, err)
}

func (suite *GetPrivilegeGroupTestSuite) TearDownSuite() {

}

func TestGetPrivilegeGroupTestSuite(t *testing.T) {
	suite.Run(t, new(GetPrivilegeGroupTestSuite))
}
