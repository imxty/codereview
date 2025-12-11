package product

import (
	"context"
	"testing"

	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/product"
	reviewmock "github.com/jinmukeji/huimaibao-service/svc/review/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	s3 "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type ListTreatmentNameByIdsTestSuite struct {
	suite.Suite
	hdl    *ProductAPIHandler
	s3     *s3.FileStore
	review *reviewmock.ReviewAPIService
}

func (suite *ListTreatmentNameByIdsTestSuite) SetupSuite() {
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

// TestListTreatmentNameByIds
func (suite *ListTreatmentNameByIdsTestSuite) TestListTreatmentNameByIds() {
	ctx := context.Background()
	req := &productpb.ListTreatmentNameByIdsRequest{
		TreatmentIds: []string{treatmentId},
	}
	resp := &productpb.ListTreatmentNameByIdsResponse{}
	err := suite.hdl.ListTreatmentNameByIds(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestTreatmentIdsIsNull
func (suite *ListTreatmentNameByIdsTestSuite) TestTreatmentIdsIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &productpb.ListTreatmentNameByIdsRequest{
		TreatmentIds: []string{},
	}
	resp := &productpb.ListTreatmentNameByIdsResponse{}
	err := suite.hdl.ListTreatmentNameByIds(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListTreatmentNameByIdsTestSuite) TearDownSuite() {}

func TestListTreatmentNameByIdsTestSuite(t *testing.T) {
	suite.Run(t, new(ListTreatmentNameByIdsTestSuite))
}
