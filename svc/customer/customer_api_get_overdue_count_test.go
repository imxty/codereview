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

type GetOverdueCountTestSuite struct {
	suite.Suite
	hdl  *CustomerAPIHandler
	user *usermock.UserAPIService
}

func (suite *GetOverdueCountTestSuite) SetupTest() {
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

func (suite *GetOverdueCountTestSuite) TestGetOverdueCount() {
	ctx := context.Background()
	suite.user.On("GetTenantEntity", mock.Anything, mock.Anything).Return(
		&userpb.GetTenantEntityResponse{
			Entity: &userpb.Entity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())
	req := &customerpb.GetOverdueCountRequest{
		TenantId: tenantId,
	}
	resp := &customerpb.GetOverdueCountResponse{}
	err := suite.hdl.GetOverdueCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(0), resp.TotalCount)

}

// TestTenantIdIsNull
func (suite *GetOverdueCountTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.GetOverdueCountRequest{
		TenantId: tenantIdIsNull,
	}
	resp := &customerpb.GetOverdueCountResponse{}
	err := suite.hdl.GetOverdueCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOverdueCountTestSuite) TearDownSuite() {

}

func TestGetOverdueCountTestSuite(t *testing.T) {
	suite.Run(t, new(GetOverdueCountTestSuite))
}
