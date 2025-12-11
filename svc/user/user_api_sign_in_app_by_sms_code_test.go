package user

import (
	"context"
	"testing"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	notificationmock "github.com/jinmukeji/huimaibao-service/svc/notification/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SignInAppBySmsCodeTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *SignInAppBySmsCodeTestSuite) SetupTest() {
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

func (suite *SignInAppBySmsCodeTestSuite) TestSignInAppBySmsCode() {
	ctx := context.Background()
	suite.notification.On("GetLastestVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.GetLatestPhoneVerificationCodeResponse{
			SmsCode: smsCode,
			TxId:    txId,
		}, nil).After(utils.RpcLatency())
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())
	req := &userpb.SignInAppBySmsCodeRequest{
		Phone:   staffPhoneIsExist,
		SmsCode: smsCode,
		TxId:    txId,
	}
	resp := &userpb.SignInAppBySmsCodeResponse{}
	err := suite.hdl.SignInAppBySmsCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPhoneIsNull
func (suite *SignInAppBySmsCodeTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInAppBySmsCodeRequest{
		Phone:   phoneIsNull,
		SmsCode: smsCode,
	}
	resp := &userpb.SignInAppBySmsCodeResponse{}
	err := suite.hdl.SignInAppBySmsCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsErr
func (suite *SignInAppBySmsCodeTestSuite) TestPhoneIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInAppBySmsCodeRequest{
		Phone:   phoneIsErr,
		SmsCode: smsCode,
	}
	resp := &userpb.SignInAppBySmsCodeResponse{}
	err := suite.hdl.SignInAppBySmsCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneNotExist
func (suite *SignInAppBySmsCodeTestSuite) TestPhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInAppBySmsCodeRequest{
		Phone:   phoneNotExist,
		SmsCode: smsCode,
	}
	resp := &userpb.SignInAppBySmsCodeResponse{}
	err := suite.hdl.SignInAppBySmsCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSmsCodeIsNull
func (suite *SignInAppBySmsCodeTestSuite) TestSmsCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInAppBySmsCodeRequest{
		Phone:   phone,
		SmsCode: smsCodeIsNull,
	}
	resp := &userpb.SignInAppBySmsCodeResponse{}
	err := suite.hdl.SignInAppBySmsCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSmsCodeIsErr
func (suite *SignInAppBySmsCodeTestSuite) TestSmsCodeIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInAppBySmsCodeRequest{
		Phone:   phone,
		SmsCode: smsCodeIsErr,
	}
	resp := &userpb.SignInAppBySmsCodeResponse{}
	err := suite.hdl.SignInAppBySmsCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SignInAppBySmsCodeTestSuite) TearDownSuite() {

}

func TestSignInAppBySmsCodeTestSuite(t *testing.T) {
	suite.Run(t, new(SignInAppBySmsCodeTestSuite))
}
