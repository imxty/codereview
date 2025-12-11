package report

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AddCustomerReportTestSuite struct {
	suite.Suite
	hdl      *ReportAPIHandler
	customer *customermock.CustomerAPIService
}

func (suite *AddCustomerReportTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	customer := &customermock.CustomerAPIService{}
	suite.customer = customer
	suite.hdl = NewReportAPIHandler(s, nil, customer, nil, nil, nil, "", "")
}

// TestAddCustomerReport 添加常客报告请求
func (suite *AddCustomerReportTestSuite) TestAddCustomerReport() {
	ctx := context.Background()
	suite.customer.On("AddCustomer", mock.Anything, mock.Anything).Return(
		&customerpb.AddCustomerResponse{
			Customer: &customerpb.Customer{
				CustomerId: customerId,
			},
		}, nil).After(utils.RpcLatency())
	req := &reportpb.AddCustomerReportRequest{
		TenantId: tenantId,
		ReportId: reportId,
		Customer: &reportpb.Customer{
			CustomerId: customerId,
			StaffId:    staffId,
			Nickname:   nickname,
			Gender:     reportpb.Gender_GENDER_FEMALE,
			AreaCode:   areaCode,
			Phone:      phone,
			Birthday: &reportpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  height,
			Weight:  weight,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &reportpb.AddCustomerReportResponse{}
	err := suite.hdl.AddCustomerReport(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *AddCustomerReportTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.AddCustomerReportRequest{
		TenantId: tenantIdIsNull,
		ReportId: reportId,
		Customer: &reportpb.Customer{
			CustomerId: customerId,
			StaffId:    staffId,
			Nickname:   nickname,
			Gender:     reportpb.Gender_GENDER_FEMALE,
			AreaCode:   areaCode,
			Phone:      phone,
			Birthday: &reportpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  height,
			Weight:  weight,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &reportpb.AddCustomerReportResponse{}
	err := suite.hdl.AddCustomerReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *AddCustomerReportTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.AddCustomerReportRequest{
		TenantId: tenantIdIsNotExist,
		ReportId: reportId,
		Customer: &reportpb.Customer{
			CustomerId: customerId,
			StaffId:    staffId,
			Nickname:   nickname,
			Gender:     reportpb.Gender_GENDER_FEMALE,
			AreaCode:   areaCode,
			Phone:      phone,
			Birthday: &reportpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  height,
			Weight:  weight,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &reportpb.AddCustomerReportResponse{}
	err := suite.hdl.AddCustomerReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdIsNull
func (suite *AddCustomerReportTestSuite) TestCustomerIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.AddCustomerReportRequest{
		TenantId: tenantId,
		ReportId: reportId,
		Customer: &reportpb.Customer{
			CustomerId: customerId,
			StaffId:    staffId,
			Nickname:   nickname,
			Gender:     reportpb.Gender_GENDER_FEMALE,
			AreaCode:   areaCode,
			Phone:      phone,
			Birthday: &reportpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  height,
			Weight:  weight,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &reportpb.AddCustomerReportResponse{}
	err := suite.hdl.AddCustomerReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestCustomerIdNotExist
func (suite *AddCustomerReportTestSuite) TestCustomerIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &reportpb.AddCustomerReportRequest{
		TenantId: tenantId,
		ReportId: reportId,
		Customer: &reportpb.Customer{
			CustomerId: customerId,
			StaffId:    staffId,
			Nickname:   nickname,
			Gender:     reportpb.Gender_GENDER_FEMALE,
			AreaCode:   areaCode,
			Phone:      phone,
			Birthday: &reportpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  height,
			Weight:  weight,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &reportpb.AddCustomerReportResponse{}
	err := suite.hdl.AddCustomerReport(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *AddCustomerReportTestSuite) TearDownSuite() {
}

func TestAddCustomerReportTestSuite(t *testing.T) {
	suite.Run(t, new(AddCustomerReportTestSuite))
}
