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

type SubmitFeedbackTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *SubmitFeedbackTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *SubmitFeedbackTestSuite) TestSubmitFeedback() {
	ctx := context.Background()
	req := &userpb.SubmitFeedbackRequest{
		TenantId: tenantId,
		Phone:    phone,
		Content:  content,
	}
	resp := &userpb.SubmitFeedbackResponse{}
	err := suite.hdl.SubmitFeedback(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *SubmitFeedbackTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SubmitFeedbackRequest{
		TenantId: tenantIdIsNull,
		Phone:    phone,
		Content:  content,
	}
	resp := &userpb.SubmitFeedbackResponse{}
	err := suite.hdl.SubmitFeedback(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPhoneIsNull
func (suite *SubmitFeedbackTestSuite) TestPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.SubmitFeedbackRequest{
		TenantId: tenantId,
		Phone:    phoneIsNull,
		Content:  content,
	}
	resp := &userpb.SubmitFeedbackResponse{}
	err := suite.hdl.SubmitFeedback(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *SubmitFeedbackTestSuite) TearDownSuite() {

}

func TestSubmitFeedbackTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitFeedbackTestSuite))
}
