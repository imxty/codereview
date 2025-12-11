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

type SendReviewNotificationTestSuite struct {
	suite.Suite
	hdl    *NotificationAPIHandler
	aliyun *alimock.SmsSend
}

func (suite *SendReviewNotificationTestSuite) SetupSuite() {
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

func (suite *SendReviewNotificationTestSuite) SetupTest() {
	suite.aliyun.On("SendNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		true, nil).After(utils.RpcLatency())

}

// TestSendReviewNotification 组织注册
func (suite *SendReviewNotificationTestSuite) TestSendReviewNotification() {
	ctx := context.Background()

	req := &notificationpb.SendReviewNotificationRequest{
		IssueNumber: 9,
		Phone:       []string{phone},
	}
	resp := &notificationpb.SendReviewNotificationResponse{}
	err := suite.hdl.SendReviewNotification(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPhoneIsNull
func (suite *SendReviewNotificationTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendReviewNotificationRequest{
		IssueNumber: 9,
		Phone:       []string{},
	}
	resp := &notificationpb.SendReviewNotificationResponse{}
	err := suite.hdl.SendReviewNotification(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SendReviewNotificationTestSuite) TearDownSuite() {
}

func TestSendReviewNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(SendReviewNotificationTestSuite))
}
