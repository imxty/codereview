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

type SubmitTreatmentToTenantTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *SubmitTreatmentToTenantTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *SubmitTreatmentToTenantTestSuite) TestSubmitTreatmentToTenant() {
	ctx := context.Background()
	req := &userpb.SubmitTreatmentToTenantRequest{
		OrganizationId: organizationId,
		TenantIds:      []string{tenantId},
		TreatmentId:    treatmentId,
	}
	resp := &userpb.SubmitTreatmentToTenantResponse{}
	err := suite.hdl.SubmitTreatmentToTenant(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestOrganizationIdIsNull
func (suite *SubmitTreatmentToTenantTestSuite) TestOrganizationIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SubmitTreatmentToTenantRequest{
		OrganizationId: organizationIdIsNull,
		TenantIds:      []string{tenantId},
		TreatmentId:    treatmentId,
	}
	resp := &userpb.SubmitTreatmentToTenantResponse{}
	err := suite.hdl.SubmitTreatmentToTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdsIsNull
func (suite *SubmitTreatmentToTenantTestSuite) TestTenantIdsIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SubmitTreatmentToTenantRequest{
		OrganizationId: organizationId,
		TenantIds:      nil,
		TreatmentId:    treatmentId,
	}
	resp := &userpb.SubmitTreatmentToTenantResponse{}
	err := suite.hdl.SubmitTreatmentToTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTreatmentIdIsNull
func (suite *SubmitTreatmentToTenantTestSuite) TestTreatmentIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SubmitTreatmentToTenantRequest{
		OrganizationId: organizationId,
		TenantIds:      []string{tenantId},
		TreatmentId:    treatmentIdIsNull,
	}
	resp := &userpb.SubmitTreatmentToTenantResponse{}
	err := suite.hdl.SubmitTreatmentToTenant(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SubmitTreatmentToTenantTestSuite) TearDownSuite() {

}

func TestSubmitTreatmentToTenantTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitTreatmentToTenantTestSuite))
}
