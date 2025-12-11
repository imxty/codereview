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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SignUpOrganizationTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *SignUpOrganizationTestSuite) SetupSuite() {
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

func (suite *SignUpOrganizationTestSuite) TestSignUpOrganization() {
	ctx := context.Background()
	suite.notification.On("GetLastestVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.GetLatestPhoneVerificationCodeResponse{
			SmsCode: smsCode,
			TxId:    txId,
		}, nil).After(utils.RpcLatency())

	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())

	req := &userpb.SignUpOrganizationRequest{
		Username:      usernameIsNew,
		PlainPassword: plainPassword,
		Phone:         phoneIsNew1,
		SmsCode:       smsCode,
		TxId:          txId,
	}
	resp := &userpb.SignUpOrganizationResponse{}
	err := suite.hdl.SignUpOrganization(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestUsernameIsNull
func (suite *SignUpOrganizationTestSuite) TestUsernameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignUpOrganizationRequest{
		Username:      usernameIsNull,
		PlainPassword: plainPassword,
		Phone:         phoneIsNew,
		SmsCode:       smsCode,
		TxId:          txId,
	}
	resp := &userpb.SignUpOrganizationResponse{}
	err := suite.hdl.SignUpOrganization(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPlainPasswordIsNull
func (suite *SignUpOrganizationTestSuite) TestPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignUpOrganizationRequest{
		Username:      usernameIsNew,
		PlainPassword: plainPasswordIsNull,
		Phone:         phoneIsNew,
		SmsCode:       smsCode,
		TxId:          txId,
	}
	resp := &userpb.SignUpOrganizationResponse{}
	err := suite.hdl.SignUpOrganization(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsNull
func (suite *SignUpOrganizationTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignUpOrganizationRequest{
		Username:      usernameIsNew,
		PlainPassword: plainPassword,
		Phone:         phoneIsNull,
		SmsCode:       smsCode,
		TxId:          txId,
	}
	resp := &userpb.SignUpOrganizationResponse{}
	err := suite.hdl.SignUpOrganization(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSmsCodeIsNull
func (suite *SignUpOrganizationTestSuite) TestSmsCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignUpOrganizationRequest{
		Username:      usernameIsNew,
		PlainPassword: plainPassword,
		Phone:         phoneIsNew,
		SmsCode:       smsCodeIsNull,
		TxId:          txId,
	}
	resp := &userpb.SignUpOrganizationResponse{}
	err := suite.hdl.SignUpOrganization(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSmsCodeIsErr
func (suite *SignUpOrganizationTestSuite) TestSmsCodeIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SignUpOrganizationRequest{
		Username:      usernameIsNew,
		PlainPassword: plainPassword,
		Phone:         phoneIsNew,
		SmsCode:       smsCodeIsErr,
		TxId:          txId,
	}
	resp := &userpb.SignUpOrganizationResponse{}
	err := suite.hdl.SignUpOrganization(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *SignUpOrganizationTestSuite) TearDownSuite() {
}

func TestSignUpOrganizationTestSuite(t *testing.T) {
	suite.Run(t, new(SignUpOrganizationTestSuite))
}
