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

type ResetPasswordTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *ResetPasswordTestSuite) SetupSuite() {
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
func (suite *ResetPasswordTestSuite) SetupTest() {
	suite.notification.On("CheckTxIdUsage", mock.Anything, mock.Anything).Return(
		&notificationpb.CheckTxIdUsageResponse{
			Action: notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD,
		}, nil).After(utils.RpcLatency())
}

func (suite *ResetPasswordTestSuite) TestResetPassword() {
	ctx := context.Background()

	req := &userpb.ResetPasswordRequest{
		Phone:            phone,
		NewPlainPassword: newPlainPassword,
		SmsCode:          smsCode,
		TxId:             txId,
	}
	resp := &userpb.ResetPasswordResponse{}
	err := suite.hdl.ResetPassword(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPhoneIsNull
func (suite *ResetPasswordTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetPasswordRequest{
		Phone:            phoneIsNull,
		NewPlainPassword: newPlainPassword,
		SmsCode:          smsCode,
		TxId:             txId,
	}
	resp := &userpb.ResetPasswordResponse{}
	err := suite.hdl.ResetPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsErr
func (suite *ResetPasswordTestSuite) TestPhoneIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetPasswordRequest{
		Phone:            phoneIsErr,
		NewPlainPassword: newPlainPassword,
		SmsCode:          smsCode,
		TxId:             txId,
	}
	resp := &userpb.ResetPasswordResponse{}
	err := suite.hdl.ResetPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrSystemUserNotFound, err)
}

// TestPlainPasswordIsNull
func (suite *ResetPasswordTestSuite) TestPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetPasswordRequest{
		Phone:            phone,
		NewPlainPassword: plainPasswordIsNull,
		SmsCode:          smsCode,
		TxId:             txId,
	}
	resp := &userpb.ResetPasswordResponse{}
	err := suite.hdl.ResetPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTxIdIsNull
func (suite *ResetPasswordTestSuite) TestTxIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetPasswordRequest{
		Phone:            phone,
		NewPlainPassword: plainPassword,
		SmsCode:          smsCode,
		TxId:             txIdIsNull,
	}
	resp := &userpb.ResetPasswordResponse{}
	err := suite.hdl.ResetPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTxIdIsErr
func (suite *ResetPasswordTestSuite) TestTxIdIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetPasswordRequest{
		Phone:            phone,
		NewPlainPassword: plainPassword,
		SmsCode:          smsCode,
		TxId:             txIdIsNull,
	}
	resp := &userpb.ResetPasswordResponse{}
	err := suite.hdl.ResetPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *ResetPasswordTestSuite) TearDownSuite() {
}

func TestResetPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(ResetPasswordTestSuite))
}
