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

type GetOrganizationTenantSummaryTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetOrganizationTenantSummaryTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetOrganizationTenantSummaryTestSuite) TestGetOrganizationTenantSummary() {
	ctx := context.Background()
	req := &userpb.GetOrganizationTenantSummaryRequest{
		OrganizationId: organizationId,
	}
	resp := &userpb.GetOrganizationTenantSummaryResponse{}
	err := suite.hdl.GetOrganizationTenantSummary(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *GetOrganizationTenantSummaryTestSuite) TestOrganizationIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &userpb.GetOrganizationTenantSummaryRequest{
		OrganizationId: organizationIdIsNull,
	}
	resp := &userpb.GetOrganizationTenantSummaryResponse{}
	err := suite.hdl.GetOrganizationTenantSummary(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationTenantSummaryTestSuite) TearDownSuite() {

}

func TestGetOrganizationTenantSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationTenantSummaryTestSuite))
}
