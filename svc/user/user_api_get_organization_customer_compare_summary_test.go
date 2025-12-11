package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetOrganizationCustomerCompareSummaryTestSuite struct {
	suite.Suite
	hdl      *UserAPIHandler
	customer *customermock.CustomerAPIService
}

func (suite *GetOrganizationCustomerCompareSummaryTestSuite) SetupSuite() {
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

func (suite *GetOrganizationCustomerCompareSummaryTestSuite) TestGetOrganizationCustomerCompareSummary() {
	ctx := context.Background()
	suite.customer.On("GetOrganizationCustomerCompareSummary", mock.Anything, mock.Anything).Return(
		&customerpb.GetOrganizationCustomerCompareSummaryResponse{
			MonthAddedCustomerCount: 1,
			YearOnYear:              0.1,
			MonthOnMonth:            0.1,
		}, nil).After(utils.RpcLatency())
	req := &userpb.GetOrganizationCustomerCompareSummaryRequest{
		OrganizationId: organizationId,
		Date: &timestamppb.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
	}
	resp := &userpb.GetOrganizationCustomerCompareSummaryResponse{}
	err := suite.hdl.GetOrganizationCustomerCompareSummary(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *GetOrganizationCustomerCompareSummaryTestSuite) TestOrganizationIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &userpb.GetOrganizationCustomerCompareSummaryRequest{
		OrganizationId: organizationIdIsNull,
		Date: &timestamppb.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
	}
	resp := &userpb.GetOrganizationCustomerCompareSummaryResponse{}
	err := suite.hdl.GetOrganizationCustomerCompareSummary(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationCustomerCompareSummaryTestSuite) TearDownSuite() {

}

func TestGetOrganizationCustomerCompareSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationCustomerCompareSummaryTestSuite))
}
