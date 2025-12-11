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

type GetCustomerTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *GetCustomerTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

func (suite *GetCustomerTestSuite) TestGetCustomer() {
	ctx := context.Background()
	req := &customerpb.GetCustomerRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
	}
	resp := &customerpb.GetCustomerResponse{}
	err := suite.hdl.GetCustomer(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(customerId, resp.Customer.CustomerId)
	suite.Assert().Equal(nickname, resp.Customer.Nickname)
	suite.Assert().Equal(staffId2, resp.Customer.StaffId)
	suite.Assert().Equal(customerpb.Gender_GENDER_FEMALE, resp.Customer.Gender)
	suite.Assert().Equal(phone, resp.Customer.Phone)
}

// TestTenantIdIsNull
func (suite *GetCustomerTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.GetCustomerRequest{
		TenantId:   tenantIdIsNull,
		CustomerId: customerId,
	}
	resp := &customerpb.GetCustomerResponse{}
	err := suite.hdl.GetCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdIsErr
func (suite *GetCustomerTestSuite) TestTenantIdIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.GetCustomerRequest{
		TenantId:   tenantIdWithNoCustomers,
		CustomerId: customerId,
	}
	resp := &customerpb.GetCustomerResponse{}
	err := suite.hdl.GetCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrCustomerNotFound, err)
}

// TestCustomerIdIsNull
func (suite *GetCustomerTestSuite) TestCustomerIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.GetCustomerRequest{
		TenantId:   tenantId,
		CustomerId: customerIdIsNull,
	}
	resp := &customerpb.GetCustomerResponse{}
	err := suite.hdl.GetCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdNotExist
func (suite *GetCustomerTestSuite) TestCustomerIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.GetCustomerRequest{
		TenantId:   tenantId,
		CustomerId: customerIdNotExist,
	}
	resp := &customerpb.GetCustomerResponse{}
	err := suite.hdl.GetCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrCustomerNotFound, err)
}

func (suite *GetCustomerTestSuite) TearDownSuite() {

}

func TestGetCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(GetCustomerTestSuite))
}
