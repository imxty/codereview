package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/suite"
)

type GetStaffNamesByIDTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetStaffNamesByIDTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetStaffNamesByIDTestSuite) TestGetStaffNamesByID() {
	ctx := context.Background()

	req := &userpb.GetStaffNamesByIDRequest{
		StaffIds: []string{staffId},
	}
	resp := &userpb.GetStaffNamesByIDResponse{}
	err := suite.hdl.GetStaffNamesByID(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestStaffIdNotExist
func (suite *GetStaffNamesByIDTestSuite) TestStaffIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetStaffNamesByIDRequest{
		StaffIds: nil,
	}
	resp := &userpb.GetStaffNamesByIDResponse{}
	err := suite.hdl.GetStaffNamesByID(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetStaffNamesByIDTestSuite) TearDownSuite() {

}

func TestGetStaffNamesByIDTestSuite(t *testing.T) {
	suite.Run(t, new(GetStaffNamesByIDTestSuite))
}
