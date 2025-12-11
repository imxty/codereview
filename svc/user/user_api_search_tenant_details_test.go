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

type SearchTenantDetailsTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *SearchTenantDetailsTestSuite) SetupSuite() {
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

func (suite *SearchTenantDetailsTestSuite) TestSearchTenantDetails() {
	ctx := context.Background()
	req := &userpb.SearchTenantDetailsRequest{
		OrganizationId: organizationId,
		TenantName:     name,
		StartTime:      &timestamppb.Timestamp{},
		EndTime:        &timestamppb.Timestamp{},
		ContactName:    contactName,
		ContactPhone:   phone,
		Status:         userpb.TenantStatus_TENANT_STATUS_USING,
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.SearchTenantDetailsResponse{}
	err := suite.hdl.SearchTenantDetails(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SearchTenantDetailsTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SearchTenantDetailsRequest{
		OrganizationId: organizationIdIsNull,
		TenantName:     name,
		StartTime:      &timestamppb.Timestamp{},
		EndTime:        &timestamppb.Timestamp{},
		ContactName:    contactName,
		ContactPhone:   phone,
		Status:         userpb.TenantStatus_TENANT_STATUS_INVALID,
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.SearchTenantDetailsResponse{}
	err := suite.hdl.SearchTenantDetails(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SearchTenantDetailsTestSuite) TearDownSuite() {
}

func TestSearchTenantDetailsTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTenantDetailsTestSuite))
}
