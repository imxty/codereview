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

type BatchGetStaffsTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *BatchGetStaffsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *BatchGetStaffsTestSuite) TestBatchGetStaffs() {
	ctx := context.Background()
	req := &userpb.BatchGetStaffsRequest{
		StaffIds: []string{staffId},
	}
	resp := &userpb.BatchGetStaffsResponse{}
	err := suite.hdl.BatchGetStaffs(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestStaffIdsIsNull
func (suite *BatchGetStaffsTestSuite) TestStaffIdsIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.BatchGetStaffsRequest{
		StaffIds: nil,
	}
	resp := &userpb.BatchGetStaffsResponse{}
	err := suite.hdl.BatchGetStaffs(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}
func (suite *BatchGetStaffsTestSuite) TearDownSuite() {

}

func TestBatchGetStaffsTestSuite(t *testing.T) {
	suite.Run(t, new(BatchGetStaffsTestSuite))
}
