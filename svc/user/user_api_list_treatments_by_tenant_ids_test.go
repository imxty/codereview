package user

import (
	"context"
	"testing"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/suite"
)

type ListTreatmentsByTenantIDsTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *ListTreatmentsByTenantIDsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *ListTreatmentsByTenantIDsTestSuite) TestListTreatmentsByTenantIDs() {
	ctx := context.Background()
	req := &userpb.ListTreatmentsByTenantIDsRequest{
		TenantIds: []string{tenantId},
	}
	resp := &userpb.ListTreatmentsByTenantIDsResponse{}
	err := suite.hdl.ListTreatmentsByTenantIDs(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdsIsNull
func (suite *ListTreatmentsByTenantIDsTestSuite) TestTenantIdsIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ListTreatmentsByTenantIDsRequest{
		TenantIds: nil,
	}
	resp := &userpb.ListTreatmentsByTenantIDsResponse{}
	err := suite.hdl.ListTreatmentsByTenantIDs(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *ListTreatmentsByTenantIDsTestSuite) TearDownSuite() {
}

func TestListTreatmentsByTenantIDsTestSuite(t *testing.T) {
	suite.Run(t, new(ListTreatmentsByTenantIDsTestSuite))
}
