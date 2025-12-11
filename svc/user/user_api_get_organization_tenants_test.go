package user

import (
	"context"
	"testing"

	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	reviewmock "github.com/jinmukeji/huimaibao-service/svc/review/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetOrganizationTenantsTestSuite struct {
	suite.Suite
	hdl    *UserAPIHandler
	review *reviewmock.ReviewAPIService
}

func (suite *GetOrganizationTenantsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	review := &reviewmock.ReviewAPIService{}
	suite.review = review
	suite.hdl = NewUserAPIHandler(s, nil, review, nil, nil, nil, nil)
}

func (suite *GetOrganizationTenantsTestSuite) TestGetOrganizationTenants() {
	ctx := context.Background()
	suite.review.On("ListCertificateReviewStatus", mock.Anything, mock.Anything).Return(
		&reviewpb.ListCertificateReviewStatusResponse{
			Results: make(map[string]*reviewpb.ReviewResult),
		}, nil).After(utils.RpcLatency())
	req := &userpb.GetOrganizationTenantsRequest{
		OrganizationId: organizationId,
	}
	resp := &userpb.GetOrganizationTenantsResponse{}
	err := suite.hdl.GetOrganizationTenants(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *GetOrganizationTenantsTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetOrganizationTenantsRequest{
		OrganizationId: organizationIdIsNull,
	}
	resp := &userpb.GetOrganizationTenantsResponse{}
	err := suite.hdl.GetOrganizationTenants(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationTenantsTestSuite) TearDownSuite() {

}

func TestGetOrganizationTenantsTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationTenantsTestSuite))
}
