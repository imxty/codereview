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

type UpdatePasswordTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *UpdatePasswordTestSuite) SetupSuite() {
}
func (suite *UpdatePasswordTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *UpdatePasswordTestSuite) TestUpdatePassword() {
	ctx := context.Background()
	req := &userpb.UpdatePasswordRequest{
		TenantId:         tenantId,
		StaffId:          staffId,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdatePasswordResponse{}
	err := suite.hdl.UpdatePassword(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *UpdatePasswordTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdatePasswordRequest{
		TenantId:         tenantIdIsNull,
		StaffId:          staffId,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdatePasswordResponse{}
	err := suite.hdl.UpdatePassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdIsNull
func (suite *UpdatePasswordTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdatePasswordRequest{
		TenantId:         tenantIdIsActivated,
		StaffId:          staffIdIsNull,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdatePasswordResponse{}
	err := suite.hdl.UpdatePassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdNotExist
func (suite *UpdatePasswordTestSuite) TestStaffIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdatePasswordRequest{
		TenantId:         tenantIdIsActivated,
		StaffId:          staffIdNotExist,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdatePasswordResponse{}
	err := suite.hdl.UpdatePassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

// TestOldPlainPasswordIsNull
func (suite *UpdatePasswordTestSuite) TestOldPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdatePasswordRequest{
		TenantId:         tenantIdIsActivated,
		StaffId:          staffId,
		OldPlainPassword: plainPasswordIsNull,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdatePasswordResponse{}
	err := suite.hdl.UpdatePassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOldPlainPasswordIsErr
func (suite *UpdatePasswordTestSuite) TestOldPlainPasswordIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdatePasswordRequest{
		TenantId:         tenantIdIsActivated,
		StaffId:          staffId,
		OldPlainPassword: plainPasswordIsErr,
		NewPlainPassword: newPlainPassword,
	}
	resp := &userpb.UpdatePasswordResponse{}
	err := suite.hdl.UpdatePassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrWrongOldPassword, err)
}

// TestNewPlainPasswordIsNull
func (suite *UpdatePasswordTestSuite) TestNewPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdatePasswordRequest{
		TenantId:         tenantIdIsActivated,
		StaffId:          staffId,
		OldPlainPassword: plainPassword,
		NewPlainPassword: newPlainPasswordIsNull,
	}
	resp := &userpb.UpdatePasswordResponse{}
	err := suite.hdl.UpdatePassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPasswordIsSame
func (suite *UpdatePasswordTestSuite) TestPasswordIsSame() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdatePasswordRequest{
		TenantId:         tenantIdIsActivated,
		StaffId:          staffId,
		OldPlainPassword: plainPassword,
		NewPlainPassword: plainPassword,
	}
	resp := &userpb.UpdatePasswordResponse{}
	err := suite.hdl.UpdatePassword(ctx, req, resp)
	suite.T().Log(resp)
	utils.AssertRpcErrorCode(t, ErrSamePasswords, err)
}

func (suite *UpdatePasswordTestSuite) TearDownSuite() {

}

func TestUpdatePasswordTestSuite(t *testing.T) {
	suite.Run(t, new(UpdatePasswordTestSuite))
}
