package review

import (
	"context"
	"testing"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	productmock "github.com/jinmukeji/huimaibao-service/svc/product/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type ConfirmReviewNotificationTestSuite struct {
	suite.Suite
	hdl     *ReviewAPIHandler
	product *productmock.ProductAPIService
}

func (suite *ConfirmReviewNotificationTestSuite) SetupSuite() {
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

// TestConfirmReviewNotification
func (suite *ConfirmReviewNotificationTestSuite) TestConfirmReviewNotification() {
	ctx := context.Background()

	req := &pb.ConfirmReviewNotificationRequest{
		NotificationId: notificationId,
	}
	resp := &pb.ConfirmReviewNotificationResponse{}
	err := suite.hdl.ConfirmReviewNotification(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestNotificationIdIsNull
func (suite *ConfirmReviewNotificationTestSuite) TestNotificationIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.ConfirmReviewNotificationRequest{
		NotificationId: notificationIdIsNull,
	}
	resp := &pb.ConfirmReviewNotificationResponse{}
	err := suite.hdl.ConfirmReviewNotification(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *ConfirmReviewNotificationTestSuite) TearDownSuite() {}

func TestConfirmReviewNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(ConfirmReviewNotificationTestSuite))
}
