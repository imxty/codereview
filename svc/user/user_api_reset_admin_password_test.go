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

type ResetAdminPasswordTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *ResetAdminPasswordTestSuite) SetupSuite() {
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

func (suite *ResetAdminPasswordTestSuite) SetupTest() {
	suite.notification.On("CheckTxIdUsage", mock.Anything, mock.Anything).Return(
		&notificationpb.CheckTxIdUsageResponse{
			Action: notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD,
		}, nil).After(utils.RpcLatency())
}

func (suite *ResetAdminPasswordTestSuite) TestResetAdminPassword() {
	ctx := context.Background()

	req := &userpb.ResetAdminPasswordRequest{
		Phone:            phone,
		NewPlainPassword: newPlainPassword,
		SmsTxId:          smsTxId,
	}
	resp := &userpb.ResetAdminPasswordResponse{}
	err := suite.hdl.ResetAdminPassword(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}
func (suite *ResetAdminPasswordTestSuite) TestTxidIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetAdminPasswordRequest{
		Phone:            phone,
		NewPlainPassword: newPlainPassword,
		SmsTxId:          txIdIsNull,
	}
	resp := &userpb.ResetAdminPasswordResponse{}
	err := suite.hdl.ResetAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsNull
func (suite *ResetAdminPasswordTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetAdminPasswordRequest{
		Phone:            phoneIsNull,
		NewPlainPassword: newPlainPassword,
		SmsTxId:          smsTxId,
	}
	resp := &userpb.ResetAdminPasswordResponse{}
	err := suite.hdl.ResetAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsErr
func (suite *ResetAdminPasswordTestSuite) TestPhoneIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetAdminPasswordRequest{
		Phone:            phoneIsErr,
		NewPlainPassword: newPlainPassword,
		SmsTxId:          smsTxId,
	}
	resp := &userpb.ResetAdminPasswordResponse{}
	err := suite.hdl.ResetAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

// TestPhoneNotExist
func (suite *ResetAdminPasswordTestSuite) TestPhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	suite.notification.On("CheckTxIdUsage", mock.Anything, mock.Anything).Return(
		&notificationpb.CheckTxIdUsageResponse{
			Action: notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD,
		}, nil).After(utils.RpcLatency())
	req := &userpb.ResetAdminPasswordRequest{
		Phone:            phoneNotExist,
		NewPlainPassword: newPlainPassword,
		SmsTxId:          smsTxId,
	}
	resp := &userpb.ResetAdminPasswordResponse{}
	err := suite.hdl.ResetAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

// TestPlainPasswordIsNull
func (suite *ResetAdminPasswordTestSuite) TestPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetAdminPasswordRequest{
		Phone:            phone,
		NewPlainPassword: newPlainPasswordIsNull,
		SmsTxId:          smsTxId,
	}
	resp := &userpb.ResetAdminPasswordResponse{}
	err := suite.hdl.ResetAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSmsTxIdIsNull
func (suite *ResetAdminPasswordTestSuite) TestSmsTxIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ResetAdminPasswordRequest{
		Phone:            phone,
		NewPlainPassword: plainPassword,
		SmsTxId:          txIdIsNull,
	}
	resp := &userpb.ResetAdminPasswordResponse{}
	err := suite.hdl.ResetAdminPassword(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ResetAdminPasswordTestSuite) TearDownSuite() {
}

func TestResetAdminPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(ResetAdminPasswordTestSuite))
}
