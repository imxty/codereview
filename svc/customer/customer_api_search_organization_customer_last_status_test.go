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

type SearchOrganizationCustomerLastStatusTestSuite struct {
	suite.Suite
	hdl  *CustomerAPIHandler
	user *usermock.UserAPIService
}

func (suite *SearchOrganizationCustomerLastStatusTestSuite) SetupSuite() {
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
func (suite *SearchOrganizationCustomerLastStatusTestSuite) SetupTest() {
	suite.user.On("SearchTenantsByName", mock.Anything, mock.Anything).Return(
		&userpb.SearchTenantsByNameResponse{
			TenantIds: []string{tenantId},
		}, nil).After(utils.RpcLatency())
	suite.user.On("BatchGetTenantNamesByIDs", mock.Anything, mock.Anything).Return(
		&userpb.BatchGetTenantNamesByIDsResponse{
			TenantNames: make(map[string]string),
		}, nil).After(utils.RpcLatency())
}

// TestSearchOrganizationCustomerLastStatus
func (suite *SearchOrganizationCustomerLastStatusTestSuite) TestSearchOrganizationCustomerLastStatus() {
	ctx := context.Background()
	req := &customerpb.SearchOrganizationCustomerLastStatusRequest{
		OrganizationId: organizationId,
		TenantName:     nickname,
		CustomerName:   nickname,
		CustomerPhone:  phone,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchOrganizationCustomerLastStatusResponse{}
	err := suite.hdl.SearchOrganizationCustomerLastStatus(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestSearchOrganizationId 根据组织ID查询
func (suite *SearchOrganizationCustomerLastStatusTestSuite) TestSearchOrganizationId() {
	ctx := context.Background()
	req := &customerpb.SearchOrganizationCustomerLastStatusRequest{
		OrganizationId: organizationId,
		TenantName:     nicknameIsNull,
		CustomerName:   nicknameIsNull,
		CustomerPhone:  phoneIsNull,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchOrganizationCustomerLastStatusResponse{}
	err := suite.hdl.SearchOrganizationCustomerLastStatus(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestSearchTenantName
func (suite *SearchOrganizationCustomerLastStatusTestSuite) TestSearchTenantName() {
	ctx := context.Background()
	t := suite.T()
	req := &customerpb.SearchOrganizationCustomerLastStatusRequest{
		OrganizationId: organizationIdIsNull,
		TenantName:     nickname,
		CustomerName:   nicknameIsNull,
		CustomerPhone:  phoneIsNull,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchOrganizationCustomerLastStatusResponse{}
	err := suite.hdl.SearchOrganizationCustomerLastStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)

}

// TestSearchCustomerName
func (suite *SearchOrganizationCustomerLastStatusTestSuite) TestSearchCustomerName() {
	ctx := context.Background()
	req := &customerpb.SearchOrganizationCustomerLastStatusRequest{
		OrganizationId: organizationId,
		TenantName:     nicknameIsNull,
		CustomerName:   nickname,
		CustomerPhone:  phoneIsNull,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchOrganizationCustomerLastStatusResponse{}
	err := suite.hdl.SearchOrganizationCustomerLastStatus(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestSearchCustomerPhone
func (suite *SearchOrganizationCustomerLastStatusTestSuite) TestSearchCustomerPhone() {
	ctx := context.Background()
	req := &customerpb.SearchOrganizationCustomerLastStatusRequest{
		OrganizationId: organizationId,
		TenantName:     nicknameIsNull,
		CustomerName:   nicknameIsNull,
		CustomerPhone:  phone,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchOrganizationCustomerLastStatusResponse{}
	err := suite.hdl.SearchOrganizationCustomerLastStatus(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *SearchOrganizationCustomerLastStatusTestSuite) TearDownSuite() {

}

func TestSearchOrganizationCustomerLastStatusTestSuite(t *testing.T) {
	suite.Run(t, new(SearchOrganizationCustomerLastStatusTestSuite))
}
