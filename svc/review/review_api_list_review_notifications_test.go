package review

import (
	"context"
	"testing"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	product "github.com/jinmukeji/huimaibao-service/svc/product/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	user "github.com/jinmukeji/huimaibao-service/svc/user/mock"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ListReviewNotificationsTestSuite struct {
	suite.Suite
	userAPI    *user.UserAPIService
	productAPI *product.ProductAPIService
	hdl        *ReviewAPIHandler
}

func (suite *ListReviewNotificationsTestSuite) SetupSuite() {
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

// TestListReviewNotificationss 正常流程
func (suite *ListReviewNotificationsTestSuite) TestListReviewNotifications() {
	ctx := context.Background()
	suite.userAPI.On("ListSystemUsersByIds", mock.Anything, mock.Anything).Return(
		&userpb.ListSystemUsersByIdsResponse{
			SystemUsers: make(map[string]*userpb.SystemUser),
		}, nil).After(utils.RpcLatency())
	req := &pb.ListReviewNotificationsRequest{
		OrganizationId: organizationId,
		TenantId:       tenantId,
	}
	resp := &pb.ListReviewNotificationsResponse{}
	err := suite.hdl.ListReviewNotifications(ctx, req, resp)
	suite.Assert().NoError(err)
}

// TestOrganizationIdIsNull
func (suite *ListReviewNotificationsTestSuite) TestOrganizationIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &pb.ListReviewNotificationsRequest{
		OrganizationId: organizationIdIsNull,
		TenantId:       tenantIdIsNull,
	}
	resp := &pb.ListReviewNotificationsResponse{}
	err := suite.hdl.ListReviewNotifications(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListReviewNotificationsTestSuite) TearDownSuite() {}

func TestListReviewNotificationsTestSuite(t *testing.T) {
	suite.Run(t, new(ListReviewNotificationsTestSuite))
}
