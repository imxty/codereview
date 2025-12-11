package review

import (
	"context"
	"testing"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	user "github.com/jinmukeji/huimaibao-service/svc/user/mock"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommitEntityCertificateTestSuite struct {
	suite.Suite
	userAPI *user.UserAPIService
	hdl     *ReviewAPIHandler
}

func (suite *CommitEntityCertificateTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reviewFile)
	s := store.NewReviewStore(conn)
	userAPI := &user.UserAPIService{}
	suite.userAPI = userAPI
	suite.hdl = NewReviewAPIHandler(s, nil, userAPI, nil)
}

// TestCommitEntityCertificate
func (suite *CommitEntityCertificateTestSuite) TestCommitEntityCertificate() {
	ctx := context.Background()
	suite.userAPI.On("GetEntityRevision", mock.Anything, mock.Anything, mock.Anything).
		Return(&userpb.GetEntityRevisionResponse{
			Entity: &userpb.TenantEntity{
				TenantId: "111111",
			},
		}, nil).After(utils.RpcLatency())
	req := &pb.CommitEntityCertificateRequest{
		OrganizationId:         organizationId,
		TenantId:               tenantId1,
		TenantEntityRevisionId: "c0ven8ollrtu4feddps0",
		IsOrganization:         false,
	}
	resp := &pb.CommitEntityCertificateResponse{}
	err := suite.hdl.CommitEntityCertificate(ctx, req, resp)
	suite.Assert().NoError(err)
}

func (suite *CommitEntityCertificateTestSuite) TearDownSuite() {}

func TestCommitEntityCertificateTestSuite(t *testing.T) {
	suite.Run(t, new(CommitEntityCertificateTestSuite))
}
