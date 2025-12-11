package product

import (
	"context"
	"testing"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reviewpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/product"
	reviewmock "github.com/jinmukeji/huimaibao-service/svc/review/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	s3 "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SubmitTreatmentToReviewTestSuite struct {
	suite.Suite
	hdl    *ProductAPIHandler
	s3     *s3.FileStore
	review *reviewmock.ReviewAPIService
}

func (suite *SubmitTreatmentToReviewTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, productFile)
	s := store.NewProductStore(conn)
	s3 := &s3.FileStore{}
	suite.s3 = s3
	review := &reviewmock.ReviewAPIService{}
	suite.review = review
	suite.hdl = NewProductAPIHandler(s, s3, review, nil)
}

// TestSubmitTreatmentToReview
func (suite *SubmitTreatmentToReviewTestSuite) TestSubmitTreatmentToReview() {
	ctx := context.Background()
	suite.review.On("SubmitTreatmentToReview", mock.Anything, mock.Anything).Return(
		&reviewpb.SubmitTreatmentToReviewResponse{}, nil).After(utils.RpcLatency())
	req := &productpb.SubmitTreatmentToReviewRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentIdIsDraft,
	}
	resp := &productpb.SubmitTreatmentToReviewResponse{}
	err := suite.hdl.SubmitTreatmentToReview(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestTreatmentIdIsPassed
func (suite *SubmitTreatmentToReviewTestSuite) TestTreatmentIdIsPassed() {
	t := suite.T()
	ctx := context.Background()

	req := &productpb.SubmitTreatmentToReviewRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentIdIsPassed,
	}
	resp := &productpb.SubmitTreatmentToReviewResponse{}
	err := suite.hdl.SubmitTreatmentToReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOrganizationIdIsNull
func (suite *SubmitTreatmentToReviewTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()

	req := &productpb.SubmitTreatmentToReviewRequest{
		OrganizationId: organizationIdIsNull,
		TreatmentId:    treatmentId,
	}
	resp := &productpb.SubmitTreatmentToReviewResponse{}
	err := suite.hdl.SubmitTreatmentToReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.DataAccessFailed, err)
}

// TestTreatmentIdIsNull
func (suite *SubmitTreatmentToReviewTestSuite) TestTreatmentIdIsNull() {
	t := suite.T()
	ctx := context.Background()

	req := &productpb.SubmitTreatmentToReviewRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentIdIsNull,
	}
	resp := &productpb.SubmitTreatmentToReviewResponse{}
	err := suite.hdl.SubmitTreatmentToReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.DataAccessFailed, err)
}

func (suite *SubmitTreatmentToReviewTestSuite) TearDownSuite() {}

func TestSubmitTreatmentToReviewTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitTreatmentToReviewTestSuite))
}
