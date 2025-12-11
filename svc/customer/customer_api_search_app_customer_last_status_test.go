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

type SearchAppCustomerLastStatusTestSuite struct {
	suite.Suite
	hdl  *CustomerAPIHandler
	user *usermock.UserAPIService
}

func (suite *SearchAppCustomerLastStatusTestSuite) SetupTest() {
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

// TestSearchAppCustomerLastStatus
func (suite *SearchAppCustomerLastStatusTestSuite) TestSearchAppCustomerLastStatus() {
	ctx := context.Background()
	suite.user.On("BatchGetTenantNamesByIDs", mock.Anything, mock.Anything).Return(
		&userpb.BatchGetTenantNamesByIDsResponse{
			TenantNames: make(map[string]string),
		}, nil).After(utils.RpcLatency())
	req := &customerpb.SearchAppCustomerLastStatusRequest{
		TenantId: tenantId,
		Key:      key,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchAppCustomerLastStatusResponse{}
	err := suite.hdl.SearchAppCustomerLastStatus(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *SearchAppCustomerLastStatusTestSuite) TestTenantIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &customerpb.SearchAppCustomerLastStatusRequest{
		TenantId: tenantIdIsNull,
		Key:      key,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		}}
	resp := &customerpb.SearchAppCustomerLastStatusResponse{}
	err := suite.hdl.SearchAppCustomerLastStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SearchAppCustomerLastStatusTestSuite) TearDownSuite() {

}

func TestSearchAppCustomerLastStatusTestSuite(t *testing.T) {
	suite.Run(t, new(SearchAppCustomerLastStatusTestSuite))
}
