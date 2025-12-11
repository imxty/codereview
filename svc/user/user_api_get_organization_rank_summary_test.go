package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	reportmock "github.com/jinmukeji/huimaibao-service/svc/report/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetOrganizationRankSummaryTestSuite struct {
	suite.Suite
	hdl      *UserAPIHandler
	calc     *reportmock.ReportAPIService
	customer *customermock.CustomerAPIService
}

func (suite *GetOrganizationRankSummaryTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	calc := &reportmock.ReportAPIService{}
	suite.calc = calc
	customer := &customermock.CustomerAPIService{}
	suite.customer = customer
	suite.hdl = NewUserAPIHandler(s, nil, nil, customer, calc, nil, nil)
}

func (suite *GetOrganizationRankSummaryTestSuite) TestGetOrganizationRankSummary() {
	ctx := context.Background()
	suite.calc.On("ListTopOrganizationReportCount", mock.Anything, mock.Anything).Return(
		&reportpb.ListTopOrganizationReportCountResponse{
			MeasurementRank: []*reportpb.TenantRank{
				&reportpb.TenantRank{
					TenantId: tenantId,
				},
			},
		}, nil).After(utils.RpcLatency())
	suite.customer.On("GetTenantsCustomerRank", mock.Anything, mock.Anything).Return(
		&customerpb.GetTenantsCustomerRankResponse{
			Tenants: make(map[string]*customerpb.TenantRank),
		}, nil).After(utils.RpcLatency())
	req := &userpb.GetOrganizationRankSummaryRequest{
		OrganizationId: organizationId,
	}
	resp := &userpb.GetOrganizationRankSummaryResponse{}
	err := suite.hdl.GetOrganizationRankSummary(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *GetOrganizationRankSummaryTestSuite) TestOrganizationIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &userpb.GetOrganizationRankSummaryRequest{
		OrganizationId: organizationIdIsNull,
	}
	resp := &userpb.GetOrganizationRankSummaryResponse{}
	err := suite.hdl.GetOrganizationRankSummary(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationRankSummaryTestSuite) TearDownSuite() {

}

func TestGetOrganizationRankSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationRankSummaryTestSuite))
}
