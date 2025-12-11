package user

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	notificationmock "github.com/jinmukeji/huimaibao-service/svc/notification/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SearchOrganizationTenantCompareTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
	customer     *customermock.CustomerAPIService
}

func (suite *SearchOrganizationTenantCompareTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	notification := &notificationmock.NotificationAPIService{}
	suite.notification = notification
	customer := &customermock.CustomerAPIService{}
	suite.customer = customer
	suite.hdl = NewUserAPIHandler(s, notification, nil, customer, nil, nil, nil)
}

func (suite *SearchOrganizationTenantCompareTestSuite) TestSearchOrganizationTenantCompare() {
	ctx := context.Background()
	suite.customer.On("SearchOrganizationTenantCompare", mock.Anything, mock.Anything).Return(
		&customerpb.SearchOrganizationTenantCompareResponse{
			Tenants: []*customerpb.TenantCustomerCompare{
				&customerpb.TenantCustomerCompare{
					TenantName: tenantName,
				},
			},
			TotalCount: 3,
		}, nil).After(utils.RpcLatency())
	req := &userpb.SearchOrganizationTenantCompareRequest{
		OrganizationId: organizationId,
		TenantName:     tenantName,
		StartTime: &timestamppb.Timestamp{
			Seconds: 1,
			Nanos:   1,
		},
		EndTime: &timestamppb.Timestamp{
			Seconds: 1,
			Nanos:   1},
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.SearchOrganizationTenantCompareResponse{}
	err := suite.hdl.SearchOrganizationTenantCompare(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SearchOrganizationTenantCompareTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SearchOrganizationTenantCompareRequest{
		OrganizationId: organizationIdIsNull,
		TenantName:     tenantName,
		StartTime: &timestamppb.Timestamp{
			Seconds: 1,
			Nanos:   1,
		},
		EndTime: &timestamppb.Timestamp{
			Seconds: 1,
			Nanos:   1},
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.SearchOrganizationTenantCompareResponse{}
	err := suite.hdl.SearchOrganizationTenantCompare(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantNameIsNull
func (suite *SearchOrganizationTenantCompareTestSuite) TestTenantNameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SearchOrganizationTenantCompareRequest{
		OrganizationId: organizationId,
		TenantName:     tenantNameIsNull,
		StartTime:      nil,
		EndTime: &timestamppb.Timestamp{
			Seconds: 1,
			Nanos:   1},
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.SearchOrganizationTenantCompareResponse{}
	err := suite.hdl.SearchOrganizationTenantCompare(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPaginationIsNull
func (suite *SearchOrganizationTenantCompareTestSuite) TestPaginationIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SearchOrganizationTenantCompareRequest{
		OrganizationId: organizationId,
		TenantName:     tenantNameIsNull,
		StartTime:      nil,
		EndTime: &timestamppb.Timestamp{
			Seconds: 1,
			Nanos:   1},
		Pagination: nil,
	}
	resp := &userpb.SearchOrganizationTenantCompareResponse{}
	err := suite.hdl.SearchOrganizationTenantCompare(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SearchOrganizationTenantCompareTestSuite) TearDownSuite() {
}

func TestSearchOrganizationTenantCompareTestSuite(t *testing.T) {
	suite.Run(t, new(SearchOrganizationTenantCompareTestSuite))
}
