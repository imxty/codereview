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

type SignInTenantAdminByPhoneTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *SignInTenantAdminByPhoneTestSuite) SetupTest() {
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

// TestSignInTenantAdminByPhone 通过手机号登陆组织管理页面请求
func (suite *SignInTenantAdminByPhoneTestSuite) TestSignInTenantAdminByPhone() {
	ctx := context.Background()
	suite.notification.On("GetLastestVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.GetLatestPhoneVerificationCodeResponse{
			SmsCode: smsCode,
			TxId:    txId,
		}, nil).After(utils.RpcLatency())
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())

	req := &userpb.SignInTenantAdminByPhoneRequest{
		Phone:   phone,
		SmsCode: smsCode,
		TxId:    txId,
	}
	resp := &userpb.SignInTenantAdminByPhoneResponse{}
	err := suite.hdl.SignInTenantAdminByPhone(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPhoneIsNull
func (suite *SignInTenantAdminByPhoneTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInTenantAdminByPhoneRequest{
		Phone:   phoneIsNull,
		SmsCode: smsCode,
		TxId:    txId,
	}
	resp := &userpb.SignInTenantAdminByPhoneResponse{}
	err := suite.hdl.SignInTenantAdminByPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsErr
func (suite *SignInTenantAdminByPhoneTestSuite) TestPhoneIsErr() {
	t := suite.T()
	ctx := context.Background()
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())

	req := &userpb.SignInTenantAdminByPhoneRequest{
		Phone:   phoneIsErr,
		SmsCode: smsCode,
		TxId:    txId,
	}
	resp := &userpb.SignInTenantAdminByPhoneResponse{}
	err := suite.hdl.SignInTenantAdminByPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

// TestPhoneNotExist
func (suite *SignInTenantAdminByPhoneTestSuite) TestPhoneNotExist() {
	t := suite.T()
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())

	ctx := context.Background()
	req := &userpb.SignInTenantAdminByPhoneRequest{
		Phone:   phoneNotExist,
		SmsCode: smsCode,
		TxId:    txId,
	}
	resp := &userpb.SignInTenantAdminByPhoneResponse{}
	err := suite.hdl.SignInTenantAdminByPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

// TestSmsCodeIsNull
func (suite *SignInTenantAdminByPhoneTestSuite) TestSmsCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInTenantAdminByPhoneRequest{
		Phone:   phone,
		SmsCode: smsCodeIsNull,
		TxId:    txId,
	}
	resp := &userpb.SignInTenantAdminByPhoneResponse{}
	err := suite.hdl.SignInTenantAdminByPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSmsCodeIsErr
func (suite *SignInTenantAdminByPhoneTestSuite) TestSmsCodeIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignInTenantAdminByPhoneRequest{
		Phone:   phone,
		SmsCode: smsCodeIsErr,
		TxId:    txId,
	}
	resp := &userpb.SignInTenantAdminByPhoneResponse{}
	err := suite.hdl.SignInTenantAdminByPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SignInTenantAdminByPhoneTestSuite) TearDownSuite() {

}

func TestSignInTenantAdminByPhoneTestSuite(t *testing.T) {
	suite.Run(t, new(SignInTenantAdminByPhoneTestSuite))
}
