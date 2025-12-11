package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	notificationmock "github.com/jinmukeji/huimaibao-service/svc/notification/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UpdateSafePhoneTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *UpdateSafePhoneTestSuite) SetupSuite() {
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

func (suite *UpdateSafePhoneTestSuite) TestUpdateSafePhone() {
	ctx := context.Background()
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())

	req := &userpb.UpdateSafePhoneRequest{
		TenantId:        tenantId,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: smsTxId,
		NewSafePhone:    phoneIsNew1,
		NewPhoneSmsCode: smsCode,
		NewPhoneSmsTxId: smsTxId,
	}
	resp := &userpb.UpdateSafePhoneResponse{}
	err := suite.hdl.UpdateSafePhone(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *UpdateSafePhoneTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateSafePhoneRequest{
		TenantId:        tenantIdIsNull,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: smsTxId,
		NewSafePhone:    phone,
		NewPhoneSmsCode: smsCode,
		NewPhoneSmsTxId: smsTxId,
	}
	resp := &userpb.UpdateSafePhoneResponse{}
	err := suite.hdl.UpdateSafePhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOldSmsCodeIsNull
func (suite *UpdateSafePhoneTestSuite) TestOldSmsCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateSafePhoneRequest{
		TenantId:        tenantId,
		OldPhoneSmsCode: smsCodeIsNull,
		OldPhoneSmsTxId: smsTxId,
		NewSafePhone:    phone,
		NewPhoneSmsCode: smsCode,
		NewPhoneSmsTxId: smsTxId,
	}
	resp := &userpb.UpdateSafePhoneResponse{}
	err := suite.hdl.UpdateSafePhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOldSmsTxIdIsNull

func (suite *UpdateSafePhoneTestSuite) TestOldSmsTxIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateSafePhoneRequest{
		TenantId:        tenantId,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: txIdIsNull,
		NewSafePhone:    phone,
		NewPhoneSmsCode: smsCode,
		NewPhoneSmsTxId: smsTxId,
	}
	resp := &userpb.UpdateSafePhoneResponse{}
	err := suite.hdl.UpdateSafePhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsNull
func (suite *UpdateSafePhoneTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateSafePhoneRequest{
		TenantId:        tenantId,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: smsTxId,
		NewSafePhone:    phoneIsNull,
		NewPhoneSmsCode: smsCode,
		NewPhoneSmsTxId: smsTxId,
	}
	resp := &userpb.UpdateSafePhoneResponse{}
	err := suite.hdl.UpdateSafePhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestNewSmsCodeIsNull
func (suite *UpdateSafePhoneTestSuite) TestNewSmsCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateSafePhoneRequest{
		TenantId:        tenantId,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: smsTxId,
		NewSafePhone:    phone,
		NewPhoneSmsCode: smsCodeIsNull,
		NewPhoneSmsTxId: smsTxId,
	}
	resp := &userpb.UpdateSafePhoneResponse{}
	err := suite.hdl.UpdateSafePhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestNewSmsTxIdIsNull

func (suite *UpdateSafePhoneTestSuite) TestNewSmsTxIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateSafePhoneRequest{
		TenantId:        tenantId,
		OldPhoneSmsCode: smsCode,
		OldPhoneSmsTxId: smsTxId,
		NewSafePhone:    phone,
		NewPhoneSmsCode: smsCode,
		NewPhoneSmsTxId: txIdIsNull,
	}
	resp := &userpb.UpdateSafePhoneResponse{}
	err := suite.hdl.UpdateSafePhone(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *UpdateSafePhoneTestSuite) TearDownSuite() {

}

func TestUpdateSafePhoneTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateSafePhoneTestSuite))
}
