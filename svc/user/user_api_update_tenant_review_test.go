package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/suite"
)

type UpdateTenantReviewTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *UpdateTenantReviewTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *UpdateTenantReviewTestSuite) TestUpdateTenantReview() {
	ctx := context.Background()

	req := &userpb.UpdateTenantReviewRequest{
		TenantId: tenantId,
		Days:     2,
	}
	resp := &userpb.UpdateTenantReviewResponse{}
	err := suite.hdl.UpdateTenantReview(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *UpdateTenantReviewTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateTenantReviewRequest{
		TenantId: tenantIdIsNull,
		Days:     2,
	}
	resp := &userpb.UpdateTenantReviewResponse{}
	err := suite.hdl.UpdateTenantReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *UpdateTenantReviewTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateTenantReviewRequest{
		TenantId: tenantIdNotExist,
		Days:     -1,
	}
	resp := &userpb.UpdateTenantReviewResponse{}
	err := suite.hdl.UpdateTenantReview(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *UpdateTenantReviewTestSuite) TearDownSuite() {

}

func TestUpdateTenantReviewTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTenantReviewTestSuite))
}
