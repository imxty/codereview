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

type GetSharedReportLinkTestSuite struct {
	suite.Suite
	hdl *ReportAPIHandler
}

func (suite *GetSharedReportLinkTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, nil, nil, "", "")
}

// TestGetSharedReportLink  获取分享报告连接
func (suite *GetSharedReportLinkTestSuite) TestGetSharedReportLink() {
	ctx := context.Background()
	req := &reportpb.GetSharedReportLinkRequest{
		TenantId: tenantId,
		ReportId: reportId,
	}
	resp := &reportpb.GetSharedReportLinkResponse{}
	err := suite.hdl.GetSharedReportLink(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *GetSharedReportLinkTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetSharedReportLinkRequest{
		TenantId: tenantIdIsNull,
		ReportId: reportId,
	}
	resp := &reportpb.GetSharedReportLinkResponse{}
	err := suite.hdl.GetSharedReportLink(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReportIdIsNull
func (suite *GetSharedReportLinkTestSuite) TestReportIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetSharedReportLinkRequest{
		TenantId: tenantId,
		ReportId: reportIdIsNull,
	}
	resp := &reportpb.GetSharedReportLinkResponse{}
	err := suite.hdl.GetSharedReportLink(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestReportIdNotExist
func (suite *GetSharedReportLinkTestSuite) TestReportIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.GetSharedReportLinkRequest{
		TenantId: tenantId,
		ReportId: reportIdIsNotExist,
	}
	resp := &reportpb.GetSharedReportLinkResponse{}
	err := suite.hdl.GetSharedReportLink(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrReportNotFound, err)
}

func (suite *GetSharedReportLinkTestSuite) TearDownSuite() {
}

func TestGetSharedReportLinkTestSuite(t *testing.T) {
	suite.Run(t, new(GetSharedReportLinkTestSuite))
}
