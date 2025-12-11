package customer

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/customer"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SearchCustomerIDByCustomerInfoTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *SearchCustomerIDByCustomerInfoTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

// TestSearchCustomerIDByCustomerInfo
func (suite *SearchCustomerIDByCustomerInfoTestSuite) TestSearchCustomerIDByCustomerInfo() {
	ctx := context.Background()
	req := &customerpb.SearchCustomerIDByCustomerInfoRequest{
		Name:  nickname,
		Phone: phone,
		Type:  customerpb.CustomerType_CUSTOMER_TYPE_CUSTOMER,
	}
	resp := &customerpb.SearchCustomerIDByCustomerInfoResponse{}
	err := suite.hdl.SearchCustomerIDByCustomerInfo(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestSearchPhone
func (suite *SearchCustomerIDByCustomerInfoTestSuite) TestSearchPhone() {
	ctx := context.Background()
	req := &customerpb.SearchCustomerIDByCustomerInfoRequest{
		Name:  nickname,
		Phone: phone,
		Type:  customerpb.CustomerType_CUSTOMER_TYPE_CUSTOMER,
	}
	resp := &customerpb.SearchCustomerIDByCustomerInfoResponse{}
	err := suite.hdl.SearchCustomerIDByCustomerInfo(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestSearchCustomerType
func (suite *SearchCustomerIDByCustomerInfoTestSuite) TestSearchCustomerType() {
	ctx := context.Background()
	req := &customerpb.SearchCustomerIDByCustomerInfoRequest{
		Name:  nickname,
		Phone: phone,
		Type:  customerpb.CustomerType_CUSTOMER_TYPE_CUSTOMER,
	}
	resp := &customerpb.SearchCustomerIDByCustomerInfoResponse{}
	err := suite.hdl.SearchCustomerIDByCustomerInfo(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestSearchCustomerName
func (suite *SearchCustomerIDByCustomerInfoTestSuite) TestSearchCustomerName() {
	ctx := context.Background()
	req := &customerpb.SearchCustomerIDByCustomerInfoRequest{
		Name:  nickname,
		Phone: phone,
		Type:  customerpb.CustomerType_CUSTOMER_TYPE_CUSTOMER,
	}
	resp := &customerpb.SearchCustomerIDByCustomerInfoResponse{}
	err := suite.hdl.SearchCustomerIDByCustomerInfo(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *SearchCustomerIDByCustomerInfoTestSuite) TearDownSuite() {

}

func TestSearchCustomerIDByCustomerInfoTestSuite(t *testing.T) {
	suite.Run(t, new(SearchCustomerIDByCustomerInfoTestSuite))
}
