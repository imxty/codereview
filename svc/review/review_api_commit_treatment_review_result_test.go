package review

import (
	"context"
	"testing"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	productmock "github.com/jinmukeji/huimaibao-service/svc/product/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommitTreatmentReviewResultTestSuite struct {
	suite.Suite
	hdl     *ReviewAPIHandler
	product *productmock.ProductAPIService
}

func (suite *CommitTreatmentReviewResultTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reviewFile)
	s := store.NewReviewStore(conn)
	product := &productmock.ProductAPIService{}
	suite.product = product
	suite.hdl = NewReviewAPIHandler(s, product, nil, nil)
}

func (suite *CommitTreatmentReviewResultTestSuite) SetupTest() {
	suite.product.On("UpdateTreatmentStatus", mock.Anything, mock.Anything, mock.Anything).
		Return(&productpb.UpdateTreatmentStatusResponse{}, nil).After(utils.RpcLatency())
	suite.product.On("ChangeProductStatusByTreatmentRevId", mock.Anything, mock.Anything, mock.Anything).
		Return(&productpb.ChangeProductStatusByTreatmentRevIdResponse{}, nil).After(utils.RpcLatency())
	suite.product.On("GetTreatmentRevByTreatmentRevId", mock.Anything, mock.Anything, mock.Anything).
		Return(&productpb.GetTreatmentRevByTreatmentRevIdResponse{
			Treatment: &productpb.Treatment{
				TreatmentName: treatmentName,
			},
		}, nil).After(utils.RpcLatency())
}

// TestCommitTreatmentReviewResult
func (suite *CommitTreatmentReviewResultTestSuite) TestCommitTreatmentReviewResult() {
	ctx := context.Background()

	req := &pb.CommitTreatmentReviewResultRequest{
		ReviewIssueId:  reviewIssueIdCertificateReview,
		ReviewerUserId: reviewerUserId,
		Pass:           true,
		Comment:        reviewComment,
	}
	resp := &pb.CommitTreatmentReviewResultResponse{}
	err := suite.hdl.CommitTreatmentReviewResult(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(false, resp.CancelReview)
}

// TestReviewStatusIsNotReviewing
func (suite *CommitTreatmentReviewResultTestSuite) TestTestReviewStatusIsNotReviewing() {
	t := suite.T()
	ctx := context.Background()

	req := &pb.CommitTreatmentReviewResultRequest{
		ReviewIssueId:  reviewIssueIdCertificateToBeReviewed,
		ReviewerUserId: reviewerUserId,
		Pass:           true,
		Comment:        reviewComment,
	}
	resp := &pb.CommitTreatmentReviewResultResponse{}
	err := suite.hdl.CommitTreatmentReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReviewIssueIsCanceled
func (suite *CommitTreatmentReviewResultTestSuite) TestReviewIssueIsCanceled() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.CommitTreatmentReviewResultRequest{
		ReviewIssueId:  reviewIssueIdCertificateFailed,
		ReviewerUserId: reviewerUserId,
		Pass:           true,
		Comment:        reviewComment,
	}
	resp := &pb.CommitTreatmentReviewResultResponse{}
	err := suite.hdl.CommitTreatmentReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
	suite.Assert().Equal(resp.CancelReview, false)
}

// TestCommentIsEmpty
func (suite *CommitTreatmentReviewResultTestSuite) TestCommentIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &pb.CommitTreatmentReviewResultRequest{
		ReviewIssueId:  reviewIssueIdCertificateReview,
		ReviewerUserId: reviewerUserId,
		Pass:           false,
		Comment:        commentIsEmpty,
	}
	resp := &pb.CommitTreatmentReviewResultResponse{}
	err := suite.hdl.CommitTreatmentReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestNoRemainingChancesOfReviewing
func (suite *CommitTreatmentReviewResultTestSuite) TestNoRemainingChancesOfReviewing() {
	ctx := context.Background()
	req := &pb.CommitTreatmentReviewResultRequest{
		ReviewIssueId:  reviewIssueIdProductReview,
		ReviewerUserId: reviewerUserId,
		Pass:           true,
		Comment:        reviewComment,
	}
	resp := &pb.CommitTreatmentReviewResultResponse{}
	err := suite.hdl.CommitTreatmentReviewResult(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(false, resp.CancelReview)
}

// TestReviewIssueStatusIsNotReviewing
func (suite *CommitTreatmentReviewResultTestSuite) TestReviewIssueStatusIsNotReviewing() {
	t := suite.T()
	ctx := context.Background()
	req := &pb.CommitTreatmentReviewResultRequest{
		ReviewIssueId:  reviewIssueIdCertificateSuccess,
		ReviewerUserId: reviewerUserId,
		Pass:           true,
		Comment:        reviewComment,
	}
	resp := &pb.CommitTreatmentReviewResultResponse{}
	err := suite.hdl.CommitTreatmentReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReviewIssueStatusIsReviewingAndIsFirstReviewAndPass
func (suite *CommitTreatmentReviewResultTestSuite) TestReviewIssueStatusIsReviewingAndIsFirstReviewAndPass() {
	ctx := context.Background()
	req := &pb.CommitTreatmentReviewResultRequest{
		ReviewIssueId:  reviewIssueId,
		ReviewerUserId: reviewerUserId,
		Pass:           true,
		Comment:        comment,
	}
	resp := &pb.CommitTreatmentReviewResultResponse{}
	err := suite.hdl.CommitTreatmentReviewResult(ctx, req, resp)
	suite.Assert().NoError(err)
	// 查询审核状态
	reviewIssue, err := suite.hdl.reviewStore.GetReviewIssueById(ctx, req.GetReviewIssueId())
	suite.Assert().NoError(err)
	suite.Assert().Equal(int32(3), reviewIssue.GetReviewStatus())
	// 查询审核结果
	result, err := suite.hdl.reviewStore.GetReviewResultByIssueId(ctx, req.GetReviewIssueId())
	suite.Assert().NoError(err)
	suite.Assert().Equal(true, result.GetResult())
	suite.Assert().Equal(comment, result.GetComment())
}

func (suite *CommitTreatmentReviewResultTestSuite) TearDownSuite() {}

func TestCommitTreatmentReviewResultTestSuite(t *testing.T) {
	suite.Run(t, new(CommitTreatmentReviewResultTestSuite))
}
