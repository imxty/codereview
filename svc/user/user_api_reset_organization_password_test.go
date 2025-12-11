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

type ResetOrganizationPasswordTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *ResetOrganizationPasswordTestSuite) SetupSuite() {
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

func (suite *ResetOrganizationPasswordTestSuite) TestResetOrganizationPassword() {
	ctx := context.Background()
	suite.notification.On("CheckTxIdUsage", mock.Anything, mock.Anything).Return(
		&notificationpb.CheckTxIdUsageResponse{
			Action: notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD,
		}, nil).After(utils.RpcLatency())
	req := &userpb.ResetOrganizationPasswordRequest{
		Phone:            phoneIsExist,
		NewPlainPassword: newPlainPassword,
		TxId:             txId,
	}
	resp := &userpb.ResetOrganizationPasswordResponse{}
	err := suite.hdl.ResetOrganizationPassword(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestPhoneIsNull
func (suite *ResetOrganizationPasswordTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetOrganizationPasswordRequest{
		Phone:            phoneIsErr,
		NewPlainPassword: newPlainPassword,
		TxId:             txId,
	}
	resp := &userpb.ResetOrganizationPasswordResponse{}
	err := suite.hdl.ResetOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}

// TestPhoneIsErr
func (suite *ResetOrganizationPasswordTestSuite) TestPhoneIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetOrganizationPasswordRequest{
		Phone:            phoneIsErr,
		NewPlainPassword: newPlainPassword,
		TxId:             txId,
	}
	resp := &userpb.ResetOrganizationPasswordResponse{}
	err := suite.hdl.ResetOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}

// TestPlainPasswordIsNull
func (suite *ResetOrganizationPasswordTestSuite) TestPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetOrganizationPasswordRequest{
		Phone:            phone,
		NewPlainPassword: plainPasswordIsNull,
		TxId:             txId,
	}
	resp := &userpb.ResetOrganizationPasswordResponse{}
	err := suite.hdl.ResetOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPlainPasswordIsErr
func (suite *ResetOrganizationPasswordTestSuite) TestPlainPasswordIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetOrganizationPasswordRequest{
		Phone:            phone,
		NewPlainPassword: plainPasswordIsErr,
		TxId:             txId,
	}
	resp := &userpb.ResetOrganizationPasswordResponse{}
	err := suite.hdl.ResetOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}

// TestPlainPasswordIsSame
func (suite *ResetOrganizationPasswordTestSuite) TestPlainPasswordIsSame() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetOrganizationPasswordRequest{
		Phone:            phone,
		NewPlainPassword: newPlainPasswordIsSame,
		TxId:             txId,
	}
	resp := &userpb.ResetOrganizationPasswordResponse{}
	err := suite.hdl.ResetOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}

// TestTxIdIsNull
func (suite *ResetOrganizationPasswordTestSuite) TestTxIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetOrganizationPasswordRequest{
		Phone:            phone,
		NewPlainPassword: plainPassword,
		TxId:             txIdIsNull,
	}
	resp := &userpb.ResetOrganizationPasswordResponse{}
	err := suite.hdl.ResetOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTxIdIsErr
func (suite *ResetOrganizationPasswordTestSuite) TestTxIdIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetOrganizationPasswordRequest{
		Phone:            phone,
		NewPlainPassword: plainPassword,
		TxId:             txIdIsErr,
	}
	resp := &userpb.ResetOrganizationPasswordResponse{}
	err := suite.hdl.ResetOrganizationPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationNotExist, err)
}
func (suite *ResetOrganizationPasswordTestSuite) TearDownSuite() {
}

func TestResetOrganizationPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(ResetOrganizationPasswordTestSuite))
}
