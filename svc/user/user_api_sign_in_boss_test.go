package user

import (
	"context"
	"testing"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SignInBossTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *SignInBossTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *SignInBossTestSuite) TestSignInBoss() {
	ctx := context.Background()
	req := &userpb.SignInBossRequest{
		Phone:    phone,
		Password: plainPassword,
	}
	resp := &userpb.SignInBossResponse{}
	err := suite.hdl.SignInBoss(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPasswordIsNull
func (suite *SignInBossTestSuite) TestPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInBossRequest{
		Phone:    phone,
		Password: plainPasswordIsNull,
	}
	resp := &userpb.SignInBossResponse{}
	err := suite.hdl.SignInBoss(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsNull
func (suite *SignInBossTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInBossRequest{
		Phone:    phoneIsNull,
		Password: plainPassword,
	}
	resp := &userpb.SignInBossResponse{}
	err := suite.hdl.SignInBoss(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsErr
func (suite *SignInBossTestSuite) TestPhoneIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInBossRequest{
		Phone:    phoneIsErr,
		Password: plainPassword,
	}
	resp := &userpb.SignInBossResponse{}
	err := suite.hdl.SignInBoss(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrSystemUserNotFound, err)
}

func (suite *SignInBossTestSuite) TearDownSuite() {

}

func TestSignInBossTestSuite(t *testing.T) {
	suite.Run(t, new(SignInBossTestSuite))
}
