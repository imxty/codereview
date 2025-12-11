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

type LeaveTenantTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *LeaveTenantTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *LeaveTenantTestSuite) TestLeaveTenant() {
	ctx := context.Background()
	req := &userpb.LeaveTenantRequest{
		TenantId: tenantId,
		StaffId:  staffId,
	}
	resp := &userpb.LeaveTenantResponse{}
	err := suite.hdl.LeaveTenant(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *LeaveTenantTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.LeaveTenantRequest{
		TenantId: tenantIdIsNull,
		StaffId:  staffId,
	}
	resp := &userpb.LeaveTenantResponse{}
	err := suite.hdl.LeaveTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdIsNull
func (suite *LeaveTenantTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.LeaveTenantRequest{
		TenantId: tenantId,
		StaffId:  staffIdIsNull,
	}
	resp := &userpb.LeaveTenantResponse{}
	err := suite.hdl.LeaveTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdNotExist
func (suite *LeaveTenantTestSuite) TestStaffIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.LeaveTenantRequest{
		TenantId: tenantId,
		StaffId:  staffIdNotExist,
	}
	resp := &userpb.LeaveTenantResponse{}
	err := suite.hdl.LeaveTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

func (suite *LeaveTenantTestSuite) TearDownSuite() {
}

func TestLeaveTenantTestSuite(t *testing.T) {
	suite.Run(t, new(LeaveTenantTestSuite))
}
