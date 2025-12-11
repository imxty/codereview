package notification

import (
	"context"
	"testing"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/notification"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type VerifyPhoneVerificationCodeTestSuite struct {
	suite.Suite
	hdl *NotificationAPIHandler
}

func (suite *VerifyPhoneVerificationCodeTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, notificationFile)
	s := store.NewNotificationStore(conn)
	suite.hdl = NewNotificationAPIHandler(s, nil, "")
}

func (suite *VerifyPhoneVerificationCodeTestSuite) TestVerifyPhoneVerificationCode() {
	ctx := context.Background()

	req := &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   phone,
		SmsCode: smsCodeNew,
		TxId:    txIdIsNew,
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS,
	}
	resp := &notificationpb.VerifyPhoneVerificationCodeResponse{}
	err := suite.hdl.VerifyPhoneVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}
func (suite *VerifyPhoneVerificationCodeTestSuite) TestVerifyPhoneVerificationCodeErr() {
	ctx := context.Background()
	t := suite.T()
	req := &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   phone,
		SmsCode: smsCodeIsSignUp,
		TxId:    txIdIsSignUp,
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION,
	}
	resp := &notificationpb.VerifyPhoneVerificationCodeResponse{}
	err := suite.hdl.VerifyPhoneVerificationCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrSmsExpired, err)
}

func (suite *VerifyPhoneVerificationCodeTestSuite) TestVerifyPhoneCodeErr() {
	ctx := context.Background()
	t := suite.T()
	req := &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   phone,
		SmsCode: smsCode,
		TxId:    txIdIsSignUp,
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION,
	}
	resp := &notificationpb.VerifyPhoneVerificationCodeResponse{}
	err := suite.hdl.VerifyPhoneVerificationCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrSmsCodeError, err)
}

func (suite *VerifyPhoneVerificationCodeTestSuite) TearDownSuite() {

}

func TestVerifyPhoneVerificationCodeTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyPhoneVerificationCodeTestSuite))
}
