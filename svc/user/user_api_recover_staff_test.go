package user

import (
	"context"
	"testing"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	notificationmock "github.com/jinmukeji/huimaibao-service/svc/notification/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type RecoverStaffTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *RecoverStaffTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	notification := &notificationmock.NotificationAPIService{}
	suite.notification = notification
	suite.hdl = NewUserAPIHandler(s, notification, nil, nil, nil, nil, nil)
}

func (suite *RecoverStaffTestSuite) TestRecoverStaff() {
	ctx := context.Background()

	req := &userpb.RecoverStaffRequest{
		TenantId: tenantIdIsActivated,
		StaffId:  staffIdNotActivated,
	}
	resp := &userpb.RecoverStaffResponse{}
	err := suite.hdl.RecoverStaff(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestStaffIdInOtherCompany
func (suite *RecoverStaffTestSuite) TestStaffIdInOtherCompany() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.RecoverStaffRequest{
		TenantId: tenantId,
		StaffId:  staffIdNotActivated1,
	}
	resp := &userpb.RecoverStaffResponse{}
	err := suite.hdl.RecoverStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffHasBeenAdded, err)
}
func (suite *RecoverStaffTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.RecoverStaffRequest{
		TenantId: tenantIdIsNull,
		StaffId:  staffId,
	}
	resp := &userpb.RecoverStaffResponse{}
	err := suite.hdl.RecoverStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdIsNull
func (suite *RecoverStaffTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.RecoverStaffRequest{
		TenantId: tenantId,
		StaffId:  staffIdIsNull,
	}
	resp := &userpb.RecoverStaffResponse{}
	err := suite.hdl.RecoverStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdNotFound
func (suite *RecoverStaffTestSuite) TestStaffIdNotFound() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.RecoverStaffRequest{
		TenantId: tenantIdIsActivated,
		StaffId:  staffIdNotExist,
	}
	resp := &userpb.RecoverStaffResponse{}
	err := suite.hdl.RecoverStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

// TestErrStaffHasActivated
func (suite *RecoverStaffTestSuite) TestErrStaffHasActivated() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.RecoverStaffRequest{
		TenantId: tenantId,
		StaffId:  staffId,
	}
	resp := &userpb.RecoverStaffResponse{}
	err := suite.hdl.RecoverStaff(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffHasActivated, err)
}

func (suite *RecoverStaffTestSuite) TearDownSuite() {
}

func TestRecoverStaffTestSuite(t *testing.T) {
	suite.Run(t, new(RecoverStaffTestSuite))
}
