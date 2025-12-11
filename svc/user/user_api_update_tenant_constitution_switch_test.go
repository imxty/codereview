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
	usermock "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/stretchr/testify/suite"
)

type UpdateTenantConstitutionSwitchTestSuite struct {
	suite.Suite
	hdl  *UserAPIHandler
	user *usermock.UserAPIService
}

func (suite *UpdateTenantConstitutionSwitchTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	user := &usermock.UserAPIService{}
	suite.user = user
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *UpdateTenantConstitutionSwitchTestSuite) TestUpdateTenantConstitutionSwitch() {
	ctx := context.Background()

	req := &userpb.UpdateTenantConstitutionSwitchRequest{
		TenantId: tenantId,
	}
	resp := &userpb.UpdateTenantConstitutionSwitchResponse{}
	err := suite.hdl.UpdateTenantConstitutionSwitch(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *UpdateTenantConstitutionSwitchTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateTenantConstitutionSwitchRequest{
		TenantId: tenantIdIsNull,
	}
	resp := &userpb.UpdateTenantConstitutionSwitchResponse{}
	err := suite.hdl.UpdateTenantConstitutionSwitch(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *UpdateTenantConstitutionSwitchTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateTenantConstitutionSwitchRequest{
		TenantId: tenantIdNotExist,
	}
	resp := &userpb.UpdateTenantConstitutionSwitchResponse{}
	err := suite.hdl.UpdateTenantConstitutionSwitch(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

func (suite *UpdateTenantConstitutionSwitchTestSuite) TearDownSuite() {

}

func TestUpdateTenantConstitutionSwitchTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTenantConstitutionSwitchTestSuite))
}
