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

type ResetAppStaffPasswordTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *ResetAppStaffPasswordTestSuite) SetupSuite() {
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

func (suite *ResetAppStaffPasswordTestSuite) SetupTest() {
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())
}

func (suite *ResetAppStaffPasswordTestSuite) TestResetAppStaffPassword() {
	ctx := context.Background()

	req := &userpb.ResetAppStaffPasswordRequest{
		Phone:            staffPhoneIsExist,
		NewPlainPassword: newPlainPassword,
		SmsCode:          smsCode,
		TxId:             txId,
	}
	resp := &userpb.ResetAppStaffPasswordResponse{}
	err := suite.hdl.ResetAppStaffPassword(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}
func (suite *ResetAppStaffPasswordTestSuite) TestTxidIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetAppStaffPasswordRequest{
		Phone:            phone,
		NewPlainPassword: newPlainPassword,
		SmsCode:          smsCode,
		TxId:             txIdIsNull,
	}
	resp := &userpb.ResetAppStaffPasswordResponse{}
	err := suite.hdl.ResetAppStaffPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsNull
func (suite *ResetAppStaffPasswordTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetAppStaffPasswordRequest{
		Phone:            phoneIsNull,
		NewPlainPassword: newPlainPassword,
		SmsCode:          smsCode,
		TxId:             smsTxId,
	}
	resp := &userpb.ResetAppStaffPasswordResponse{}
	err := suite.hdl.ResetAppStaffPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsErr
func (suite *ResetAppStaffPasswordTestSuite) TestPhoneIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetAppStaffPasswordRequest{
		Phone:            phoneIsErr,
		NewPlainPassword: newPlainPassword,
		SmsCode:          smsCode,
		TxId:             smsTxId,
	}
	resp := &userpb.ResetAppStaffPasswordResponse{}
	err := suite.hdl.ResetAppStaffPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

// TestPhoneNotExist
func (suite *ResetAppStaffPasswordTestSuite) TestPhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	suite.notification.On("CheckTxIdUsage", mock.Anything, mock.Anything).Return(
		&notificationpb.CheckTxIdUsageResponse{
			Action: notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD,
		}, nil).After(utils.RpcLatency())
	req := &userpb.ResetAppStaffPasswordRequest{
		Phone:            phoneNotExist,
		NewPlainPassword: newPlainPassword,
		SmsCode:          smsCode,
		TxId:             smsTxId,
	}
	resp := &userpb.ResetAppStaffPasswordResponse{}
	err := suite.hdl.ResetAppStaffPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrStaffNotExist, err)
}

// TestPlainPasswordIsNull
func (suite *ResetAppStaffPasswordTestSuite) TestPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetAppStaffPasswordRequest{
		Phone:            phone,
		NewPlainPassword: newPlainPasswordIsNull,
		SmsCode:          smsCode,
		TxId:             smsTxId,
	}
	resp := &userpb.ResetAppStaffPasswordResponse{}
	err := suite.hdl.ResetAppStaffPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSmsCodeIsNull
func (suite *ResetAppStaffPasswordTestSuite) TestSmsCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetAppStaffPasswordRequest{
		Phone:            phone,
		NewPlainPassword: plainPassword,
		SmsCode:          smsCodeIsNull,
		TxId:             txId,
	}
	resp := &userpb.ResetAppStaffPasswordResponse{}
	err := suite.hdl.ResetAppStaffPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ResetAppStaffPasswordTestSuite) TearDownSuite() {
}

func TestResetAppStaffPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(ResetAppStaffPasswordTestSuite))
}
