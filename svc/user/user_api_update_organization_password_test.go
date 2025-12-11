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

type UpdateOrganizationPasswordTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *UpdateOrganizationPasswordTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *UpdateOrganizationPasswordTestSuite) TestUpdateOrganizationPassword() {
	ctx := context.Background()
	req := &userpb.UpdateOrganizationPasswordRequest{
		OrganizationId:   organizationId,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdateOrganizationPasswordResponse{}
	err := suite.hdl.UpdateOrganizationPassword(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *UpdateOrganizationPasswordTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOrganizationPasswordRequest{
		OrganizationId:   organizationIdIsNull,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdateOrganizationPasswordResponse{}
	err := suite.hdl.UpdateOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOrganizationIdNotExist
func (suite *UpdateOrganizationPasswordTestSuite) TestOrganizationIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOrganizationPasswordRequest{
		OrganizationId:   organizationIdNotExist,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdateOrganizationPasswordResponse{}
	err := suite.hdl.UpdateOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}

// TestOldPlainPasswordIsNull
func (suite *UpdateOrganizationPasswordTestSuite) TestOldPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOrganizationPasswordRequest{
		OrganizationId:   organizationId,
		OldPlainPassword: plainPasswordIsErr,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdateOrganizationPasswordResponse{}
	err := suite.hdl.UpdateOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrWrongOldPassword, err)
}

// TestOldPlainPasswordIsErr
func (suite *UpdateOrganizationPasswordTestSuite) TestOldPlainPasswordIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOrganizationPasswordRequest{
		OrganizationId:   organizationId,
		OldPlainPassword: plainPasswordIsErr,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdateOrganizationPasswordResponse{}
	err := suite.hdl.UpdateOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrWrongOldPassword, err)
}

// TestNewPlainPasswordIsNull
func (suite *UpdateOrganizationPasswordTestSuite) TestNewPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOrganizationPasswordRequest{
		OrganizationId:   organizationId,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPasswordIsNull,
	}
	resp := &userpb.UpdateOrganizationPasswordResponse{}
	err := suite.hdl.UpdateOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *UpdateOrganizationPasswordTestSuite) TearDownSuite() {

}

func TestUpdateOrganizationPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateOrganizationPasswordTestSuite))
}
