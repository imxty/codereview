package review

import (
	"context"
	"testing"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	product "github.com/jinmukeji/huimaibao-service/svc/product/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	user "github.com/jinmukeji/huimaibao-service/svc/user/mock"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ListReviewIssueTestSuite struct {
	suite.Suite
	userAPI    *user.UserAPIService
	productAPI *product.ProductAPIService
	hdl        *ReviewAPIHandler
}

func (suite *ListReviewIssueTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reviewFile)
	s := store.NewReviewStore(conn)
	userAPI := &user.UserAPIService{}
	suite.userAPI = userAPI
	productAPI := &product.ProductAPIService{}
	suite.productAPI = productAPI
	suite.hdl = NewReviewAPIHandler(s, productAPI, userAPI, nil)
}

// TestListReviewIssues 正常流程
func (suite *ListReviewIssueTestSuite) TestListReviewIssues() {
	ctx := context.Background()
	suite.userAPI.On("ListSystemUsersByIds", mock.Anything, mock.Anything).Return(
		&userpb.ListSystemUsersByIdsResponse{
			SystemUsers: make(map[string]*userpb.SystemUser),
		}, nil).After(utils.RpcLatency())
	suite.userAPI.On("GetOrganizationsName", mock.Anything, mock.Anything).Return(
		&userpb.GetOrganizationsNameResponse{
			OrganizationName: make(map[string]*userpb.Organization),
		}, nil).After(utils.RpcLatency())
	req := &pb.ListReviewIssuesRequest{
		Pagination: &pb.Pagination{
			Size:   10,
			Offset: 0,
		},
		SearchAll: true,
	}
	resp := &pb.ListReviewIssuesResponse{}
	err := suite.hdl.ListReviewIssues(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(9, len(resp.GetReviewIssues()))
}

// TestListReviewIssuesPending
func (suite *ListReviewIssueTestSuite) TestListReviewIssuesPending() {
	ctx := context.Background()
	req := &pb.ListReviewIssuesRequest{
		Pagination: &pb.Pagination{
			Size:   10,
			Offset: 0,
		},
		Status: pb.ReviewStatus_REVIEW_STATUS_PENDING_REVIEW,
	}
	resp := &pb.ListReviewIssuesResponse{}
	err := suite.hdl.ListReviewIssues(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(2, len(resp.GetReviewIssues()))
}

func (suite *ListReviewIssueTestSuite) TearDownSuite() {}

func TestListReviewIssueTestSuite(t *testing.T) {
	suite.Run(t, new(ListReviewIssueTestSuite))
}
