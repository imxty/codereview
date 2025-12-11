package notification

import (
	"context"
	"testing"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/notification"
	alimock "github.com/jinmukeji/huimaibao-service/svc/notification/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SendBillingNotificationTestSuite struct {
	suite.Suite
	hdl    *NotificationAPIHandler
	aliyun *alimock.SmsSend
}

func (suite *SendBillingNotificationTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, notificationFile)
	s := store.NewNotificationStore(conn)
	aliyun := &alimock.SmsSend{}
	suite.aliyun = aliyun
	suite.hdl = NewNotificationAPIHandler(s, aliyun, "")
}
func (suite *SendBillingNotificationTestSuite) SetupTest() {
	suite.aliyun.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		true, nil).After(utils.RpcLatency())

}

// TestSendBillingNotification
func (suite *SendBillingNotificationTestSuite) TestSendBillingNotification() {
	ctx := context.Background()

	req := &notificationpb.SendBillingNotificationRequest{
		Notifications: []*notificationpb.BillingNotification{
			&notificationpb.BillingNotification{
				Phone:            phone,
				OrganizationId:   organizationId,
				OrganizationName: organizationName,
				Amount:           10000,
			},
		},
	}
	resp := &notificationpb.SendBillingNotificationResponse{}
	err := suite.hdl.SendBillingNotification(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPhoneIsNull
func (suite *SendBillingNotificationTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendBillingNotificationRequest{
		Notifications: []*notificationpb.BillingNotification{
			&notificationpb.BillingNotification{
				Phone:            phoneIsNull,
				OrganizationId:   organizationId,
				OrganizationName: organizationName,
				Amount:           10000,
			},
		},
	}
	resp := &notificationpb.SendBillingNotificationResponse{}
	err := suite.hdl.SendBillingNotification(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOrganizationIdIsNull
func (suite *SendBillingNotificationTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendBillingNotificationRequest{
		Notifications: []*notificationpb.BillingNotification{
			&notificationpb.BillingNotification{
				Phone:            phone,
				OrganizationId:   organizationIdIsNull,
				OrganizationName: organizationName,
				Amount:           10000,
			},
		},
	}
	resp := &notificationpb.SendBillingNotificationResponse{}
	err := suite.hdl.SendBillingNotification(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOrganizationNameIsNull
func (suite *SendBillingNotificationTestSuite) TestOrganizationNameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendBillingNotificationRequest{
		Notifications: []*notificationpb.BillingNotification{
			&notificationpb.BillingNotification{
				Phone:            phone,
				OrganizationId:   organizationId,
				OrganizationName: organizationNameIsNull,
				Amount:           10000,
			},
		},
	}
	resp := &notificationpb.SendBillingNotificationResponse{}
	err := suite.hdl.SendBillingNotification(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestAmountIsNull
func (suite *SendBillingNotificationTestSuite) TestAmountIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendBillingNotificationRequest{
		Notifications: []*notificationpb.BillingNotification{
			&notificationpb.BillingNotification{
				Phone:            phone,
				OrganizationId:   organizationId,
				OrganizationName: organizationName,
				Amount:           -1,
			},
		},
	}
	resp := &notificationpb.SendBillingNotificationResponse{}
	err := suite.hdl.SendBillingNotification(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestNotifivationIsNull
func (suite *SendBillingNotificationTestSuite) TestNotifivationIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendBillingNotificationRequest{
		Notifications: []*notificationpb.BillingNotification{}}
	resp := &notificationpb.SendBillingNotificationResponse{}
	err := suite.hdl.SendBillingNotification(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SendBillingNotificationTestSuite) TearDownSuite() {
}

func TestSendBillingNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(SendBillingNotificationTestSuite))
}
