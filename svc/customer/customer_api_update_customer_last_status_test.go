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

type UpdateCustomerLastStatusTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *UpdateCustomerLastStatusTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

// TestUpdateCustomerLastStatus

func (suite *UpdateCustomerLastStatusTestSuite) TestUpdateCustomerLastStatus() {
	ctx := context.Background()
	req := &customerpb.UpdateCustomerLastStatusRequest{
		CustomerId: customerId,
		ReportId:   reportId,
	}
	resp := &customerpb.UpdateCustomerLastStatusResponse{}
	err := suite.hdl.UpdateCustomerLastStatus(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestCustomerIdIsNull
func (suite *UpdateCustomerLastStatusTestSuite) TestCustomerIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.UpdateCustomerLastStatusRequest{
		CustomerId: customerIdIsNull,
		ReportId:   reportId,
	}
	resp := &customerpb.UpdateCustomerLastStatusResponse{}
	err := suite.hdl.UpdateCustomerLastStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReportIdIsNull
func (suite *UpdateCustomerLastStatusTestSuite) TestReportIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &customerpb.UpdateCustomerLastStatusRequest{
		CustomerId: customerId,
		ReportId:   reportIdIsNull,
	}
	resp := &customerpb.UpdateCustomerLastStatusResponse{}
	err := suite.hdl.UpdateCustomerLastStatus(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *UpdateCustomerLastStatusTestSuite) TearDownSuite() {

}

func TestUpdateCustomerLastStatusTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateCustomerLastStatusTestSuite))
}
