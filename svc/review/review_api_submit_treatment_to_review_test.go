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

type SubmitTreatmentToReviewTestSuite struct {
	suite.Suite
	hdl     *ReviewAPIHandler
	product *productmock.ProductAPIService
}

func (suite *SubmitTreatmentToReviewTestSuite) SetupSuite() {
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

// TestSubmitTreatmentToReview
func (suite *SubmitTreatmentToReviewTestSuite) TestSubmitTreatmentToReview() {
	ctx := context.Background()
	suite.product.On("UpdateTreatmentStatus", mock.Anything, mock.Anything).Return(
		&productpb.UpdateTreatmentStatusResponse{}, nil).After(utils.RpcLatency())
	req := &pb.SubmitTreatmentToReviewRequest{
		OrganizationId: organizationId,
		TreatmentRevId: tenantTreatmentId1,
	}
	resp := &pb.SubmitTreatmentToReviewResponse{}
	err := suite.hdl.SubmitTreatmentToReview(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestOrganizationIdIsNull
func (suite *SubmitTreatmentToReviewTestSuite) TestOrganizationIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.SubmitTreatmentToReviewRequest{
		OrganizationId: organizationIdIsNull,
		TreatmentRevId: tenantTreatmentId1,
	}
	resp := &pb.SubmitTreatmentToReviewResponse{}
	err := suite.hdl.SubmitTreatmentToReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentRevIdIsNull
func (suite *SubmitTreatmentToReviewTestSuite) TestTreatmentRevIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.SubmitTreatmentToReviewRequest{
		OrganizationId: organizationId,
		TreatmentRevId: treatmentIdIsNull,
	}
	resp := &pb.SubmitTreatmentToReviewResponse{}
	err := suite.hdl.SubmitTreatmentToReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SubmitTreatmentToReviewTestSuite) TearDownSuite() {}

func TestSubmitTreatmentToReviewTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitTreatmentToReviewTestSuite))
}
