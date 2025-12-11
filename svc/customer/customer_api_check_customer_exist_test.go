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

type CheckCustomerExistTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *CheckCustomerExistTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

// TestTenantIdIsNUll
func (suite *CheckCustomerExistTestSuite) TestTenantIdIsNUll() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.CheckCustomerExistRequest{
		TenantId:   tenantIdIsNull,
		CustomerId: customerId,
	}
	resp := &customerpb.CheckCustomerExistResponse{}
	err := suite.hdl.CheckCustomerExist(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdIsNull
func (suite *CheckCustomerExistTestSuite) TestCustomerIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.CheckCustomerExistRequest{
		TenantId:   tenantId,
		CustomerId: customerIdIsNull,
	}
	resp := &customerpb.CheckCustomerExistResponse{}
	err := suite.hdl.CheckCustomerExist(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdIsErr
func (suite *CheckCustomerExistTestSuite) TestCustomerIdIsErr() {

	ctx := context.Background()
	req := &customerpb.CheckCustomerExistRequest{
		TenantId:   tenantId,
		CustomerId: customerIdNotExist,
	}
	resp := &customerpb.CheckCustomerExistResponse{}
	err := suite.hdl.CheckCustomerExist(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(false, resp.CustomerExisting)
}

// TestCheckCustomerExist 查询常客
func (suite *CheckCustomerExistTestSuite) TestCheckCustomerExist() {
	ctx := context.Background()
	req := &customerpb.CheckCustomerExistRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
	}
	resp := &customerpb.CheckCustomerExistResponse{}
	err := suite.hdl.CheckCustomerExist(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(true, resp.CustomerExisting)
}

func (suite *CheckCustomerExistTestSuite) TearDownSuite() {

}

func TestCheckCustomerExistTestSuite(t *testing.T) {
	suite.Run(t, new(CheckCustomerExistTestSuite))
}
