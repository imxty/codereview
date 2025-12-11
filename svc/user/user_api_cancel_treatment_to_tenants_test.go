package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/suite"
)

type CancelTreatmentToTenantsTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *CancelTreatmentToTenantsTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *CancelTreatmentToTenantsTestSuite) TestCancelTreatmentToTenants() {
	ctx := context.Background()
	req := &userpb.CancelTreatmentToTenantsRequest{
		OrganizationId: organizationId,
		TenantIds:      []string{tenantId},
		TreatmentId:    treatmentId,
	}
	resp := &userpb.CancelTreatmentToTenantsResponse{}
	err := suite.hdl.CancelTreatmentToTenants(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *CancelTreatmentToTenantsTestSuite) TestOrganizationIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &userpb.CancelTreatmentToTenantsRequest{
		OrganizationId: organizationIdIsNull,
		TenantIds:      []string{tenantId},
		TreatmentId:    treatmentId,
	}
	resp := &userpb.CancelTreatmentToTenantsResponse{}
	err := suite.hdl.CancelTreatmentToTenants(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdIsNull
func (suite *CancelTreatmentToTenantsTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.CancelTreatmentToTenantsRequest{
		OrganizationId: organizationId,
		TenantIds:      nil,
		TreatmentId:    treatmentId,
	}
	resp := &userpb.CancelTreatmentToTenantsResponse{}
	err := suite.hdl.CancelTreatmentToTenants(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentIdIsNull
func (suite *CancelTreatmentToTenantsTestSuite) TestTreatmentIdIsNull() {
	ctx := context.Background()
	req := &userpb.CancelTreatmentToTenantsRequest{
		OrganizationId: organizationId,
		TenantIds:      []string{tenantId},
		TreatmentId:    treatmentIdIsNull,
	}
	resp := &userpb.CancelTreatmentToTenantsResponse{}
	err := suite.hdl.CancelTreatmentToTenants(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *CancelTreatmentToTenantsTestSuite) TearDownSuite() {

}

func TestCancelTreatmentToTenantsTestSuite(t *testing.T) {
	suite.Run(t, new(CancelTreatmentToTenantsTestSuite))
}
