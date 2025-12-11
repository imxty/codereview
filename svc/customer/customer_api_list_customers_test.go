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

type ListCustomersTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *ListCustomersTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

// TestListCustomers
func (suite *ListCustomersTestSuite) TestListCustomers() {
	ctx := context.Background()
	req := &customerpb.ListCustomersRequest{
		TenantId: tenantId,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
		StaffId: staffId,
	}
	resp := &customerpb.ListCustomersResponse{}
	err := suite.hdl.ListCustomers(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *ListCustomersTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.ListCustomersRequest{
		TenantId: tenantIdIsNull,
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
		StaffId: staffId,
	}
	resp := &customerpb.ListCustomersResponse{}
	err := suite.hdl.ListCustomers(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdIsNull
func (suite *ListCustomersTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.ListCustomersRequest{
		TenantId: tenantIdNotExist,
		Pagination: &customerpb.Pagination{
			Offset: offsetIsNull,
			Size:   sizeIsNull,
		},
		StaffId: staffIdIsNull,
	}
	resp := &customerpb.ListCustomersResponse{}
	err := suite.hdl.ListCustomers(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestHavingNoCustomers
func (suite *ListCustomersTestSuite) TestHavingNoCustomers() {
	ctx := context.Background()
	req := &customerpb.ListCustomersRequest{
		TenantId: tenantIdWithNoCustomers,
		Pagination: &customerpb.Pagination{
			Offset: offset1,
			Size:   size1,
		},
		StaffId: staffId,
	}
	resp := &customerpb.ListCustomersResponse{}
	err := suite.hdl.ListCustomers(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(0, len(resp.GetCustomers()))
}

func (suite *ListCustomersTestSuite) TearDownSuite() {

}

func TestListCustomersTestSuite(t *testing.T) {
	suite.Run(t, new(ListCustomersTestSuite))
}
