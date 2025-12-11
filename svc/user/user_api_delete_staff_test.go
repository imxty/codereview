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

type DeleteStaffTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *DeleteStaffTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *DeleteStaffTestSuite) TestDeleteStaff() {
	ctx := context.Background()
	req := &userpb.DeleteStaffRequest{
		StaffId: staffId,
	}
	resp := &userpb.DeleteStaffResponse{}
	err := suite.hdl.DeleteStaff(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestStaffIdIsNull
func (suite *DeleteStaffTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.DeleteStaffRequest{
		StaffId: staffIdIsNull,
	}
	resp := &userpb.DeleteStaffResponse{}
	err := suite.hdl.DeleteStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdNotExist
func (suite *DeleteStaffTestSuite) TestStaffIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.DeleteStaffRequest{
		StaffId: staffIdNotExist,
	}
	resp := &userpb.DeleteStaffResponse{}
	err := suite.hdl.DeleteStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}
func (suite *DeleteStaffTestSuite) TearDownSuite() {

}

func TestDeleteStaffTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteStaffTestSuite))
}
