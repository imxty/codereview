package user

import (
	"context"
	"testing"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type SubmitSubscriptionTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *SubmitSubscriptionTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *SubmitSubscriptionTestSuite) TestSubmitSubscription() {
	ctx := context.Background()
	req := &userpb.SubmitSubscriptionRequest{
		TenantId:  []string{tenantId},
		TimeCount: 3,
		UserId:    userId,
	}
	resp := &userpb.SubmitSubscriptionResponse{}
	err := suite.hdl.SubmitSubscription(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdsIsNull
func (suite *SubmitSubscriptionTestSuite) TestTenantIdsIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SubmitSubscriptionRequest{
		TenantId:  nil,
		TimeCount: 3,
		UserId:    userId,
	}
	resp := &userpb.SubmitSubscriptionResponse{}
	err := suite.hdl.SubmitSubscription(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestUserIdIsNull
func (suite *SubmitSubscriptionTestSuite) TestUserIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SubmitSubscriptionRequest{
		TenantId:  []string{tenantId},
		TimeCount: 3,
		UserId:    userIdIsNull,
	}
	resp := &userpb.SubmitSubscriptionResponse{}
	err := suite.hdl.SubmitSubscription(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SubmitSubscriptionTestSuite) TearDownSuite() {

}

func TestSubmitSubscriptionTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitSubscriptionTestSuite))
}
