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

type BatchGetOrganizationNamesByIdsTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *BatchGetOrganizationNamesByIdsTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *BatchGetOrganizationNamesByIdsTestSuite) TestBatchGetOrganizationNamesbyIds() {
	ctx := context.Background()
	req := &userpb.BatchGetOrganizationNamesByIDsRequest{
		OrganizationId: []string{organizationId},
	}
	resp := &userpb.BatchGetOrganizationNamesByIDsResponse{}
	err := suite.hdl.BatchGetOrganizationNamesByIDs(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *BatchGetOrganizationNamesByIdsTestSuite) TearDownSuite() {

}

func TestBatchGetOrganizationNamesByIdsTestSuite(t *testing.T) {
	suite.Run(t, new(BatchGetOrganizationNamesByIdsTestSuite))
}
