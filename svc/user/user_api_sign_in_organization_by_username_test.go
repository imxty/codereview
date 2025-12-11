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

type SignInOrganizationByUsernameTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *SignInOrganizationByUsernameTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *SignInOrganizationByUsernameTestSuite) TestSignInOrganizationByUsername() {
	ctx := context.Background()
	req := &userpb.SignInOrganizationByUsernameRequest{
		Username: username,
		Password: password,
	}
	resp := &userpb.SignInOrganizationByUsernameResponse{}
	err := suite.hdl.SignInOrganizationByUsername(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestUsernameIsNull
func (suite *SignInOrganizationByUsernameTestSuite) TestUsernameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInOrganizationByUsernameRequest{
		Username: usernameIsNull,
		Password: password,
	}
	resp := &userpb.SignInOrganizationByUsernameResponse{}
	err := suite.hdl.SignInOrganizationByUsername(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestUsernameNotExist
func (suite *SignInOrganizationByUsernameTestSuite) TestUsernameNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInOrganizationByUsernameRequest{
		Username: usernameNotExist,
		Password: password,
	}
	resp := &userpb.SignInOrganizationByUsernameResponse{}
	err := suite.hdl.SignInOrganizationByUsername(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}

// TestPasswordIsNull
func (suite *SignInOrganizationByUsernameTestSuite) TestPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInOrganizationByUsernameRequest{
		Username: username,
		Password: passwordIsNull,
	}
	resp := &userpb.SignInOrganizationByUsernameResponse{}
	err := suite.hdl.SignInOrganizationByUsername(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SignInOrganizationByUsernameTestSuite) TearDownSuite() {
}

func TestSignInOrganizationByUsernameTestSuite(t *testing.T) {
	suite.Run(t, new(SignInOrganizationByUsernameTestSuite))
}
