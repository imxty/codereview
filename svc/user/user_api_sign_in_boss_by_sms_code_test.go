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

type SignInBossBySmsCodeTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *SignInBossBySmsCodeTestSuite) SetupTest() {
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

func (suite *SignInBossBySmsCodeTestSuite) TestSignInJMBySmsCode() {
	ctx := context.Background()
	suite.notification.On("GetLastestVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.GetLatestPhoneVerificationCodeResponse{
			SmsCode: smsCode,
			TxId:    txId,
		}, nil).After(utils.RpcLatency())
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())
	req := &userpb.SignInBossBySmsCodeRequest{
		Phone:   phone,
		SmsCode: smsCode,
		TxId:    txId,
	}
	resp := &userpb.SignInBossBySmsCodeResponse{}
	err := suite.hdl.SignInBossBySmsCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *SignInBossBySmsCodeTestSuite) TestSignInBossBySmsCode() {
	ctx := context.Background()
	suite.notification.On("GetLastestVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.GetLatestPhoneVerificationCodeResponse{
			SmsCode: smsCode,
			TxId:    txId,
		}, nil).After(utils.RpcLatency())
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())
	req := &userpb.SignInBossBySmsCodeRequest{
		Phone:   phone1,
		SmsCode: smsCode,
		TxId:    txId,
	}
	resp := &userpb.SignInBossBySmsCodeResponse{}
	err := suite.hdl.SignInBossBySmsCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPhoneIsNull
func (suite *SignInBossBySmsCodeTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInBossBySmsCodeRequest{
		Phone:   phoneIsNull,
		SmsCode: smsCode,
	}
	resp := &userpb.SignInBossBySmsCodeResponse{}
	err := suite.hdl.SignInBossBySmsCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsErr
func (suite *SignInBossBySmsCodeTestSuite) TestPhoneIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInBossBySmsCodeRequest{
		Phone:   phoneIsErr,
		SmsCode: smsCode,
	}
	resp := &userpb.SignInBossBySmsCodeResponse{}
	err := suite.hdl.SignInBossBySmsCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneNotExist
func (suite *SignInBossBySmsCodeTestSuite) TestPhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInBossBySmsCodeRequest{
		Phone:   phoneNotExist,
		SmsCode: smsCode,
	}
	resp := &userpb.SignInBossBySmsCodeResponse{}
	err := suite.hdl.SignInBossBySmsCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSmsCodeIsNull
func (suite *SignInBossBySmsCodeTestSuite) TestSmsCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInBossBySmsCodeRequest{
		Phone:   phone,
		SmsCode: smsCodeIsNull,
	}
	resp := &userpb.SignInBossBySmsCodeResponse{}
	err := suite.hdl.SignInBossBySmsCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSmsCodeIsErr
func (suite *SignInBossBySmsCodeTestSuite) TestSmsCodeIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInBossBySmsCodeRequest{
		Phone:   phone,
		SmsCode: smsCodeIsErr,
	}
	resp := &userpb.SignInBossBySmsCodeResponse{}
	err := suite.hdl.SignInBossBySmsCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SignInBossBySmsCodeTestSuite) TearDownSuite() {

}

func TestSignInBossBySmsCodeTestSuite(t *testing.T) {
	suite.Run(t, new(SignInBossBySmsCodeTestSuite))
}
