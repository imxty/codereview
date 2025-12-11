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

type SignInTenantAdminTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *SignInTenantAdminTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *SignInTenantAdminTestSuite) TestSignInTenantAdmin() {
	ctx := context.Background()
	req := &userpb.SignInTenantAdminRequest{
		Phone:    phone,
		Password: password,
	}
	resp := &userpb.SignInTenantAdminResponse{}
	err := suite.hdl.SignInTenantAdmin(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPhoneIsNull
func (suite *SignInTenantAdminTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInTenantAdminRequest{
		Phone:    phoneIsNull,
		Password: password,
	}
	resp := &userpb.SignInTenantAdminResponse{}
	err := suite.hdl.SignInTenantAdmin(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestphoneNotExist
func (suite *SignInTenantAdminTestSuite) TestPhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInTenantAdminRequest{
		Phone:    phoneNotExist,
		Password: password,
	}
	resp := &userpb.SignInTenantAdminResponse{}
	err := suite.hdl.SignInTenantAdmin(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

// TestPasswordIsNull
func (suite *SignInTenantAdminTestSuite) TestPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInTenantAdminRequest{
		Phone:    phone,
		Password: passwordIsNull,
	}
	resp := &userpb.SignInTenantAdminResponse{}
	err := suite.hdl.SignInTenantAdmin(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPasswordIsErr
func (suite *SignInTenantAdminTestSuite) TestPasswordIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInTenantAdminRequest{
		Phone:    phone,
		Password: passwordIsErr,
	}
	resp := &userpb.SignInTenantAdminResponse{}
	err := suite.hdl.SignInTenantAdmin(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrWrongPassword, err)
}

func (suite *SignInTenantAdminTestSuite) TearDownSuite() {
}

func TestSignInTenantAdminTestSuite(t *testing.T) {
	suite.Run(t, new(SignInTenantAdminTestSuite))
}
