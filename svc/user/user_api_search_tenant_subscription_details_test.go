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

type SearchTenantSubscriptionDetailsTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *SearchTenantSubscriptionDetailsTestSuite) SetupSuite() {
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

func (suite *SearchTenantSubscriptionDetailsTestSuite) TestSearchTenantSubscriptionDetails() {
	ctx := context.Background()
	req := &userpb.SearchTenantSubscriptionDetailsRequest{
		OrganizationId: organizationId,
		TenantId:       tenantId,
		StartTime:      &timestamppb.Timestamp{},
		EndTime:        &timestamppb.Timestamp{},
		ContactName:    contactName,
		ContactPhone:   phone,
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.SearchTenantSubscriptionDetailsResponse{}
	err := suite.hdl.SearchTenantSubscriptionDetails(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SearchTenantSubscriptionDetailsTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SearchTenantSubscriptionDetailsRequest{
		OrganizationId: organizationIdIsNull,
		TenantId:       tenantId,
		StartTime:      nil,
		EndTime:        nil,
		ContactName:    contactName,
		ContactPhone:   phone,
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.SearchTenantSubscriptionDetailsResponse{}
	err := suite.hdl.SearchTenantSubscriptionDetails(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SearchTenantSubscriptionDetailsTestSuite) TearDownSuite() {
}

func TestSearchTenantSubscriptionDetailsTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTenantSubscriptionDetailsTestSuite))
}
