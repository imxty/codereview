package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	notificationmock "github.com/jinmukeji/huimaibao-service/svc/notification/mock"
	reviewmock "github.com/jinmukeji/huimaibao-service/svc/review/mock"
	s3mock "github.com/jinmukeji/huimaibao-service/svc/user/mock"

	utils "github.com/jinmukeji/huimaibao-service/svc/testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateTenantTestSuite struct {
	suite.Suite
	hdl          *UserAPIHandler
	s3           *s3mock.FileStore
	notification *notificationmock.NotificationAPIService
	review       *reviewmock.ReviewAPIService
}

func (suite *CreateTenantTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	s3 := &s3mock.FileStore{}
	suite.s3 = s3
	notification := &notificationmock.NotificationAPIService{}
	suite.notification = notification
	review := &reviewmock.ReviewAPIService{}
	suite.review = review
	suite.hdl = NewUserAPIHandler(s, notification, review, nil, nil, nil, s3)
}

func (suite *CreateTenantTestSuite) TestCreateTenant() {
	ctx := context.Background()
	image, _ := utils.ReadFile(filePath)
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return("", nil).After(utils.RpcLatency())
	suite.notification.On("VerifyPhoneVerificationCode", mock.Anything, mock.Anything).Return(
		&notificationpb.VerifyPhoneVerificationCodeResponse{}, nil).After(utils.RpcLatency())
	suite.review.On("CommitEntityCertificate", mock.Anything, mock.Anything).Return(
		&reviewpb.CommitEntityCertificateResponse{}, nil).After(utils.RpcLatency())
	req := &userpb.CreateTenantRequest{
		OrganizationId: organizationId,
		Tenant: &userpb.TenantEntity{
			Name: entityName,
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

			ContactPhone: phoneIsNew,
			SafePhone:    phoneIsNew,
		},
		TxId:          txId,
		PlainPassword: plainPassword,
		SmsCode:       smsCode,
	}
	resp := &userpb.CreateTenantResponse{}
	err := suite.hdl.CreateTenant(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *CreateTenantTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	image, _ := utils.ReadFile(filePath)
	req := &userpb.CreateTenantRequest{
		OrganizationId: organizationIdIsNull,
		Tenant: &userpb.TenantEntity{
			TenantId: tenantIdIsNull,
			Name:     entityName,
			Address: &userpb.Address{
				Province: province,
				City:     city,
				District: district,
				Street:   street,
			},
			ContactPhone: phone,
			BusinessLicense: &userpb.UploadingImage{
				Mime:     mime,
				Image:    image,
				Filename: filename,
			},
			BusinessLicenseUrl: businessLicenseUrl,
			SocialCreditCode:   socialCreditCode,
			TenantStatus:       userpb.TenantStatus_TENANT_STATUS_PENDING,
			FailReason:         failReason,
		},
		TxId:          txId,
		PlainPassword: plainPassword,
		SmsCode:       smsCode,
	}
	resp := &userpb.CreateTenantResponse{}
	err := suite.hdl.CreateTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *CreateTenantTestSuite) TearDownSuite() {

}

func TestCreateTenantTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTenantTestSuite))
}
