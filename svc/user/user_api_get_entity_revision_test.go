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

type GetEntityRevisionTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetEntityRevisionTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetEntityRevisionTestSuite) TestGetEntityRevision() {
	ctx := context.Background()
	req := &userpb.GetEntityRevisionRequest{
		RevisionId: revisionId,
	}
	resp := &userpb.GetEntityRevisionResponse{}
	err := suite.hdl.GetEntityRevision(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
	suite.Assert().Equal(tenantId, resp.Entity.TenantId)
	suite.Assert().Equal(phone, resp.Entity.ContactPhone)
}

// TestTenantIdIsNull
func (suite *GetEntityRevisionTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.GetEntityRevisionRequest{
		RevisionId: revisionIdIsNull,
	}
	resp := &userpb.GetEntityRevisionResponse{}
	err := suite.hdl.GetEntityRevision(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetEntityRevisionTestSuite) TearDownSuite() {

}

func TestGetEntityRevisionTestSuite(t *testing.T) {
	suite.Run(t, new(GetEntityRevisionTestSuite))
}
