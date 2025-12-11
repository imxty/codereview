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

type ListPrivilegeOrganizationNamesTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *ListPrivilegeOrganizationNamesTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *ListPrivilegeOrganizationNamesTestSuite) TestListPrivilegeOrganizationNames() {
	ctx := context.Background()
	req := &userpb.ListPrivilegeOrganizationNamesRequest{
		PrivilegeId: privilegeId,
	}
	resp := &userpb.ListPrivilegeOrganizationNamesResponse{}
	err := suite.hdl.ListPrivilegeOrganizationNames(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}
func (suite *ListPrivilegeOrganizationNamesTestSuite) TestPrivilegeIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &userpb.ListPrivilegeOrganizationNamesRequest{
		PrivilegeId: privilegeIdIsNull,
	}
	resp := &userpb.ListPrivilegeOrganizationNamesResponse{}
	err := suite.hdl.ListPrivilegeOrganizationNames(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListPrivilegeOrganizationNamesTestSuite) TearDownSuite() {
}

func TestListPrivilegeOrganizationNamesTestSuite(t *testing.T) {
	suite.Run(t, new(ListPrivilegeOrganizationNamesTestSuite))
}
