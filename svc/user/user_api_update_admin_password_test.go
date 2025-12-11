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

type UpdateAdminPasswordTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *UpdateAdminPasswordTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *UpdateAdminPasswordTestSuite) TestUpdateAdminPassword() {
	ctx := context.Background()
	req := &userpb.UpdateAdminPasswordRequest{
		TenantId:         tenantId,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdateAdminPasswordResponse{}
	err := suite.hdl.UpdateAdminPassword(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *UpdateAdminPasswordTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateAdminPasswordRequest{
		TenantId:         tenantIdIsNull,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdateAdminPasswordResponse{}
	err := suite.hdl.UpdateAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *UpdateAdminPasswordTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateAdminPasswordRequest{
		TenantId:         tenantIdNotExist,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdateAdminPasswordResponse{}
	err := suite.hdl.UpdateAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

// TestOldPlainPasswordIsNull
func (suite *UpdateAdminPasswordTestSuite) TestOldPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateAdminPasswordRequest{
		TenantId:         tenantId,
		OldPlainPassword: plainPasswordIsErr,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdateAdminPasswordResponse{}
	err := suite.hdl.UpdateAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrWrongOldPassword, err)
}

// TestOldPlainPasswordIsErr
func (suite *UpdateAdminPasswordTestSuite) TestOldPlainPasswordIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateAdminPasswordRequest{
		TenantId:         tenantId,
		OldPlainPassword: plainPasswordIsErr,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdateAdminPasswordResponse{}
	err := suite.hdl.UpdateAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrWrongOldPassword, err)
}

// TestNewPlainPasswordIsNull
func (suite *UpdateAdminPasswordTestSuite) TestNewPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateAdminPasswordRequest{
		TenantId:         tenantId,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPasswordIsNull,
	}
	resp := &userpb.UpdateAdminPasswordResponse{}
	err := suite.hdl.UpdateAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *UpdateAdminPasswordTestSuite) TearDownSuite() {

}

func TestUpdateAdminPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateAdminPasswordTestSuite))
}
