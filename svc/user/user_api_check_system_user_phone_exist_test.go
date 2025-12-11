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

type CheckSystemUserPhoneExistTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *CheckSystemUserPhoneExistTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *CheckSystemUserPhoneExistTestSuite) TestCheckSystemUserPhoneExist() {
	ctx := context.Background()
	req := &userpb.CheckSystemUserPhoneExistRequest{
		Phone: phone,
	}
	resp := &userpb.CheckSystemUserPhoneExistResponse{}
	err := suite.hdl.CheckSystemUserPhoneExist(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPhoneIsNull
func (suite *CheckSystemUserPhoneExistTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.CheckSystemUserPhoneExistRequest{
		Phone: phoneIsNull,
	}
	resp := &userpb.CheckSystemUserPhoneExistResponse{}
	err := suite.hdl.CheckSystemUserPhoneExist(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *CheckSystemUserPhoneExistTestSuite) TearDownSuite() {

}

func TestCheckSystemUserPhoneExistTestSuite(t *testing.T) {
	suite.Run(t, new(CheckSystemUserPhoneExistTestSuite))
}
