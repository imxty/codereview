package report

import (
	"context"
	"testing"

	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ModifyReportStaffRemarkTestSuite struct {
	suite.Suite
	hdl *ReportAPIHandler
}

func (suite *ModifyReportStaffRemarkTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, nil, nil, "", "")
}

// TestModifyReportStaffRemark  修改员工备注
func (suite *ModifyReportStaffRemarkTestSuite) TestModifyReportStaffRemark() {
	ctx := context.Background()
	req := &reportpb.ModifyReportRemarkRequest{
		TenantId: tenantId,
		ReportId: reportId,
		Remark:   remark,
	}
	resp := &reportpb.ModifyReportRemarkResponse{}
	err := suite.hdl.ModifyReportRemark(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *ModifyReportStaffRemarkTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.ModifyReportRemarkRequest{
		TenantId: tenantIdIsNull,
		ReportId: reportId,
		Remark:   remark,
	}
	resp := &reportpb.ModifyReportRemarkResponse{}
	err := suite.hdl.ModifyReportRemark(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReportIdIsNull
func (suite *ModifyReportStaffRemarkTestSuite) TestReportIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.ModifyReportRemarkRequest{
		TenantId: tenantId,
		ReportId: reportIdIsNull,
		Remark:   remark,
	}
	resp := &reportpb.ModifyReportRemarkResponse{}
	err := suite.hdl.ModifyReportRemark(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReportNotExist
func (suite *ModifyReportStaffRemarkTestSuite) TestReportNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.ModifyReportRemarkRequest{
		TenantId: tenantId,
		ReportId: reportIdIsNotExist,
		Remark:   remark,
	}
	resp := &reportpb.ModifyReportRemarkResponse{}
	err := suite.hdl.ModifyReportRemark(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrReportNotFound, err)
}
func (suite *ModifyReportStaffRemarkTestSuite) TearDownSuite() {
}

func TestModifyReportStaffRemarkTestSuite(t *testing.T) {
	suite.Run(t, new(ModifyReportStaffRemarkTestSuite))
}
