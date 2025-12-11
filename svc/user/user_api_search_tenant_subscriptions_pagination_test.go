package user

import (
	"context"
	"testing"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	notificationmock "github.com/jinmukeji/huimaibao-service/svc/notification/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type SearchTenantSubscriptionsPaginationTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *SearchTenantSubscriptionsPaginationTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	notification := &notificationmock.NotificationAPIService{}
	suite.notification = notification
	suite.hdl = NewUserAPIHandler(s, notification, nil, nil, nil, nil, nil)
}

func (suite *SearchTenantSubscriptionsPaginationTestSuite) TestSearchTenantSubscriptionsPagination() {
	ctx := context.Background()
	req := &userpb.SearchTenantSubscriptionsPaginationRequest{
		OrganizationId: organizationId,
		StartTime:      &timestamppb.Timestamp{},
		EndTime:        &timestamppb.Timestamp{},
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.SearchTenantSubscriptionsPaginationResponse{}
	err := suite.hdl.SearchTenantSubscriptionsPagination(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SearchTenantSubscriptionsPaginationTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SearchTenantSubscriptionsPaginationRequest{
		OrganizationId: organizationIdIsNull,
		StartTime:      nil,
		EndTime:        &timestamppb.Timestamp{},
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.SearchTenantSubscriptionsPaginationResponse{}
	err := suite.hdl.SearchTenantSubscriptionsPagination(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SearchTenantSubscriptionsPaginationTestSuite) TearDownSuite() {
}

func TestSearchTenantSubscriptionsPaginationTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTenantSubscriptionsPaginationTestSuite))
}
