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

type GetTreatmentByTreatmentRevIdTestSuite struct {
	suite.Suite
	hdl    *ProductAPIHandler
	s3     *s3.FileStore
	review *reviewmock.ReviewAPIService
}

func (suite *GetTreatmentByTreatmentRevIdTestSuite) SetupSuite() {
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

// TestGetTreatment
func (suite *GetTreatmentByTreatmentRevIdTestSuite) TestGetTreatmentByTreatmentRevIdTestSuite() {
	ctx := context.Background()
	suite.review.On("GetReviewResultByTargetId", mock.Anything, mock.Anything).Return(
		&reviewpb.GetReviewResultByTargetIdResponse{
			Result: &reviewpb.ReviewResult{
				ReviewType:  reviewpb.ReviewType_REVIEW_TYPE_CERTIFICATE,
				ReiewStatus: reviewpb.ReviewStatus_REVIEW_STATUS_PENDING_REVIEW,
				FailReason:  "",
			},
		}, nil).After(utils.RpcLatency())
	req := &productpb.GetTreatmentRevByTreatmentRevIdRequest{
		OrganizationId: organizationId1,
		TreatmentRevId: treatmentRevId,
	}
	resp := &productpb.GetTreatmentRevByTreatmentRevIdResponse{}
	err := suite.hdl.GetTreatmentRevByTreatmentRevId(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp.Treatment)
	suite.Assert().Equal(treatmentName, resp.GetTreatment().GetTreatmentName())
}

// TestOrganizationIdIsNull
func (suite *GetTreatmentByTreatmentRevIdTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.GetTreatmentRevByTreatmentRevIdRequest{
		OrganizationId: organizationIdIsNull,
		TreatmentRevId: treatmentRevId,
	}
	resp := &productpb.GetTreatmentRevByTreatmentRevIdResponse{}
	err := suite.hdl.GetTreatmentRevByTreatmentRevId(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentRevIdIsEmpty
func (suite *GetTreatmentByTreatmentRevIdTestSuite) TestTreatmentRevIdIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.GetTreatmentRevByTreatmentRevIdRequest{
		OrganizationId: organizationId1,
		TreatmentRevId: treatmentIdIsNull,
	}
	resp := &productpb.GetTreatmentRevByTreatmentRevIdResponse{}
	err := suite.hdl.GetTreatmentRevByTreatmentRevId(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetTreatmentByTreatmentRevIdTestSuite) TearDownSuite() {}

func TestGetTreatmentByTreatmentRevIdTestSuite(t *testing.T) {
	suite.Run(t, new(GetTreatmentByTreatmentRevIdTestSuite))
}
