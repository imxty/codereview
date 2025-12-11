package user

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetOrganizationOverdueSummaryTestSuite struct {
	suite.Suite
	hdl      *UserAPIHandler
	customer *customermock.CustomerAPIService
}

func (suite *GetOrganizationOverdueSummaryTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	customer := &customermock.CustomerAPIService{}
	suite.customer = customer
	suite.hdl = NewUserAPIHandler(s, nil, nil, customer, nil, nil, nil)
}

func (suite *GetOrganizationOverdueSummaryTestSuite) TestGetOrganizationOverdueSummary() {
	ctx := context.Background()
	suite.customer.On("GetOrganizationOverdueCount", mock.Anything, mock.Anything).Return(
		&customerpb.GetOrganizationOverdueCountResponse{
			TotalCount: 2,
		}, nil).After(utils.RpcLatency())
	req := &userpb.GetOrganizationOverdueSummaryRequest{
		OrganizationId: organizationId,
	}
	resp := &userpb.GetOrganizationOverdueSummaryResponse{}
	err := suite.hdl.GetOrganizationOverdueSummary(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(resp.OverdueCustomerTotalCount, int32(2))
}

func (suite *GetOrganizationOverdueSummaryTestSuite) TestOrganizationIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &userpb.GetOrganizationOverdueSummaryRequest{
		OrganizationId: organizationIdIsNull,
	}
	resp := &userpb.GetOrganizationOverdueSummaryResponse{}
	err := suite.hdl.GetOrganizationOverdueSummary(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationOverdueSummaryTestSuite) TearDownSuite() {

}

func TestGetOrganizationOverdueSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationOverdueSummaryTestSuite))
}
