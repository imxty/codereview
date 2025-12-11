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

type GetOrganizationsNameTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetOrganizationsNameTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetOrganizationsNameTestSuite) TestGetOrganizationsName() {
	ctx := context.Background()
	req := &userpb.GetOrganizationsNameRequest{
		OrganizationIds: []string{organizationId},
	}
	resp := &userpb.GetOrganizationsNameResponse{}
	err := suite.hdl.GetOrganizationsName(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *GetOrganizationsNameTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetOrganizationsNameRequest{
		OrganizationIds: []string{},
	}
	resp := &userpb.GetOrganizationsNameResponse{}
	err := suite.hdl.GetOrganizationsName(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationsNameTestSuite) TearDownSuite() {

}

func TestGetOrganizationsNameTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationsNameTestSuite))
}
