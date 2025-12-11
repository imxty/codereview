package customer

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/customer"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	usermock "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GetTenantsCustomerRankTestSuite struct {
	suite.Suite
	hdl  *CustomerAPIHandler
	user *usermock.UserAPIService
}

func (suite *GetTenantsCustomerRankTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	user := &usermock.UserAPIService{}
	suite.user = user
	suite.hdl = NewCustomerAPIHandler(s, user)
}

func (suite *GetTenantsCustomerRankTestSuite) TestGetTenantsCustomerRank() {
	ctx := context.Background()
	suite.user.On("BatchGetTenantNamesByIDs", mock.Anything, mock.Anything).Return(
		&userpb.BatchGetTenantNamesByIDsResponse{
			TenantNames: make(map[string]string),
		}, nil).After(utils.RpcLatency())
	req := &customerpb.GetTenantsCustomerRankRequest{
		TenantIds: []string{tenantId},
	}
	resp := &customerpb.GetTenantsCustomerRankResponse{}
	err := suite.hdl.GetTenantsCustomerRank(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)

}

// TestTenantIdIsNull
func (suite *GetTenantsCustomerRankTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.GetTenantsCustomerRankRequest{
		TenantIds: nil,
	}
	resp := &customerpb.GetTenantsCustomerRankResponse{}
	err := suite.hdl.GetTenantsCustomerRank(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetTenantsCustomerRankTestSuite) TearDownSuite() {

}

func TestGetTenantsCustomerRankTestSuite(t *testing.T) {
	suite.Run(t, new(GetTenantsCustomerRankTestSuite))
}
