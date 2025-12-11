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

type SearchCustomerLastStatusTestSuite struct {
	suite.Suite
	hdl  *CustomerAPIHandler
	user *usermock.UserAPIService
}

func (suite *SearchCustomerLastStatusTestSuite) SetupTest() {
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

// TestSearchCustomerLastStatus 通过nickname 查询
func (suite *SearchCustomerLastStatusTestSuite) TestSearchCustomerLastStatus() {
	ctx := context.Background()
	suite.user.On("SearchTenantByNameAndOrganizationID", mock.Anything, mock.Anything).Return(
		&userpb.SearchTenantByNameAndOrganizationIDResponse{
			TenantIds: []string{tenantId},
		}, nil).After(utils.RpcLatency())
	suite.user.On("BatchGetTenantNamesByIDs", mock.Anything, mock.Anything).Return(
		&userpb.BatchGetTenantNamesByIDsResponse{
			TenantNames: make(map[string]string),
		}, nil).After(utils.RpcLatency())
	req := &customerpb.SearchCustomerLastStatusRequest{
		OrganizationId: organizationId,
		TenantName:     nickname,
		CustomerName:   nickname,
		CustomerPhone:  phone,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchCustomerLastStatusResponse{}
	err := suite.hdl.SearchCustomerLastStatus(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SearchCustomerLastStatusTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.SearchCustomerLastStatusRequest{
		OrganizationId: organizationIdIsNull,
		TenantName:     nickname,
		CustomerName:   nickname,
		CustomerPhone:  phone,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchCustomerLastStatusResponse{}
	err := suite.hdl.SearchCustomerLastStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)

}

func (suite *SearchCustomerLastStatusTestSuite) TearDownSuite() {

}

func TestSearchCustomerLastStatusTestSuite(t *testing.T) {
	suite.Run(t, new(SearchCustomerLastStatusTestSuite))
}
