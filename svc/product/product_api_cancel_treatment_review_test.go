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

type CancelTreatmentReviewTestSuite struct {
	suite.Suite
	hdl    *ProductAPIHandler
	s3     *s3.FileStore
	review *reviewmock.ReviewAPIService
}

func (suite *CancelTreatmentReviewTestSuite) SetupSuite() {
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

// TestCancelTreatmentReview
func (suite *CancelTreatmentReviewTestSuite) TestCancelTreatmentReview() {
	ctx := context.Background()
	suite.review.On("CancelTreatmentReview", mock.Anything, mock.Anything).Return(
		&reviewpb.CancelTreatmentReviewResponse{
			ReviewIsTerminated: true,
		}, nil).After(utils.RpcLatency())
	req := &productpb.CancelTreatmentReviewRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentId,
	}
	resp := &productpb.CancelTreatmentReviewResponse{}
	err := suite.hdl.CancelTreatmentReview(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(true, resp.ReviewIsTerminated)
}

// TestOrganizationIdIsNull
func (suite *CancelTreatmentReviewTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.CancelTreatmentReviewRequest{
		OrganizationId: organizationIdIsNull,
		TreatmentId:    treatmentId,
	}
	resp := &productpb.CancelTreatmentReviewResponse{}
	err := suite.hdl.CancelTreatmentReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentIdIsNull
func (suite *CancelTreatmentReviewTestSuite) TestTreatmentIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.CancelTreatmentReviewRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentIdIsNull,
	}
	resp := &productpb.CancelTreatmentReviewResponse{}
	err := suite.hdl.CancelTreatmentReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *CancelTreatmentReviewTestSuite) TearDownSuite() {}

func TestCancelTreatmentReviewTestSuite(t *testing.T) {
	suite.Run(t, new(CancelTreatmentReviewTestSuite))
}
