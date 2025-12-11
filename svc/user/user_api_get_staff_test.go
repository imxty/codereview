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

type GetStaffTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetStaffTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetStaffTestSuite) TestGetStaff() {
	ctx := context.Background()

	req := &userpb.GetStaffRequest{
		TenantId: tenantId,
		StaffId:  staffId,
	}
	resp := &userpb.GetStaffResponse{}
	err := suite.hdl.GetStaff(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(staffId, resp.Staff.StaffId)
	suite.Assert().Equal(name, resp.Staff.Name)
	suite.Assert().Equal(staffPhoneIsExist, resp.Staff.Phone)
}

// TestTenantIdIsNull
func (suite *GetStaffTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetStaffRequest{
		TenantId: tenantIdIsNull,
		StaffId:  staffId,
	}
	resp := &userpb.GetStaffResponse{}
	err := suite.hdl.GetStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

// TestStaffIdIsNull
func (suite *GetStaffTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetStaffRequest{
		TenantId: tenantId,
		StaffId:  staffIdIsNull,
	}
	resp := &userpb.GetStaffResponse{}
	err := suite.hdl.GetStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdNotExist
func (suite *GetStaffTestSuite) TestStaffIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetStaffRequest{
		TenantId: tenantId,
		StaffId:  staffIdNotExist,
	}
	resp := &userpb.GetStaffResponse{}
	err := suite.hdl.GetStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

func (suite *GetStaffTestSuite) TearDownSuite() {

}

func TestGetStaffTestSuite(t *testing.T) {
	suite.Run(t, new(GetStaffTestSuite))
}
