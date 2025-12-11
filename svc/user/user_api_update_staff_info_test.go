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

type UpdateStaffInfoTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *UpdateStaffInfoTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *UpdateStaffInfoTestSuite) TestUpdateStaffInfo() {
	ctx := context.Background()

	req := &userpb.UpdateStaffInfoRequest{
		UserId:      userId,
		StaffId:     staffId,
		StaffName:   name,
		StaffPhone:  phone,
		PrivilegeId: privilegeId,
		Remark:      remark,
	}
	resp := &userpb.UpdateStaffInfoResponse{}
	err := suite.hdl.UpdateStaffInfo(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestUserIdIsNull
func (suite *UpdateStaffInfoTestSuite) TestUserIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffInfoRequest{
		UserId:      userIdIsNull,
		StaffId:     staffId,
		StaffName:   name,
		StaffPhone:  phone,
		PrivilegeId: privilegeId,
		Remark:      remark,
	}
	resp := &userpb.UpdateStaffInfoResponse{}
	err := suite.hdl.UpdateStaffInfo(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdIsNull
func (suite *UpdateStaffInfoTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffInfoRequest{
		UserId:      userId,
		StaffId:     staffIdIsNull,
		StaffName:   name,
		StaffPhone:  phone,
		PrivilegeId: privilegeId,
		Remark:      remark,
	}
	resp := &userpb.UpdateStaffInfoResponse{}
	err := suite.hdl.UpdateStaffInfo(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffNameIsNull
func (suite *UpdateStaffInfoTestSuite) TestStaffNameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffInfoRequest{
		UserId:      userId,
		StaffId:     staffId,
		StaffName:   nameIsNull,
		StaffPhone:  phone,
		PrivilegeId: privilegeId,
		Remark:      remark,
	}
	resp := &userpb.UpdateStaffInfoResponse{}
	err := suite.hdl.UpdateStaffInfo(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffPhoneIsNull

func (suite *UpdateStaffInfoTestSuite) TestStaffPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffInfoRequest{
		UserId:      userId,
		StaffId:     staffId,
		StaffName:   name,
		StaffPhone:  phoneIsNull,
		PrivilegeId: privilegeId,
		Remark:      remark,
	}
	resp := &userpb.UpdateStaffInfoResponse{}
	err := suite.hdl.UpdateStaffInfo(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPrivilegeIdIsNull
func (suite *UpdateStaffInfoTestSuite) TestPrivilegeIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateStaffInfoRequest{
		UserId:      userId,
		StaffId:     staffId,
		StaffName:   name,
		StaffPhone:  phone,
		PrivilegeId: privilegeIdIsNull,
		Remark:      remark,
	}
	resp := &userpb.UpdateStaffInfoResponse{}
	err := suite.hdl.UpdateStaffInfo(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *UpdateStaffInfoTestSuite) TearDownSuite() {

}

func TestUpdateStaffInfoTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateStaffInfoTestSuite))
}
