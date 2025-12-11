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

type UpdateStaffTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *UpdateStaffTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *UpdateStaffTestSuite) TestUpdateStaff() {
	ctx := context.Background()

	req := &userpb.UpdateStaffRequest{
		StaffId:       staffId,
		Phone:         phone,
		Name:          name,
		PlainPassword: plainPassword,
	}
	resp := &userpb.UpdateStaffResponse{}
	err := suite.hdl.UpdateStaff(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestStaffIdIsNull
func (suite *UpdateStaffTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffRequest{
		StaffId:       staffIdIsNull,
		Phone:         phone,
		Name:          name,
		PlainPassword: plainPassword,
	}
	resp := &userpb.UpdateStaffResponse{}
	err := suite.hdl.UpdateStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdNotExist
func (suite *UpdateStaffTestSuite) TestStaffIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffRequest{
		StaffId:       staffIdNotExist,
		Phone:         phone,
		Name:          name,
		PlainPassword: plainPassword,
	}
	resp := &userpb.UpdateStaffResponse{}
	err := suite.hdl.UpdateStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

// TestPhoneIsNull
func (suite *UpdateStaffTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffRequest{
		StaffId:       staffId,
		Phone:         phoneIsNull,
		Name:          name,
		PlainPassword: plainPassword,
	}
	resp := &userpb.UpdateStaffResponse{}
	err := suite.hdl.UpdateStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestNameIsNull
func (suite *UpdateStaffTestSuite) TestNameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffRequest{
		StaffId:       staffId,
		Phone:         phone,
		Name:          nameIsNull,
		PlainPassword: plainPassword,
	}
	resp := &userpb.UpdateStaffResponse{}
	err := suite.hdl.UpdateStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *UpdateStaffTestSuite) TearDownSuite() {

}

func TestUpdateStaffTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateStaffTestSuite))
}
