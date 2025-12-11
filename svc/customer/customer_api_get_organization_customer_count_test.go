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

type GetOrganizationCustomerCountTestSuite struct {
	suite.Suite
	hdl  *CustomerAPIHandler
	user *usermock.UserAPIService
}

func (suite *GetOrganizationCustomerCountTestSuite) SetupTest() {
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

func (suite *GetOrganizationCustomerCountTestSuite) TestGetOrganizationCustomerCount() {
	ctx := context.Background()
	suite.user.On("GetOrganizationTenants", mock.Anything, mock.Anything).Return(
		&userpb.GetOrganizationTenantsResponse{
			Tenants: []*userpb.TenantEntity{
				&userpb.TenantEntity{
					OrganizationId: organizationId,
				},
			},
		}, nil).After(utils.RpcLatency())
	req := &customerpb.GetOrganizationCustomerCountRequest{
		OrganizationId: organizationId,
	}
	resp := &customerpb.GetOrganizationCustomerCountResponse{}
	err := suite.hdl.GetOrganizationCustomerCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(0), resp.TodayAddedCustomerCount)
	suite.Assert().Equal(int32(0), resp.TodayCustomerMeasurementCount)
	suite.Assert().Equal(int32(0), resp.TodayTempCustomerMeasurementCount)
	suite.Assert().Equal(int32(0), resp.CustomerTotalCount)
}

// TestOrganizationIdIsNull
func (suite *GetOrganizationCustomerCountTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.GetOrganizationCustomerCountRequest{
		OrganizationId: organizationIdIsNull,
	}
	resp := &customerpb.GetOrganizationCustomerCountResponse{}
	err := suite.hdl.GetOrganizationCustomerCount(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationCustomerCountTestSuite) TearDownSuite() {

}

func TestGetOrganizationCustomerCountTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationCustomerCountTestSuite))
}
