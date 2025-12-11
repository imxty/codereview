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

type ReCreateTenantTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	notification *notificationmock.NotificationAPIService
}

func (suite *ReCreateTenantTestSuite) SetupSuite() {
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

func (suite *ReCreateTenantTestSuite) TestReCreateTenant() {
	ctx := context.Background()
	image, _ := utils.ReadFile(filePath)
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())

	req := &userpb.ReCreateTenantRequest{
		Tenant: &userpb.TenantEntity{
			TenantId:       tenantId,
			OrganizationId: organizationId,
			ContactPhone:   phone,
			Name:           tenantName,
			Address: &userpb.Address{
				Province: province,
				City:     city,
				District: district,
				Street:   street,
			},
			ContactName: contactName,
			BusinessLicense: &userpb.UploadingImage{
				Mime:     mime,
				Image:    image,
				Filename: filename,
			},
			BusinessLicenseUrl: businessLicenseUrl,
			SocialCreditCode:   socialCreditCode,
			SafePhone:          phoneIsNew,
		},
		PlainPassword: password,
		SmsCode:       smsCode,
		TxId:          txId,
	}
	resp := &userpb.ReCreateTenantResponse{}
	err := suite.hdl.ReCreateTenant(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestBusinessLicenseIsNull
func (suite *ReCreateTenantTestSuite) TestBusinessLicenseIsNull() {
	t := suite.T()
	ctx := context.Background()
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())

	req := &userpb.ReCreateTenantRequest{
		Tenant: &userpb.TenantEntity{
			TenantId:       tenantId,
			OrganizationId: organizationId,
			ContactPhone:   phone,
			Name:           tenantName,
			Address: &userpb.Address{
				Province: province,
				City:     city,
				District: district,
				Street:   street,
			},
			ContactName:        contactName,
			BusinessLicense:    nil,
			BusinessLicenseUrl: businessLicenseUrl,
			SocialCreditCode:   socialCreditCode,
			SafePhone:          phoneIsNew,
		},
		PlainPassword: password,
		SmsCode:       smsCode,
		TxId:          txId,
	}
	resp := &userpb.ReCreateTenantResponse{}
	err := suite.hdl.ReCreateTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrOrganizationTenantNotExist, err)
}

// TestAddressIsErr
func (suite *ReCreateTenantTestSuite) TestAddressIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ReCreateTenantRequest{
		Tenant: &userpb.TenantEntity{
			TenantId:           tenantId,
			OrganizationId:     organizationId,
			ContactPhone:       phone,
			Name:               tenantName,
			Address:            nil,
			ContactName:        contactName,
			BusinessLicense:    nil,
			BusinessLicenseUrl: businessLicenseUrl,
			SocialCreditCode:   socialCreditCode,
			SafePhone:          phoneIsNew,
		},
		PlainPassword: password,
		SmsCode:       smsCode,
		TxId:          txId,
	}
	resp := &userpb.ReCreateTenantResponse{}
	err := suite.hdl.ReCreateTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdIsNull
func (suite *ReCreateTenantTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ReCreateTenantRequest{
		Tenant: &userpb.TenantEntity{
			TenantId:       tenantIdIsNull,
			OrganizationId: organizationIdIsNull,
			ContactPhone:   phone,
			Name:           tenantName,
			Address: &userpb.Address{
				Province: province,
				City:     city,
				District: district,
				Street:   street,
			},
			ContactName:        contactName,
			BusinessLicense:    nil,
			BusinessLicenseUrl: businessLicenseUrl,
			SocialCreditCode:   socialCreditCode,
			SafePhone:          phoneIsNew,
		},
		PlainPassword: password,
		SmsCode:       smsCode,
		TxId:          txId,
	}
	resp := &userpb.ReCreateTenantResponse{}
	err := suite.hdl.ReCreateTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPlainPasswordIsNull
func (suite *ReCreateTenantTestSuite) TestPlainPasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ReCreateTenantRequest{
		Tenant: &userpb.TenantEntity{
			TenantId:       tenantIdIsNull,
			OrganizationId: organizationIdIsNull,
			ContactPhone:   phone,
			Name:           tenantName,
			Address: &userpb.Address{
				Province: province,
				City:     city,
				District: district,
				Street:   street,
			},
			ContactName:        contactName,
			BusinessLicense:    nil,
			BusinessLicenseUrl: businessLicenseUrl,
			SocialCreditCode:   socialCreditCode,
			SafePhone:          phoneIsNew,
		},
		PlainPassword: passwordIsNull,
		SmsCode:       smsCode,
		TxId:          txId,
	}
	resp := &userpb.ReCreateTenantResponse{}
	err := suite.hdl.ReCreateTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestSmsTxIdIsNull
func (suite *ReCreateTenantTestSuite) TestSmsTxIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ReCreateTenantRequest{
		Tenant: &userpb.TenantEntity{
			TenantId:       tenantIdIsNull,
			OrganizationId: organizationIdIsNull,
			ContactPhone:   phone,
			Name:           tenantName,
			Address: &userpb.Address{
				Province: province,
				City:     city,
				District: district,
				Street:   street,
			},
			ContactName:        contactName,
			BusinessLicense:    nil,
			BusinessLicenseUrl: businessLicenseUrl,
			SocialCreditCode:   socialCreditCode,
			SafePhone:          phoneIsNew,
		},
		PlainPassword: password,
		SmsCode:       smsCode,
		TxId:          txIdIsNull,
	}
	resp := &userpb.ReCreateTenantResponse{}
	err := suite.hdl.ReCreateTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ReCreateTenantTestSuite) TearDownSuite() {
}

func TestReCreateTenantTestSuite(t *testing.T) {
	suite.Run(t, new(ReCreateTenantTestSuite))
}
