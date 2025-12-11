package customer

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/customer"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SearchCustomersTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *SearchCustomersTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

// TestSearchCustomers 通过nickname 查询
func (suite *SearchCustomersTestSuite) TestSearchCustomers() {
	ctx := context.Background()
	req := &customerpb.SearchCustomersRequest{
		TenantId: tenantId,
		Keywords: nickname,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchCustomersResponse{}
	err := suite.hdl.SearchCustomers(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestSearchCustomersByPhone 通过phone查询
func (suite *SearchCustomersTestSuite) TestSearchCustomersByPhone() {
	ctx := context.Background()
	req := &customerpb.SearchCustomersRequest{
		TenantId: tenantId,
		Keywords: phone,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchCustomersResponse{}
	err := suite.hdl.SearchCustomers(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *SearchCustomersTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.SearchCustomersRequest{
		TenantId: tenantIdIsNull,
		Keywords: phone,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchCustomersResponse{}
	err := suite.hdl.SearchCustomers(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)

}

// TestKeywordsIsNull
func (suite *SearchCustomersTestSuite) TestKeywordsIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.SearchCustomersRequest{
		TenantId:   tenantId,
		Keywords:   "",
		Pagination: nil,
	}
	resp := &customerpb.SearchCustomersResponse{}
	err := suite.hdl.SearchCustomers(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)

}

func (suite *SearchCustomersTestSuite) TearDownSuite() {

}

func TestSearchCustomersTestSuite(t *testing.T) {
	suite.Run(t, new(SearchCustomersTestSuite))
}
