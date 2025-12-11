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

type ModifyReportRemarkTestSuite struct {
	suite.Suite
	hdl *ReportAPIHandler
}

func (suite *ModifyReportRemarkTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, nil, nil, "", "")
}

// TestModifyReportRemark  修改报告备注
func (suite *ModifyReportRemarkTestSuite) TestModifyReportRemark() {
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
func (suite *ModifyReportRemarkTestSuite) TestTenantIdIsNull() {
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
func (suite *ModifyReportRemarkTestSuite) TestReportIdIsNull() {
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
func (suite *ModifyReportRemarkTestSuite) TestReportNotExist() {
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

func (suite *ModifyReportRemarkTestSuite) TearDownSuite() {
}

func TestModifyReportRemarkTestSuite(t *testing.T) {
	suite.Run(t, new(ModifyReportRemarkTestSuite))
}
