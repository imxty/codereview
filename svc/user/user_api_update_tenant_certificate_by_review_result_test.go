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

type UpdateTenantCertificateByReviewResultTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *UpdateTenantCertificateByReviewResultTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *UpdateTenantCertificateByReviewResultTestSuite) TestUpdateTenantCertificateByReviewResult() {
	ctx := context.Background()

	req := &userpb.UpdateTenantCertificateByReviewResultRequest{
		TenantId:   tenantId,
		Status:     true,
		FailReason: failReason,
	}
	resp := &userpb.UpdateTenantCertificateByReviewResultResponse{}
	err := suite.hdl.UpdateTenantCertificateByReviewResult(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestTenantIdIsNull
func (suite *UpdateTenantCertificateByReviewResultTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateTenantCertificateByReviewResultRequest{
		TenantId:   tenantIdIsNull,
		Status:     true,
		FailReason: failReason,
	}
	resp := &userpb.UpdateTenantCertificateByReviewResultResponse{}
	err := suite.hdl.UpdateTenantCertificateByReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestTenantIdNotExist
func (suite *UpdateTenantCertificateByReviewResultTestSuite) TestTenantIdNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &userpb.UpdateTenantCertificateByReviewResultRequest{
		TenantId:   tenantIdNotExist,
		Status:     true,
		FailReason: failReason,
	}
	resp := &userpb.UpdateTenantCertificateByReviewResultResponse{}
	err := suite.hdl.UpdateTenantCertificateByReviewResult(ctx, req, resp)
	utils.AssertRpcErrorCode(t, ErrTenantNotFound, err)
}

func (suite *UpdateTenantCertificateByReviewResultTestSuite) TearDownSuite() {

}

func TestUpdateTenantCertificateByReviewResultTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTenantCertificateByReviewResultTestSuite))
}
