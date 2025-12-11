package customer

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/customer"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GetOrganizationCustomerCompareSummaryTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *GetOrganizationCustomerCompareSummaryTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

func (suite *GetOrganizationCustomerCompareSummaryTestSuite) TestGetOrganizationCustomerCompareSummary() {
	ctx := context.Background()
	req := &customerpb.GetOrganizationCustomerCompareSummaryRequest{
		OrganizationId: organizationId,
		Date: &timestamppb.Timestamp{
			Seconds: 0,
			Nanos:   1,
		},
	}
	resp := &customerpb.GetOrganizationCustomerCompareSummaryResponse{}
	err := suite.hdl.GetOrganizationCustomerCompareSummary(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(0), resp.MonthAddedCustomerCount)

}

// TestOrganizationIdIsNull
func (suite *GetOrganizationCustomerCompareSummaryTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.GetOrganizationCustomerCompareSummaryRequest{
		OrganizationId: organizationIdIsNull,
		Date: &timestamppb.Timestamp{
			Seconds: 0,
			Nanos:   1,
		},
	}
	resp := &customerpb.GetOrganizationCustomerCompareSummaryResponse{}
	err := suite.hdl.GetOrganizationCustomerCompareSummary(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestDateIsNull
func (suite *GetOrganizationCustomerCompareSummaryTestSuite) TestDateIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.GetOrganizationCustomerCompareSummaryRequest{
		OrganizationId: organizationId,
		Date:           nil,
	}
	resp := &customerpb.GetOrganizationCustomerCompareSummaryResponse{}
	err := suite.hdl.GetOrganizationCustomerCompareSummary(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationCustomerCompareSummaryTestSuite) TearDownSuite() {

}

func TestGetOrganizationCustomerCompareSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationCustomerCompareSummaryTestSuite))
}
