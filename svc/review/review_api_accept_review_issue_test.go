package review

import (
	"context"
	"testing"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	productmock "github.com/jinmukeji/huimaibao-service/svc/product/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	usermock "github.com/jinmukeji/huimaibao-service/svc/user/mock"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AcceptReviewIssueTestSuite struct {
	suite.Suite
	user    *usermock.UserAPIService
	product *productmock.ProductAPIService
	hdl     *ReviewAPIHandler
}

func (suite *AcceptReviewIssueTestSuite) SetupSuite() {
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

// TestAcceptCertificateReviewIssue 正常流程
func (suite *AcceptReviewIssueTestSuite) TestAcceptCertificateReviewIssue() {
	ctx := context.Background()
	suite.user.On("GetEntityRevision", mock.Anything, mock.Anything).Return(
		&userpb.GetEntityRevisionResponse{
			Entity: &userpb.TenantEntity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())

	suite.product.On("GetTreatmentRevByTreatmentRevId", mock.Anything, mock.Anything).Return(
		&productpb.GetTreatmentRevByTreatmentRevIdResponse{
			Treatment: &productpb.Treatment{
				OrganizationId: organizationId,
				TreatmentId:    tenantTreatmentId1,
			},
		}, nil).After(utils.RpcLatency())
	req := &pb.AcceptReviewIssueRequest{
		// 工单ID
		ReviewIssueId: reviewIssueIdCertificateToBeReviewed,
		// 审查人ID
		ReviewerUserId: reviewerUserId,
	}
	resp := &pb.AcceptReviewIssueResponse{}
	err := suite.hdl.AcceptReviewIssue(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(reviewIssueIdCertificateToBeReviewed, resp.ReviewIssueId)
	suite.Assert().Equal(tenantId, resp.Entity.TenantId)
	suite.Assert().Equal(tenantId, resp.SubmitterTenantId)
	suite.Assert().Equal(organizationId, resp.SubmitterOrganizationId)
}

// TestAcceptFailedCertificateReviewIssue
func (suite *AcceptReviewIssueTestSuite) TestAcceptFailedCertificateReviewIssue() {
	t := suite.T()
	ctx := context.Background()
	suite.user.On("GetEntityRevision", mock.Anything, mock.Anything).Return(
		&userpb.GetEntityRevisionResponse{
			Entity: &userpb.TenantEntity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())
	req := &pb.AcceptReviewIssueRequest{
		// 工单ID
		ReviewIssueId: reviewIssueIdCertificateFailed,
		// 审查人ID
		ReviewerUserId: reviewerUserId,
	}
	resp := &pb.AcceptReviewIssueResponse{}
	err := suite.hdl.AcceptReviewIssue(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidOperation, err)
}

// TestAcceptProductReviewIssue 正常流程
func (suite *AcceptReviewIssueTestSuite) TestAcceptProductReviewIssue() {
	ctx := context.Background()
	suite.product.On("GetTreatment", mock.Anything, mock.Anything, mock.Anything).
		Return(&productpb.GetTreatmentResponse{
			Treatment: &productpb.Treatment{
				TreatmentItemsRiskyDisease: []*productpb.TreatmentItem{{
					Symptom: symptomKey,
					Products: []*productpb.Product{
						{
							ProductId:     productId1,
							ProductType:   productpb.ProductType_PRODUCT_TYPE_CPD,
							ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_TO_BE_USED,
						},
						{
							ProductId:     productId2,
							ProductType:   productpb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT,
							ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_USING,
						},
					},
				}},
			},
		}, nil).After(utils.RpcLatency())
	req := &pb.AcceptReviewIssueRequest{
		// 工单ID
		ReviewIssueId: reviewIssueIdProductToBeReviewed,
		// 审查人ID
		ReviewerUserId: reviewerUserId,
	}
	resp := &pb.AcceptReviewIssueResponse{}
	err := suite.hdl.AcceptReviewIssue(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestAcceptFailedProductReviewIssue
func (suite *AcceptReviewIssueTestSuite) TestAcceptFailedProductReviewIssue() {
	ctx := context.Background()
	t := suite.T()
	suite.product.On("GetTreatment", mock.Anything, mock.Anything, mock.Anything).
		Return(&productpb.GetTreatmentResponse{
			Treatment: &productpb.Treatment{
				TreatmentItemsRiskyDisease: []*productpb.TreatmentItem{{
					Symptom: symptomKey,
					Products: []*productpb.Product{
						{
							ProductId:     productId1,
							ProductType:   productpb.ProductType_PRODUCT_TYPE_CPD,
							ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_TO_BE_USED,
						},
						{
							ProductId:     productId2,
							ProductType:   productpb.ProductType_PRODUCT_TYPE_HEALTH_CARE_PRODUCT,
							ProductStatus: productpb.ProductStatus_PRODUCT_STATUS_USING,
						},
					},
				}},
			},
		}, nil).After(utils.RpcLatency())
	req := &pb.AcceptReviewIssueRequest{
		// 工单ID
		ReviewIssueId: reviewIssueIdProductFailed,
		// 审查人ID
		ReviewerUserId: reviewerUserId,
	}
	resp := &pb.AcceptReviewIssueResponse{}
	err := suite.hdl.AcceptReviewIssue(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidOperation, err)
}

// TestReviewIssueIdIsNull
func (suite *AcceptReviewIssueTestSuite) TestReviewIssueIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.AcceptReviewIssueRequest{
		// 工单ID
		ReviewIssueId: reviewIssueIdIsNull,
		// 审查人ID
		ReviewerUserId: reviewerUserId,
	}
	resp := &pb.AcceptReviewIssueResponse{}
	err := suite.hdl.AcceptReviewIssue(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReviewUserIdIsNull
func (suite *AcceptReviewIssueTestSuite) TestReviewUserIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.AcceptReviewIssueRequest{
		// 工单ID
		ReviewIssueId: reviewIssueIdCertificateFailed,
		// 审查人ID
		ReviewerUserId: reviewerUserIdIsNull,
	}
	resp := &pb.AcceptReviewIssueResponse{}
	err := suite.hdl.AcceptReviewIssue(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestAcceptSuccessCertificateReviewIssue 已审核通过
func (suite *AcceptReviewIssueTestSuite) TestAcceptSuccessCertificateReviewIssue() {
	ctx := context.Background()
	t := suite.T()
	suite.user.On("GetEntityRevision", mock.Anything, mock.Anything).Return(
		&userpb.GetEntityRevisionResponse{
			Entity: &userpb.TenantEntity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())
	req := &pb.AcceptReviewIssueRequest{
		// 工单ID
		ReviewIssueId: reviewIssueIdCertificateSuccess,
		// 审查人ID
		ReviewerUserId: reviewerUserId,
	}
	resp := &pb.AcceptReviewIssueResponse{}
	err := suite.hdl.AcceptReviewIssue(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidOperation, err)
}

// TestAcceptReviewCertificateReviewIssue
func (suite *AcceptReviewIssueTestSuite) TestAcceptReviewCertificateReviewIssue() {
	ctx := context.Background()
	t := suite.T()
	suite.user.On("GetEntityRevision", mock.Anything, mock.Anything).Return(
		&userpb.GetEntityRevisionResponse{
			Entity: &userpb.TenantEntity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())
	req := &pb.AcceptReviewIssueRequest{
		// 工单ID
		ReviewIssueId: reviewIssueIdCertificateReview,
		// 审查人ID
		ReviewerUserId: reviewerUserId,
	}
	resp := &pb.AcceptReviewIssueResponse{}
	err := suite.hdl.AcceptReviewIssue(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidOperation, err)
}

func (suite *AcceptReviewIssueTestSuite) TearDownSuite() {}

func TestAcceptReviewIssueTestSuite(t *testing.T) {
	suite.Run(t, new(AcceptReviewIssueTestSuite))
}
