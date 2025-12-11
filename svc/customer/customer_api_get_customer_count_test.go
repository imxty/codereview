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

type GetCustomerCountTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *GetCustomerCountTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

func (suite *GetCustomerCountTestSuite) TestGetCustomerCount() {
	ctx := context.Background()
	req := &customerpb.GetCustomerCountRequest{
		TenantId: tenantId,
		StaffIds: []string{staffId},
	}
	resp := &customerpb.GetCustomerCountResponse{}
	err := suite.hdl.GetCustomerCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(0), resp.TodayAddedCustomerCount)
	suite.Assert().Equal(int32(5), resp.CustomerTotalCount)
}

// TestGetCustomerCountMore 多个员工
func (suite *GetCustomerCountTestSuite) TestGetCustomerCountMore() {
	ctx := context.Background()
	req := &customerpb.GetCustomerCountRequest{
		TenantId: tenantId,
		StaffIds: []string{staffId, staffId1},
	}
	resp := &customerpb.GetCustomerCountResponse{}
	err := suite.hdl.GetCustomerCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(0), resp.TodayAddedCustomerCount)
	suite.Assert().Equal(int32(5), resp.CustomerTotalCount)
}

// TestTenantIdIsNull
func (suite *GetCustomerCountTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.GetCustomerCountRequest{
		TenantId: tenantIdIsNull,
		StaffIds: []string{staffId},
	}
	resp := &customerpb.GetCustomerCountResponse{}
	err := suite.hdl.GetCustomerCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdIsNull
func (suite *GetCustomerCountTestSuite) TestStaffIdIsNull() {
	ctx := context.Background()
	req := &customerpb.GetCustomerCountRequest{
		TenantId: tenantId,
		StaffIds: []string{},
	}
	resp := &customerpb.GetCustomerCountResponse{}
	err := suite.hdl.GetCustomerCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *GetCustomerCountTestSuite) TearDownSuite() {

}

func TestGetCustomerCountTestSuite(t *testing.T) {
	suite.Run(t, new(GetCustomerCountTestSuite))
}
