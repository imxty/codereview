package user

import (
	"context"
	"testing"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type ListPrivilegeGroupsTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *ListPrivilegeGroupsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *ListPrivilegeGroupsTestSuite) TestListPrivilegeGroups() {
	ctx := context.Background()
	req := &userpb.ListPrivilegeGroupsRequest{}
	resp := &userpb.ListPrivilegeGroupsResponse{}
	err := suite.hdl.ListPrivilegeGroups(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *ListPrivilegeGroupsTestSuite) TearDownSuite() {
}

func TestListPrivilegeGroupsTestSuite(t *testing.T) {
	suite.Run(t, new(ListPrivilegeGroupsTestSuite))
}
