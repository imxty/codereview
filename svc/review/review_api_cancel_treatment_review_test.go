package review

import (
	"context"
	"testing"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type CancelTreatmentReviewTestSuite struct {
	suite.Suite
	hdl *ReviewAPIHandler
}

func (suite *CancelTreatmentReviewTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reviewFile)
	s := store.NewReviewStore(conn)
	suite.hdl = NewReviewAPIHandler(s, nil, nil, nil)
}

// TestCancelTreatmentReview
func (suite *CancelTreatmentReviewTestSuite) TestCancelTreatmentReview() {
	ctx := context.Background()
	req := &pb.CancelTreatmentReviewRequest{
		OrganizationId: organizationId,
		TreatmentId:    treatmentId,
	}
	resp := &pb.CancelTreatmentReviewResponse{}
	err := suite.hdl.CancelTreatmentReview(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(false, resp.ReviewIsTerminated)
}

// TestTreatmentReviewIsTerminated
func (suite *CancelTreatmentReviewTestSuite) TestTreatmentReviewIsTerminated() {
	t := suite.T()
	ctx := context.Background()
	req := &pb.CancelTreatmentReviewRequest{
		OrganizationId: organizationId,
		TreatmentId:    tenantTreatmentId3,
	}
	resp := &pb.CancelTreatmentReviewResponse{}
	err := suite.hdl.CancelTreatmentReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestNoRemainingChances
func (suite *CancelTreatmentReviewTestSuite) TestNoRemainingChances() {
	ctx := context.Background()
	req := &pb.CancelTreatmentReviewRequest{
		OrganizationId: organizationId,
		TreatmentId:    tenantTreatmentId1,
	}
	resp := &pb.CancelTreatmentReviewResponse{}
	err := suite.hdl.CancelTreatmentReview(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().Equal(false, resp.ReviewIsTerminated)
}

// TestTenantTreatmentIdSuccess
func (suite *CancelTreatmentReviewTestSuite) TestTenantTreatmentIdSuccess() {
	t := suite.T()
	ctx := context.Background()
	req := &pb.CancelTreatmentReviewRequest{
		OrganizationId: organizationId,
		TreatmentId:    tenantTreatmentId4,
	}
	resp := &pb.CancelTreatmentReviewResponse{}
	err := suite.hdl.CancelTreatmentReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestOrganizationIdIsEmpty
func (suite *CancelTreatmentReviewTestSuite) TestOrganizationIdIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &pb.CancelTreatmentReviewRequest{
		OrganizationId: organizationIdIsNull,
		TreatmentId:    tenantTreatmentId4,
	}
	resp := &pb.CancelTreatmentReviewResponse{}
	err := suite.hdl.CancelTreatmentReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentIdIsEmpty
func (suite *CancelTreatmentReviewTestSuite) TestTreatmentIdIsEmpty() {
	t := suite.T()
	ctx := context.Background()
	req := &pb.CancelTreatmentReviewRequest{
		OrganizationId: organizationId,
		TreatmentId:    treatmentIdIsNull,
	}
	resp := &pb.CancelTreatmentReviewResponse{}
	err := suite.hdl.CancelTreatmentReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *CancelTreatmentReviewTestSuite) TearDownSuite() {}

func TestCancelTreatmentReviewTestSuite(t *testing.T) {
	suite.Run(t, new(CancelTreatmentReviewTestSuite))
}
