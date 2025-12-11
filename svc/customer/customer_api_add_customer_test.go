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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AddCustomerTestSuite struct {
	suite.Suite
	hdl  *CustomerAPIHandler
	user *usermock.UserAPIService
}

func (suite *AddCustomerTestSuite) SetupTest() {
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

// TestAddCustomer 添加常客
func (suite *AddCustomerTestSuite) TestAddCustomer() {
	ctx := context.Background()
	suite.user.On("GetTenantEntity", mock.Anything, mock.Anything).Return(
		&userpb.GetTenantEntityResponse{
			Entity: &userpb.Entity{
				TenantId: tenantId,
			},
		}, nil).After(utils.RpcLatency())
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantId,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phoneIsNew,
			Birthday: &customerpb.Date{
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
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(nickname, resp.Customer.Nickname)
	suite.Assert().Equal(staffId, resp.Customer.StaffId)
	suite.Assert().Equal(customerpb.Gender_GENDER_FEMALE, resp.Customer.Gender)
	suite.Assert().Equal(phoneIsNew, resp.Customer.Phone)
}

// TestTenantIdIsNull
func (suite *AddCustomerTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantIdIsNull,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phone,
			Birthday: &customerpb.Date{
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
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdIsNull
func (suite *AddCustomerTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantId,
		Customer: &customerpb.Customer{
			StaffId:  staffIdIsNull,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phone,
			Birthday: &customerpb.Date{
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
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestNicknameIsNull
func (suite *AddCustomerTestSuite) TestNicknameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantIdIsNull,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nicknameIsNull,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phoneIsNew,
			Birthday: &customerpb.Date{
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
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsNull
func (suite *AddCustomerTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantIdIsNull,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phoneIsNull,
			Birthday: &customerpb.Date{
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
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsErr
func (suite *AddCustomerTestSuite) TestPhoneIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantIdIsNull,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phoneIsErr,
			Birthday: &customerpb.Date{
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
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsExist
func (suite *AddCustomerTestSuite) TestPhoneIsExist() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantIdIsNull,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phoneIsExist,
			Birthday: &customerpb.Date{
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
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestBirthdayIsNull
func (suite *AddCustomerTestSuite) TestBirthdayIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantIdIsNull,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phone,
			Birthday: &customerpb.Date{},
			Height:   height,
			Weight:   weight,
			Pmh:      pmh,
			Remarks:  remarks,
		},
	}
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestHeightIsNull
func (suite *AddCustomerTestSuite) TestHeightIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantId,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phone,
			Birthday: &customerpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  heightIsBeyondLimit,
			Weight:  weight,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestHeightIsErr
func (suite *AddCustomerTestSuite) TestHeightIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantId,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phone,
			Birthday: &customerpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  heightIsErr,
			Weight:  weight,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestWeightIsErr
func (suite *AddCustomerTestSuite) TestWeightIsErr() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantId,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phone,
			Birthday: &customerpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  height,
			Weight:  weightIsErr,
			Pmh:     pmh,
			Remarks: remarks,
		},
	}
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestRemarksErr
func (suite *AddCustomerTestSuite) TestRemarksErr() {
	t := suite.T()
	ctx := context.Background()
	req := &customerpb.AddCustomerRequest{
		TenantId: tenantId,
		Customer: &customerpb.Customer{
			StaffId:  staffId,
			Nickname: nickname,
			Gender:   customerpb.Gender_GENDER_FEMALE,
			Phone:    phone,
			Birthday: &customerpb.Date{
				Year:  year,
				Month: month,
				Day:   day,
			},
			Height:  height,
			Weight:  weight,
			Pmh:     pmh,
			Remarks: remarksExcel,
		},
	}
	resp := &customerpb.AddCustomerResponse{}
	err := suite.hdl.AddCustomer(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *AddCustomerTestSuite) TearDownSuite() {

}

func TestAddCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(AddCustomerTestSuite))
}
