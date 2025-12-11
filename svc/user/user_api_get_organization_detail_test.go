package user

import (
	"context"
	"testing"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GetOrganizationDetailTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetOrganizationDetailTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetOrganizationDetailTestSuite) TestGetOrganizationDetail() {
	ctx := context.Background()
	req := &userpb.GetOrganizationDetailRequest{
		OrganizationId: organizationId,
	}
	resp := &userpb.GetOrganizationDetailResponse{}
	err := suite.hdl.GetOrganizationDetail(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(organizationId, resp.Organization.OrganizationId)
	suite.Assert().Equal(nameIsNew, resp.Organization.Name)
	suite.Assert().Equal(phoneIsNew, resp.Organization.Phone)
}

func (suite *GetOrganizationDetailTestSuite) TestArrearsIsTrue() {
	ctx := context.Background()
	req := &userpb.GetOrganizationDetailRequest{
		OrganizationId: organizationHasStencil,
	}
	resp := &userpb.GetOrganizationDetailResponse{}
	err := suite.hdl.GetOrganizationDetail(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(organizationHasStencil, resp.Organization.OrganizationId)
	suite.Assert().Equal(name1, resp.Organization.Name)
	suite.Assert().Equal(phoneIsExist, resp.Organization.Phone)
}

// TestOrganizationIdIsNull
func (suite *GetOrganizationDetailTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetOrganizationDetailRequest{
		OrganizationId: organizationIdIsNull,
	}
	resp := &userpb.GetOrganizationDetailResponse{}
	err := suite.hdl.GetOrganizationDetail(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOrganizationIdNotExist
func (suite *GetOrganizationDetailTestSuite) TestOrganizationIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetOrganizationDetailRequest{
		OrganizationId: organizationIdNotExist,
	}
	resp := &userpb.GetOrganizationDetailResponse{}
	err := suite.hdl.GetOrganizationDetail(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}

func (suite *GetOrganizationDetailTestSuite) TearDownSuite() {

}

func TestGetOrganizationDetailTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationDetailTestSuite))
}
