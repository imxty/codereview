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

type AddStaffTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *AddStaffTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *AddStaffTestSuite) TestAddStaff() {
	ctx := context.Background()
	req := &userpb.AddStaffRequest{
		TenantId:      tenantId,
		Phone:         staffPhone,
		Name:          name,
		PlainPassword: plainPassword,
	}
	resp := &userpb.AddStaffResponse{}
	err := suite.hdl.AddStaff(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *AddStaffTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.AddStaffRequest{
		TenantId:      tenantIdIsNull,
		Phone:         staffPhone,
		Name:          name,
		PlainPassword: plainPassword,
	}
	resp := &userpb.AddStaffResponse{}
	err := suite.hdl.AddStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *AddStaffTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.AddStaffRequest{
		TenantId:      tenantIdNotExist,
		Phone:         staffPhone,
		Name:          name,
		PlainPassword: plainPassword,
	}
	resp := &userpb.AddStaffResponse{}
	err := suite.hdl.AddStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

// TestPhoneIsNull
func (suite *AddStaffTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.AddStaffRequest{
		TenantId:      tenantId,
		Phone:         phoneIsNull,
		Name:          name,
		PlainPassword: plainPassword,
	}
	resp := &userpb.AddStaffResponse{}
	err := suite.hdl.AddStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsExist

func (suite *AddStaffTestSuite) TestPhoneIsExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.AddStaffRequest{
		TenantId:      tenantId,
		Phone:         staffPhoneIsExist,
		Name:          name,
		PlainPassword: plainPassword,
	}
	resp := &userpb.AddStaffResponse{}
	err := suite.hdl.AddStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffPhoneHasBeenUsed, err)
}

// TestNameIsNull
func (suite *AddStaffTestSuite) TestNameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.AddStaffRequest{
		TenantId:      tenantId,
		Phone:         staffPhone,
		Name:          nameIsNull,
		PlainPassword: plainPassword,
	}
	resp := &userpb.AddStaffResponse{}
	err := suite.hdl.AddStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPlainPasswordIsNull
func (suite *AddStaffTestSuite) TestPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.AddStaffRequest{
		TenantId:      tenantId,
		Phone:         staffPhone,
		Name:          name,
		PlainPassword: plainPasswordIsNull,
	}
	resp := &userpb.AddStaffResponse{}
	err := suite.hdl.AddStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPlainPasswordIsErr
func (suite *AddStaffTestSuite) TestPlainPasswordIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.AddStaffRequest{
		TenantId:      tenantId,
		Phone:         staffPhone,
		Name:          name,
		PlainPassword: plainPasswordIsErr,
	}
	resp := &userpb.AddStaffResponse{}
	err := suite.hdl.AddStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestErrStaffExceedLimit
func (suite *AddStaffTestSuite) TestErrStaffExceedLimit() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.AddStaffRequest{
		TenantId:      tenantIdIsActivated,
		Phone:         phoneIsNew,
		Name:          name,
		PlainPassword: plainPassword,
	}
	resp := &userpb.AddStaffResponse{}
	err := suite.hdl.AddStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffPhoneHasBeenUsed, err)
}
func (suite *AddStaffTestSuite) TearDownSuite() {

}

func TestAddStaffTestSuite(t *testing.T) {
	suite.Run(t, new(AddStaffTestSuite))
}
