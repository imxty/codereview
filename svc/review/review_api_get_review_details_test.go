package review

import (
	"context"
	"testing"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	product "github.com/jinmukeji/huimaibao-service/svc/product/mock"
	productmock "github.com/jinmukeji/huimaibao-service/svc/product/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	user "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	usermock "github.com/jinmukeji/huimaibao-service/svc/user/mock"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetReviewDetailsTestSuite struct {
	suite.Suite
	user    *user.UserAPIService
	product *product.ProductAPIService
	hdl     *ReviewAPIHandler
}

func (suite *GetReviewDetailsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reviewFile)
	s := store.NewReviewStore(conn)
	user := &usermock.UserAPIService{}
	suite.user = user
	product := &productmock.ProductAPIService{}
	suite.product = product
	suite.hdl = NewReviewAPIHandler(s, product, user, nil)
}

// TestGetReviewDetails
func (suite *GetReviewDetailsTestSuite) TestGetReviewDetails() {
	ctx := context.Background()
	suite.user.On("GetEntityRevision", mock.Anything, mock.Anything).Return(
		&userpb.GetEntityRevisionResponse{
			Entity: &userpb.TenantEntity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())
	req := &pb.GetReviewDetailsRequest{
		ReviewIssueId: reviewIssueIdCertificateReview,
	}
	resp := &pb.GetReviewDetailsResponse{}
	err := suite.hdl.GetReviewDetails(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestReviewIssueIdIsNull
func (suite *GetReviewDetailsTestSuite) TestReviewIssueIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &pb.GetReviewDetailsRequest{
		ReviewIssueId: reviewIssueIdIsNull,
	}
	resp := &pb.GetReviewDetailsResponse{}
	err := suite.hdl.GetReviewDetails(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetReviewDetailsTestSuite) TearDownSuite() {}

func TestGetReviewDetailsTestSuite(t *testing.T) {
	suite.Run(t, new(GetReviewDetailsTestSuite))
}
