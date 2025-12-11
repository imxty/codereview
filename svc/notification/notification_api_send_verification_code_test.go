package notification

import (
	"context"
	"testing"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/notification"
	alimock "github.com/jinmukeji/huimaibao-service/svc/notification/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SendVerificationCodeTestSuite struct {
	suite.Suite
	hdl    *NotificationAPIHandler
	aliyun *alimock.SmsSend
}

func (suite *SendVerificationCodeTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, notificationFile)
	s := store.NewNotificationStore(conn)
	aliyun := &alimock.SmsSend{}
	suite.aliyun = aliyun
	suite.hdl = NewNotificationAPIHandler(s, aliyun, "")
}

func (suite *SendVerificationCodeTestSuite) SetupTest() {
	suite.aliyun.On("SendSms", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		true, nil).After(utils.RpcLatency())

}

// TestOrganizationSignUP 组织注册
func (suite *SendVerificationCodeTestSuite) TestOrganizationSignUP() {
	ctx := context.Background()

	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION,
		Language:       notificationpb.Language_LANGUAGE_ENGLISH,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationSignIn 组织登陆
func (suite *SendVerificationCodeTestSuite) TestOrganizationSignIn() {
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_ORGANIZATION,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationResetPassword 组织重置密码
func (suite *SendVerificationCodeTestSuite) TestOrganizationResetPassword() {
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestSendVerificationCodeBindPhone  组织绑定手机号
func (suite *SendVerificationCodeTestSuite) TestSendVerificationCodeBindPhone() {
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationBilling   组织账单
func (suite *SendVerificationCodeTestSuite) TestOrganizationBilling() {
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestCreateTenant   商户创建成功
func (suite *SendVerificationCodeTestSuite) TestCreateTenant() {
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestCreateTenantFailed
func (suite *SendVerificationCodeTestSuite) TestCreateTenantFailed() {
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantBindPhone  商户绑定手机号
func (suite *SendVerificationCodeTestSuite) TestTenantBindPhone() {
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantSignIn
func (suite *SendVerificationCodeTestSuite) TestTenantSignIn() {
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNIN_TENANT,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationResetPhone
func (suite *SendVerificationCodeTestSuite) TestOrganizationResetPhone() {
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTemplateActionIsInvalid
func (suite *SendVerificationCodeTestSuite) TestTemplateActionIsInvalid() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_INVALID,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTemplateActionIsUnset
func (suite *SendVerificationCodeTestSuite) TestTemplateActionIsUnset() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_UNSET,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsNull
func (suite *SendVerificationCodeTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phoneIsNull,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE,
		Language:       notificationpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestLanguageIsInvalid
func (suite *SendVerificationCodeTestSuite) TestLanguageIsInvalid() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE,
		Language:       notificationpb.Language_LANGUAGE_INVALID,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestLanguageIsUnset
func (suite *SendVerificationCodeTestSuite) TestLanguageIsUnset() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.SendVerificationCodeRequest{
		Phone:          phone,
		TemplateAction: notificationpb.TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE,
		Language:       notificationpb.Language_LANGUAGE_UNSET,
		Params:         make(map[string]string),
	}
	resp := &notificationpb.SendVerificationCodeResponse{}
	err := suite.hdl.SendVerificationCode(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SendVerificationCodeTestSuite) TearDownSuite() {
}

func TestSendVerificationCodeTestSuite(t *testing.T) {
	suite.Run(t, new(SendVerificationCodeTestSuite))
}
