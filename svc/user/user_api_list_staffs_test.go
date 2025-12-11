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

type ListStaffsTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *ListStaffsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *ListStaffsTestSuite) TestListStaffs() {
	ctx := context.Background()
	req := &userpb.ListStaffsRequest{
		TenantId: tenantId,
	}
	resp := &userpb.ListStaffsResponse{}
	err := suite.hdl.ListStaffs(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(int32(2), resp.ActivatedStaffCount)
	suite.Assert().Equal(int32(6), resp.MaxStaffCount)
}

// TestTenantIdIsNull
func (suite *ListStaffsTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ListStaffsRequest{
		TenantId: tenantIdIsNull,
	}
	resp := &userpb.ListStaffsResponse{}
	err := suite.hdl.ListStaffs(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *ListStaffsTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.ListStaffsRequest{
		TenantId: tenantIdNotExist,
	}
	resp := &userpb.ListStaffsResponse{}
	err := suite.hdl.ListStaffs(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

func (suite *ListStaffsTestSuite) TearDownSuite() {
}

func TestListStaffsTestSuite(t *testing.T) {
	suite.Run(t, new(ListStaffsTestSuite))
}
