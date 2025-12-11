package review

import (
	"context"
	"testing"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	user "github.com/jinmukeji/huimaibao-service/svc/user/mock"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// 提交证书结果
type CommitCertificateReviewResultTestSuite struct {
	suite.Suite
	userAPI *user.UserAPIService
	hdl     *ReviewAPIHandler
}

func (suite *CommitCertificateReviewResultTestSuite) SetupSuite() {
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

// TestCommitCertificateReviewResult,通过
func (suite *CommitCertificateReviewResultTestSuite) TestCommitCertificateReviewResult() {
	ctx := context.Background()
	suite.userAPI.On("UpdateTenantCertificateByReviewResult", mock.Anything, mock.Anything).Return(
		&userpb.UpdateTenantCertificateByReviewResultResponse{}, nil).After(utils.RpcLatency())
	req := &pb.CommitCertificateReviewResultRequest{
		ReviewIssueId:  reviewIssueIdCertificateReview,
		ReviewerUserId: reviewerUserId,
		Pass:           true,
		Comment:        comment,
	}
	resp := &pb.CommitCertificateReviewResultResponse{}
	err := suite.hdl.CommitCertificateReviewResult(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestCommitCertificateReviewResultFalse,通过
func (suite *CommitCertificateReviewResultTestSuite) TestCommitCertificateReviewResultFalse() {
	ctx := context.Background()
	t := suite.T()
	suite.userAPI.On("GetTenantEntity", mock.Anything, mock.Anything).Return(
		&userpb.GetTenantEntityResponse{
			Entity: &userpb.Entity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())
	req := &pb.CommitCertificateReviewResultRequest{
		ReviewIssueId:  reviewIssueIdCertificateReview,
		ReviewerUserId: reviewerUserId,
		Pass:           false,
		Comment:        comment,
	}
	resp := &pb.CommitCertificateReviewResultResponse{}
	err := suite.hdl.CommitCertificateReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCommitCertificateReviewResultNotPass,不通过
func (suite *CommitCertificateReviewResultTestSuite) TestCommitCertificateReviewResultNotPass() {
	ctx := context.Background()
	t := suite.T()
	suite.userAPI.On("GetTenantEntity", mock.Anything, mock.Anything).Return(
		&userpb.GetTenantEntityResponse{
			Entity: &userpb.Entity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())
	req := &pb.CommitCertificateReviewResultRequest{
		ReviewIssueId:  "c071jr6vvhfsr0viotng",
		ReviewerUserId: "c0vfijbipt3agu1o3jl0",
		Pass:           false,
		Comment:        "ssss",
	}
	resp := &pb.CommitCertificateReviewResultResponse{}
	err := suite.hdl.CommitCertificateReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrReviewIssueNotFound, err)
}

// TestReviewIssueIdIsNull
func (suite *CommitCertificateReviewResultTestSuite) TestReviewIssueIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.CommitCertificateReviewResultRequest{
		ReviewIssueId:  reviewIssueIdIsNull,
		ReviewerUserId: reviewerUserId,
		Pass:           false,
		Comment:        comment,
	}
	resp := &pb.CommitCertificateReviewResultResponse{}
	err := suite.hdl.CommitCertificateReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReviewerUserIdIsNull
func (suite *CommitCertificateReviewResultTestSuite) TestReviewerUserIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.CommitCertificateReviewResultRequest{
		ReviewIssueId:  reviewIssueIdCertificateFailed,
		ReviewerUserId: reviewerUserIdIsNull,
		Pass:           false,
		Comment:        comment,
	}
	resp := &pb.CommitCertificateReviewResultResponse{}
	err := suite.hdl.CommitCertificateReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCommentIsNull
func (suite *CommitCertificateReviewResultTestSuite) TestCommentIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.CommitCertificateReviewResultRequest{
		ReviewIssueId:  reviewIssueIdCertificateFailed,
		ReviewerUserId: reviewerUserId,
		Pass:           false,
		Comment:        commentIsEmpty,
	}
	resp := &pb.CommitCertificateReviewResultResponse{}
	err := suite.hdl.CommitCertificateReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *CommitCertificateReviewResultTestSuite) TearDownSuite() {}

func TestCommitCertificateReviewResultTestSuite(t *testing.T) {
	suite.Run(t, new(CommitCertificateReviewResultTestSuite))
}
