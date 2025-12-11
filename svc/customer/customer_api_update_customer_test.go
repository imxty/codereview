package customer

import (
	"context"
	"testing"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/customer"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UpdateCustomerTestSuite struct {
	suite.Suite
	hdl *CustomerAPIHandler
}

func (suite *UpdateCustomerTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, customerFile)
	s := store.NewCustomerStore(conn)
	suite.hdl = NewCustomerAPIHandler(s, nil)
}

// TestUpdateCustomer
// - nickname -- 常客昵称
// - phone -- 常客手机号
// - birthday -- 常客生日
// - height -- 常客身高
// - weight -- 常客体重
// - pmh -- 常客既往病史
// - remarks -- 常客其他备注
func (suite *UpdateCustomerTestSuite) TestUpdateCustomer() {
	ctx := context.Background()
	req := &customerpb.UpdateCustomerRequest{
		TenantId: tenantId,
		Customer: &customerpb.Customer{
			CustomerId: customerId,
			StaffId:    staffId,
			Nickname:   nickname,
			Phone:      phone,
			Birthday: &customerpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  height,
			Weight:  weight,
			Gender:  customerpb.Gender_GENDER_FEMALE,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &customerpb.UpdateCustomerResponse{}
	err := suite.hdl.UpdateCustomer(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *UpdateCustomerTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.UpdateCustomerRequest{
		TenantId: tenantIdIsNull,
		Customer: &customerpb.Customer{
			CustomerId: customerId,
			StaffId:    staffId,
			Nickname:   nickname,
			Phone:      phone,
			Birthday: &customerpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  height,
			Weight:  weight,
			Gender:  customerpb.Gender_GENDER_FEMALE,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &customerpb.UpdateCustomerResponse{}
	err := suite.hdl.UpdateCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestFieldMaskWithOnePath
func (suite *UpdateCustomerTestSuite) TestFieldMaskWithOnePath() {
	ctx := context.Background()
	req := &customerpb.UpdateCustomerRequest{
		TenantId: tenantId,
		Customer: &customerpb.Customer{
			CustomerId: customerId,
			StaffId:    staffId,
			Nickname:   nickname,
			Phone:      phoneNotExist,
			Birthday: &customerpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  height,
			Weight:  weight,
			Gender:  customerpb.Gender_GENDER_FEMALE,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &customerpb.UpdateCustomerResponse{}
	err := suite.hdl.UpdateCustomer(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}
func (suite *UpdateCustomerTestSuite) TearDownSuite() {

}

func TestUpdateCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateCustomerTestSuite))
}
