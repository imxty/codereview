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

type BatchGetCustomersByIdsTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *BatchGetCustomersByIdsTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

// TestCustomerIdIsNull
func (suite *BatchGetCustomersByIdsTestSuite) TestCustomerIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.BatchGetCustomersByIdsRequest{
		CustomerIds: nil,
	}
	resp := &customerpb.BatchGetCustomersByIdsResponse{}
	err := suite.hdl.BatchGetCustomersByIds(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdIsErr
func (suite *BatchGetCustomersByIdsTestSuite) TestCustomerIdIsErr() {

	ctx := context.Background()
	req := &customerpb.BatchGetCustomersByIdsRequest{
		CustomerIds: []string{customerIdNotExist},
	}
	resp := &customerpb.BatchGetCustomersByIdsResponse{}
	err := suite.hdl.BatchGetCustomersByIds(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestBatchGetCustomersByIds 查询常客
func (suite *BatchGetCustomersByIdsTestSuite) TestBatchGetCustomersByIds() {
	ctx := context.Background()
	req := &customerpb.BatchGetCustomersByIdsRequest{
		CustomerIds: []string{customerId},
	}
	resp := &customerpb.BatchGetCustomersByIdsResponse{}
	err := suite.hdl.BatchGetCustomersByIds(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *BatchGetCustomersByIdsTestSuite) TearDownSuite() {

}

func TestBatchGetCustomersByIdsTestSuite(t *testing.T) {
	suite.Run(t, new(BatchGetCustomersByIdsTestSuite))
}
