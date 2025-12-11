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

type ListSystemUsersByIdsTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *ListSystemUsersByIdsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *ListSystemUsersByIdsTestSuite) TestListSystemUsersByIds() {
	ctx := context.Background()
	req := &userpb.ListSystemUsersByIdsRequest{
		SystemUserIds: []string{userId},
	}
	resp := &userpb.ListSystemUsersByIdsResponse{}
	err := suite.hdl.ListSystemUsersByIds(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestSystemUserIdsIsNull
func (suite *ListSystemUsersByIdsTestSuite) TestSystemUserIdsIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ListSystemUsersByIdsRequest{
		SystemUserIds: []string{},
	}
	resp := &userpb.ListSystemUsersByIdsResponse{}
	err := suite.hdl.ListSystemUsersByIds(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListSystemUsersByIdsTestSuite) TearDownSuite() {
}

func TestListSystemUsersByIdsTestSuite(t *testing.T) {
	suite.Run(t, new(ListSystemUsersByIdsTestSuite))
}
