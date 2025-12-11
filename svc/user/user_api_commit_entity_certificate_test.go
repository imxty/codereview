package user

import (
	"context"
	"testing"

	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	reviewmock "github.com/jinmukeji/huimaibao-service/svc/review/mock"
	s3mock "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommitEntityCertificateTestSuite struct {
	suite.Suite
	hdl    *UserAPIHandler
	s3     *s3mock.FileStore
	review *reviewmock.ReviewAPIService
}

func (suite *CommitEntityCertificateTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	s3 := &s3mock.FileStore{}
	suite.s3 = s3
	review := &reviewmock.ReviewAPIService{}
	suite.review = review
	suite.hdl = NewUserAPIHandler(s, nil, review, nil, nil, nil, s3)
}

func (suite *CommitEntityCertificateTestSuite) TestCommitEntityCertificate() {
	ctx := context.Background()
	image, _ := utils.ReadFile(filePath)
	suite.s3.On("Save", mock.Anything, mock.Anything, mock.Anything).Return("", nil).After(utils.RpcLatency())
	suite.review.On("CommitEntityCertificate", mock.Anything, mock.Anything).Return(
		&reviewpb.CommitEntityCertificateResponse{}, nil).After(utils.RpcLatency())
	req := &userpb.CommitEntityCertificateRequest{
		OrganizationId: organizationId,
		Tenant: &userpb.TenantEntity{
			TenantId: tenantId,
			Name:     tenantName,
			Address: &userpb.Address{
				Province: province,
				City:     city,
				District: district,
				Street:   street,
			},
			SafePhone: phone,
			BusinessLicense: &userpb.UploadingImage{
				Mime:     mime,
				Image:    image,
				Filename: filename,
			},
			SocialCreditCode:   socialCreditCode,
			TenantReviewStatus: userpb.TenantReviewStatus_TENANT_REVIEW_STATUS_INVALID,
			FailReason:         failReason,
			TenantStatus:       userpb.TenantStatus_TENANT_STATUS_PENDING,
			BusinessLicenseUrl: businessLicenseUrl,
			ContactPhone:       phone,
		},
		IsOrganization: false,
	}
	resp := &userpb.CommitEntityCertificateResponse{}
	err := suite.hdl.CommitEntityCertificate(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *CommitEntityCertificateTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	image, _ := utils.ReadFile(filePath)
	req := &userpb.CommitEntityCertificateRequest{
		OrganizationId: organizationId,
		Tenant: &userpb.TenantEntity{
			Name: entityName,
			Address: &userpb.Address{
				Province: province,
				City:     city,
				District: district,
				Street:   street,
			},
			SafePhone: phone,
			BusinessLicense: &userpb.UploadingImage{
				Mime:     mime,
				Image:    image,
				Filename: filename,
			},
			SocialCreditCode:   socialCreditCode,
			FailReason:         failReason,
			TenantStatus:       userpb.TenantStatus_TENANT_STATUS_PENDING,
			BusinessLicenseUrl: businessLicenseUrl,
			ContactPhone:       phone,
		},
		IsOrganization: false,
	}
	resp := &userpb.CommitEntityCertificateResponse{}
	err := suite.hdl.CommitEntityCertificate(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *CommitEntityCertificateTestSuite) TearDownSuite() {

}

func TestCommitEntityCertificateTestSuite(t *testing.T) {
	suite.Run(t, new(CommitEntityCertificateTestSuite))
}
