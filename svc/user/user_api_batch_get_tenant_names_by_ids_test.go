package user

import (
	"context"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/suite"
)

type BatchGetTenantNamesByIdsTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *BatchGetTenantNamesByIdsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *BatchGetTenantNamesByIdsTestSuite) TestBatchGetTenantNamesByIds() {
	ctx := context.Background()
	req := &userpb.BatchGetTenantNamesByIDsRequest{
		TenantIds: []string{tenantId},
	}
	resp := &userpb.BatchGetTenantNamesByIDsResponse{}
	err := suite.hdl.BatchGetTenantNamesByIDs(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *BatchGetTenantNamesByIdsTestSuite) TearDownSuite() {

}

func TestBatchGetTenantNamesByIdsTestSuite(t *testing.T) {
	suite.Run(t, new(BatchGetTenantNamesByIdsTestSuite))
}
