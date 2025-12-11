package review

import (
	"context"
	"testing"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	product "github.com/jinmukeji/huimaibao-service/svc/product/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	user "github.com/jinmukeji/huimaibao-service/svc/user/mock"

	"github.com/stretchr/testify/suite"
)

type ListReviewResultByTargetIdsTestSuite struct {
	suite.Suite
	userAPI    *user.UserAPIService
	productAPI *product.ProductAPIService
	hdl        *ReviewAPIHandler
}

func (suite *ListReviewResultByTargetIdsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reviewFile)
	s := store.NewReviewStore(conn)
	userAPI := &user.UserAPIService{}
	suite.userAPI = userAPI
	productAPI := &product.ProductAPIService{}
	suite.productAPI = productAPI
	suite.hdl = NewReviewAPIHandler(s, productAPI, userAPI, nil)
}

// TestListReviewResultByTargetIds 正常流程
func (suite *ListReviewResultByTargetIdsTestSuite) TestListReviewResultByTargetIds() {
	ctx := context.Background()
	req := &pb.ListReviewResultByTargetIdsRequest{
		TargetId: []string{targetId},
	}
	resp := &pb.ListReviewResultByTargetIdsResponse{}
	err := suite.hdl.ListReviewResultByTargetIds(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestTargetIdIsNull
func (suite *ListReviewResultByTargetIdsTestSuite) TestTargetIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.ListReviewResultByTargetIdsRequest{
		TargetId: []string{},
	}
	resp := &pb.ListReviewResultByTargetIdsResponse{}
	err := suite.hdl.ListReviewResultByTargetIds(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListReviewResultByTargetIdsTestSuite) TearDownSuite() {}

func TestListReviewResultByTargetIdsTestSuite(t *testing.T) {
	suite.Run(t, new(ListReviewResultByTargetIdsTestSuite))
}
