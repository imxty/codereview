package user

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GetEntityCustomerCompareSummaryTestSuite struct {
	suite.Suite
	hdl      *UserAPIHandler
	customer *customermock.CustomerAPIService
}

func (suite *GetEntityCustomerCompareSummaryTestSuite) SetupSuite() {
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

func (suite *GetEntityCustomerCompareSummaryTestSuite) TestGetEntityCustomerCompareSummary() {
	ctx := context.Background()
	suite.customer.On("GetEntityCustomerCompareSummary", mock.Anything, mock.Anything).Return(
		&customerpb.GetEntityCustomerCompareSummaryResponse{
			MonthAddedCustomerCount: 1,
			YearOnYear:              0.1,
			MonthOnMonth:            0.1,
		}, nil).After(utils.RpcLatency())
	req := &userpb.GetEntityCustomerCompareSummaryRequest{
		TenantId: tenantId,
		Date: &timestamppb.Timestamp{
			Seconds: 9,
			Nanos:   1,
		},
	}
	resp := &userpb.GetEntityCustomerCompareSummaryResponse{}
	err := suite.hdl.GetEntityCustomerCompareSummary(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *GetEntityCustomerCompareSummaryTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetEntityCustomerCompareSummaryRequest{
		TenantId: tenantIdIsNull,
		Date: &timestamppb.Timestamp{
			Seconds: 9,
			Nanos:   1,
		}}
	resp := &userpb.GetEntityCustomerCompareSummaryResponse{}
	err := suite.hdl.GetEntityCustomerCompareSummary(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetEntityCustomerCompareSummaryTestSuite) TearDownSuite() {

}

func TestGetEntityCustomerCompareSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(GetEntityCustomerCompareSummaryTestSuite))
}
