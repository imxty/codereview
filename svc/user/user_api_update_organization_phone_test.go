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

type UpdateOrganizationPhoneTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *UpdateOrganizationPhoneTestSuite) SetupSuite() {
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

func (suite *UpdateOrganizationPhoneTestSuite) TestUpdateOrganizationPhone() {
	ctx := context.Background()
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())

	req := &userpb.UpdateOrganizationPhoneRequest{
		OrganizationId:  organizationId,
		NewSafePhone:    phoneIsNew1,
		NewPhoneSmsCode: smsCode,
		NewPhoneSmsTxId: txId,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: txId,
	}
	resp := &userpb.UpdateOrganizationPhoneResponse{}
	err := suite.hdl.UpdateOrganizationPhone(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *UpdateOrganizationPhoneTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOrganizationPhoneRequest{
		OrganizationId:  organizationIdIsNull,
		NewSafePhone:    phone,
		NewPhoneSmsCode: smsCode,
		NewPhoneSmsTxId: txId,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: txId,
	}
	resp := &userpb.UpdateOrganizationPhoneResponse{}
	err := suite.hdl.UpdateOrganizationPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOrganizationIdNotExist
func (suite *UpdateOrganizationPhoneTestSuite) TestOrganizationIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOrganizationPhoneRequest{
		OrganizationId:  organizationIdNotExist,
		NewSafePhone:    phone,
		NewPhoneSmsCode: smsCode,
		NewPhoneSmsTxId: txId,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: txId,
	}
	resp := &userpb.UpdateOrganizationPhoneResponse{}
	err := suite.hdl.UpdateOrganizationPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}

// TestTestOldPhoneSmsCodeIsNull
func (suite *UpdateOrganizationPhoneTestSuite) TestTestOldPhoneSmsCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())

	req := &userpb.UpdateOrganizationPhoneRequest{
		OrganizationId:  organizationId,
		NewSafePhone:    phone,
		NewPhoneSmsCode: smsCodeIsNull,
		NewPhoneSmsTxId: txId,
		OldPhoneSmsCode: smsCodeIsNull,
		OldPhoneSmsTxId: txId,
	}
	resp := &userpb.UpdateOrganizationPhoneResponse{}
	err := suite.hdl.UpdateOrganizationPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOldPhoneSmsCodeIsErr
func (suite *UpdateOrganizationPhoneTestSuite) TestOldPhoneSmsCodeIsErr() {
	t := suite.T()
	ctx := context.Background()
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())

	req := &userpb.UpdateOrganizationPhoneRequest{
		OrganizationId:  organizationId,
		NewSafePhone:    phone,
		NewPhoneSmsCode: smsCode,
		NewPhoneSmsTxId: txId,
		OldPhoneSmsCode: smsCodeIsNull,
		OldPhoneSmsTxId: txId,
	}
	resp := &userpb.UpdateOrganizationPhoneResponse{}
	err := suite.hdl.UpdateOrganizationPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestNewPlainPasswordIsNull
func (suite *UpdateOrganizationPhoneTestSuite) TestNewPhoneSmsCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOrganizationPhoneRequest{
		OrganizationId:  organizationId,
		NewSafePhone:    phone,
		NewPhoneSmsCode: smsCodeIsNull,
		NewPhoneSmsTxId: txId,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: txId,
	}
	resp := &userpb.UpdateOrganizationPhoneResponse{}
	err := suite.hdl.UpdateOrganizationPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestNewSafePhoneIsExist
func (suite *UpdateOrganizationPhoneTestSuite) TestNewSafePhoneIsExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateOrganizationPhoneRequest{
		OrganizationId:  organizationId,
		NewSafePhone:    phoneIsExist,
		NewPhoneSmsCode: smsCodeIsNull,
		NewPhoneSmsTxId: txId,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: txId,
	}
	resp := &userpb.UpdateOrganizationPhoneResponse{}
	err := suite.hdl.UpdateOrganizationPhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *UpdateOrganizationPhoneTestSuite) TearDownSuite() {

}

func TestUpdateOrganizationPhoneTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateOrganizationPhoneTestSuite))
}
