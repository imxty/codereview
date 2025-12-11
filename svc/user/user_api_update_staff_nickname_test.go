package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/suite"
)

type UpdateStaffNicknameTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *UpdateStaffNicknameTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *UpdateStaffNicknameTestSuite) TestUpdateStaffNickname() {
	ctx := context.Background()

	req := &userpb.UpdateStaffNicknameRequest{
		TenantId:    tenantId,
		StaffId:     staffId,
		NewNickname: name,
	}
	resp := &userpb.UpdateStaffNicknameResponse{}
	err := suite.hdl.UpdateStaffNickname(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *UpdateStaffNicknameTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffNicknameRequest{
		TenantId:    tenantIdIsNull,
		StaffId:     staffId,
		NewNickname: name,
	}
	resp := &userpb.UpdateStaffNicknameResponse{}
	err := suite.hdl.UpdateStaffNickname(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdIsNull
func (suite *UpdateStaffNicknameTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffNicknameRequest{
		TenantId:    tenantId,
		StaffId:     staffIdIsNull,
		NewNickname: name,
	}
	resp := &userpb.UpdateStaffNicknameResponse{}
	err := suite.hdl.UpdateStaffNickname(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdNotExist
func (suite *UpdateStaffNicknameTestSuite) TestStaffIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffNicknameRequest{
		TenantId:    tenantId,
		StaffId:     staffIdNotExist,
		NewNickname: name,
	}
	resp := &userpb.UpdateStaffNicknameResponse{}
	err := suite.hdl.UpdateStaffNickname(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

// TestStaffIdIsErr
func (suite *UpdateStaffNicknameTestSuite) TestStaffIdIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffNicknameRequest{
		TenantId:    tenantId,
		StaffId:     staffIdNotActivated1,
		NewNickname: name,
	}
	resp := &userpb.UpdateStaffNicknameResponse{}
	err := suite.hdl.UpdateStaffNickname(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffHasNotActivated, err)
}
func (suite *UpdateStaffNicknameTestSuite) TearDownSuite() {

}

func TestUpdateStaffNicknameTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateStaffNicknameTestSuite))
}
