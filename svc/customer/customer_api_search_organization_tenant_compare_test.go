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
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SearchOrganizationTenantCompareTestSuite struct {
	suite.Suite
	hdl  *CustomerAPIHandler
	user *usermock.UserAPIService
}

func (suite *SearchOrganizationTenantCompareTestSuite) SetupTest() {
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

// TestSearchOrganizationTenantCompare 通过nickname 查询
func (suite *SearchOrganizationTenantCompareTestSuite) TestSearchOrganizationTenantCompare() {
	ctx := context.Background()
	suite.user.On("SearchTenantByNameAndOrganizationID", mock.Anything, mock.Anything).Return(
		&userpb.SearchTenantByNameAndOrganizationIDResponse{
			TenantIds: []string{tenantId},
		}, nil).After(utils.RpcLatency())
	suite.user.On("BatchGetTenantNamesByIDs", mock.Anything, mock.Anything).Return(
		&userpb.BatchGetTenantNamesByIDsResponse{
			TenantNames: make(map[string]string),
		}, nil).After(utils.RpcLatency())
	req := &customerpb.SearchOrganizationTenantCompareRequest{
		OrganizationId: organizationId,
		TenantName:     nickname,
		StartTime:      &timestamppb.Timestamp{},
		EndTime:        &timestamppb.Timestamp{},
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchOrganizationTenantCompareResponse{}
	err := suite.hdl.SearchOrganizationTenantCompare(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SearchOrganizationTenantCompareTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.SearchOrganizationTenantCompareRequest{
		OrganizationId: organizationIdIsNull,
		TenantName:     nickname,
		StartTime:      &timestamppb.Timestamp{},
		EndTime:        &timestamppb.Timestamp{},
		Pagination: &customerpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &customerpb.SearchOrganizationTenantCompareResponse{}
	err := suite.hdl.SearchOrganizationTenantCompare(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)

}

func (suite *SearchOrganizationTenantCompareTestSuite) TearDownSuite() {

}

func TestSearchOrganizationTenantCompareTestSuite(t *testing.T) {
	suite.Run(t, new(SearchOrganizationTenantCompareTestSuite))
}
