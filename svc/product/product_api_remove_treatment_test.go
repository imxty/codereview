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

type RemoveTreatmentTestSuite struct {
	suite.Suite
	hdl    *ProductAPIHandler
	s3     *s3.FileStore
	review *reviewmock.ReviewAPIService
}

func (suite *RemoveTreatmentTestSuite) SetupSuite() {
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

// TestRemoveTreatment
func (suite *RemoveTreatmentTestSuite) TestRemoveTreatment() {
	ctx := context.Background()
	suite.review.On("ListReviewResultByTargetIds", mock.Anything, mock.Anything).Return(
		&reviewpb.ListReviewResultByTargetIdsResponse{
			Results: make(map[string]*reviewpb.ReviewResult)}, nil).After(utils.RpcLatency())
	req := &productpb.RemoveTreatmentRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentId,
	}
	resp := &productpb.RemoveTreatmentResponse{}
	err := suite.hdl.RemoveTreatment(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestOrganizationIdIsNull
func (suite *RemoveTreatmentTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.RemoveTreatmentRequest{
		OrganizationId: organizationIdIsNull,
		TreatmentId:    treatmentId,
	}
	resp := &productpb.RemoveTreatmentResponse{}
	err := suite.hdl.RemoveTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdIsNull
func (suite *RemoveTreatmentTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.RemoveTreatmentRequest{
		OrganizationId: organizationId1,
		TreatmentId:    treatmentIdIsNull,
	}
	resp := &productpb.RemoveTreatmentResponse{}
	err := suite.hdl.RemoveTreatment(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *RemoveTreatmentTestSuite) TearDownSuite() {}

func TestRemoveTreatmentTestSuite(t *testing.T) {
	suite.Run(t, new(RemoveTreatmentTestSuite))
}
