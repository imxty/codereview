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

type CreateSystemUserTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *CreateSystemUserTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)

	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *CreateSystemUserTestSuite) TestCreateSystemUser() {
	ctx := context.Background()
	req := &userpb.CreateSystemUserRequest{
		StaffName:   name,
		StaffPhone:  phoneIsExist,
		PrivilegeId: privilegeId,
		Remark:      remark,
	}
	resp := &userpb.CreateSystemUserResponse{}
	err := suite.hdl.CreateSystemUser(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestStaffNameIsNull
func (suite *CreateSystemUserTestSuite) TestStaffNameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.CreateSystemUserRequest{
		StaffName:   nameIsNull,
		StaffPhone:  phone,
		PrivilegeId: privilegeId,
		Remark:      remark,
	}
	resp := &userpb.CreateSystemUserResponse{}
	err := suite.hdl.CreateSystemUser(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsNull
func (suite *CreateSystemUserTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.CreateSystemUserRequest{
		StaffName:   name,
		StaffPhone:  phoneIsNull,
		PrivilegeId: privilegeId,
		Remark:      remark,
	}
	resp := &userpb.CreateSystemUserResponse{}
	err := suite.hdl.CreateSystemUser(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPrivilegeIdIsNull
func (suite *CreateSystemUserTestSuite) TestPrivilegeIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.CreateSystemUserRequest{
		StaffName:   name,
		StaffPhone:  phone,
		PrivilegeId: privilegeIdIsNull,
		Remark:      remark,
	}
	resp := &userpb.CreateSystemUserResponse{}
	err := suite.hdl.CreateSystemUser(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *CreateSystemUserTestSuite) TearDownSuite() {

}

func TestCreateSystemUserTestSuite(t *testing.T) {
	suite.Run(t, new(CreateSystemUserTestSuite))
}
