package report

import (
	"context"
	"testing"

	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GetCustomerConstitutionTestSuite struct {
	suite.Suite
	hdl *ReportAPIHandler
}

func (suite *GetCustomerConstitutionTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, nil, nil, "", "")
}

// TestGetCustomerConstitution 获取常客的体质报告.
func (suite *GetCustomerConstitutionTestSuite) TestGetCustomerConstitution() {
	ctx := context.Background()
	xid := xid.New()
	suite.T().Log(xid)
	req := &reportpb.GetCustomerConstitutionRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
	}
	resp := &reportpb.GetCustomerConstitutionResponse{}

	err := suite.hdl.GetCustomerConstitution(ctx, req, resp)
	suite.Assert().True(true, resp.GetIsReportAvailable)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *GetCustomerConstitutionTestSuite) TestGetReportNotAvailable() {
	ctx := context.Background()
	xid := xid.New()
	suite.T().Log(xid)
	req := &reportpb.GetCustomerConstitutionRequest{
		TenantId:   tenantId1,
		CustomerId: customerId,
	}
	resp := &reportpb.GetCustomerConstitutionResponse{}

	err := suite.hdl.GetCustomerConstitution(ctx, req, resp)
	suite.Assert().Equal(false, resp.IsReportAvailable)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *GetCustomerConstitutionTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetCustomerConstitutionRequest{
		TenantId:   tenantIdIsNull,
		CustomerId: customerId,
	}
	resp := &reportpb.GetCustomerConstitutionResponse{}
	err := suite.hdl.GetCustomerConstitution(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdIsNull
func (suite *GetCustomerConstitutionTestSuite) TestCustomerIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetCustomerConstitutionRequest{
		TenantId:   tenantId,
		CustomerId: customerIdIsNull,
	}
	resp := &reportpb.GetCustomerConstitutionResponse{}
	err := suite.hdl.GetCustomerConstitution(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdNotExist
func (suite *GetCustomerConstitutionTestSuite) TestCustomerIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetCustomerConstitutionRequest{
		TenantId:   tenantId,
		CustomerId: customerIdIsNotExist,
	}
	resp := &reportpb.GetCustomerConstitutionResponse{}
	err := suite.hdl.GetCustomerConstitution(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetCustomerConstitutionTestSuite) TearDownSuite() {
}

func TestGetCustomerConstitutionTestSuite(t *testing.T) {
	suite.Run(t, new(GetCustomerConstitutionTestSuite))
}
