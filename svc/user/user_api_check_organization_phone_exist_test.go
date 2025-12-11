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

type CheckOrganizationPhoneExistTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *CheckOrganizationPhoneExistTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *CheckOrganizationPhoneExistTestSuite) TestCheckOrganizationPhoneExist() {
	ctx := context.Background()
	req := &userpb.CheckOrganizationPhoneExistRequest{
		Phone: phoneIsNew,
	}
	resp := &userpb.CheckOrganizationPhoneExistResponse{}
	err := suite.hdl.CheckOrganizationPhoneExist(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(true, resp.Exist)
}

func (suite *CheckOrganizationPhoneExistTestSuite) TestCheckOrganizationPhoneNotExist() {
	ctx := context.Background()
	req := &userpb.CheckOrganizationPhoneExistRequest{
		Phone: phoneNotExist,
	}
	resp := &userpb.CheckOrganizationPhoneExistResponse{}
	err := suite.hdl.CheckOrganizationPhoneExist(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(false, resp.Exist)
}

// TestPhoneIsNull
func (suite *CheckOrganizationPhoneExistTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.CheckOrganizationPhoneExistRequest{
		Phone: phoneIsNull,
	}
	resp := &userpb.CheckOrganizationPhoneExistResponse{}
	err := suite.hdl.CheckOrganizationPhoneExist(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *CheckOrganizationPhoneExistTestSuite) TearDownSuite() {

}

func TestCheckOrganizationPhoneExistTestSuite(t *testing.T) {
	suite.Run(t, new(CheckOrganizationPhoneExistTestSuite))
}
