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

type SignInAppByPasswordTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *SignInAppByPasswordTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *SignInAppByPasswordTestSuite) TestSignInAppByPassword() {
	ctx := context.Background()
	req := &userpb.SignInAppByPasswordRequest{
		Phone:         staffPhoneIsExist,
		PlainPassword: plainPassword,
	}
	resp := &userpb.SignInAppByPasswordResponse{}
	err := suite.hdl.SignInAppByPassword(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPhoneIsNull
func (suite *SignInAppByPasswordTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInAppByPasswordRequest{
		Phone:         phoneIsNull,
		PlainPassword: plainPassword,
	}
	resp := &userpb.SignInAppByPasswordResponse{}
	err := suite.hdl.SignInAppByPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsErr
func (suite *SignInAppByPasswordTestSuite) TestPhoneIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInAppByPasswordRequest{
		Phone:         phoneIsErr,
		PlainPassword: plainPassword,
	}
	resp := &userpb.SignInAppByPasswordResponse{}
	err := suite.hdl.SignInAppByPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

// TestPhoneNotExist
func (suite *SignInAppByPasswordTestSuite) TestPhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInAppByPasswordRequest{
		Phone:         phoneNotExist,
		PlainPassword: plainPassword,
	}
	resp := &userpb.SignInAppByPasswordResponse{}
	err := suite.hdl.SignInAppByPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

// TestPlainPasswordIsNull
func (suite *SignInAppByPasswordTestSuite) TestPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInAppByPasswordRequest{
		Phone:         phone,
		PlainPassword: plainPasswordIsNull,
	}
	resp := &userpb.SignInAppByPasswordResponse{}
	err := suite.hdl.SignInAppByPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPlainPasswordIsErr
func (suite *SignInAppByPasswordTestSuite) TestPlainPasswordIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInAppByPasswordRequest{
		Phone:         phone,
		PlainPassword: plainPasswordIsErr,
	}
	resp := &userpb.SignInAppByPasswordResponse{}
	err := suite.hdl.SignInAppByPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

func (suite *SignInAppByPasswordTestSuite) TearDownSuite() {
}

func TestSignInAppByPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(SignInAppByPasswordTestSuite))
}
