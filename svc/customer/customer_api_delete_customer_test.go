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

type DeleteCustomerTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *DeleteCustomerTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

func (suite *DeleteCustomerTestSuite) TestDeleteCustomer() {
	ctx := context.Background()
	req := &customerpb.DeleteCustomerRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
	}
	resp := &customerpb.DeleteCustomerResponse{}
	err := suite.hdl.DeleteCustomer(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *DeleteCustomerTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.DeleteCustomerRequest{
		TenantId:   tenantIdIsNull,
		CustomerId: customerId,
	}
	resp := &customerpb.DeleteCustomerResponse{}
	err := suite.hdl.DeleteCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdIsNull
func (suite *DeleteCustomerTestSuite) TestCustomerIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.DeleteCustomerRequest{
		TenantId:   tenantId,
		CustomerId: customerIdIsNull,
	}
	resp := &customerpb.DeleteCustomerResponse{}
	err := suite.hdl.DeleteCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdNotExist
func (suite *DeleteCustomerTestSuite) TestCustomerIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.DeleteCustomerRequest{
		TenantId:   tenantId,
		CustomerId: customerIdNotExist,
	}
	resp := &customerpb.DeleteCustomerResponse{}
	err := suite.hdl.DeleteCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrCustomerNotFound, err)
}

func (suite *DeleteCustomerTestSuite) TearDownSuite() {

}

func TestDeleteCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteCustomerTestSuite))
}
