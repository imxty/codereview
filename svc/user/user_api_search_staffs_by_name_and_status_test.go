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
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type SearchStaffsByNameAndStatusTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *SearchStaffsByNameAndStatusTestSuite) SetupSuite() {
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

func (suite *SearchStaffsByNameAndStatusTestSuite) TestSearchStaffsByNameAndStatus() {
	ctx := context.Background()
	req := &userpb.SearchStaffsByNameAndStatusRequest{
		TenantId:  tenantId,
		StaffName: name,
		IsActivated: &wrapperspb.BoolValue{
			Value: true,
		},
	}
	resp := &userpb.SearchStaffsByNameAndStatusResponse{}
	err := suite.hdl.SearchStaffsByNameAndStatus(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *SearchStaffsByNameAndStatusTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SearchStaffsByNameAndStatusRequest{
		TenantId:  tenantIdIsNull,
		StaffName: name,
		IsActivated: &wrapperspb.BoolValue{
			Value: true,
		},
	}
	resp := &userpb.SearchStaffsByNameAndStatusResponse{}
	err := suite.hdl.SearchStaffsByNameAndStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SearchStaffsByNameAndStatusTestSuite) TearDownSuite() {
}

func TestSearchStaffsByNameAndStatusTestSuite(t *testing.T) {
	suite.Run(t, new(SearchStaffsByNameAndStatusTestSuite))
}
