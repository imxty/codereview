package notification

import (
	"context"
	"testing"

	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/notification"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CheckTxIdUsageTestSuite struct {
	suite.Suite
	hdl *NotificationAPIHandler
}

func (suite *CheckTxIdUsageTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, notificationFile)
	s := store.NewNotificationStore(conn)
	suite.hdl = NewNotificationAPIHandler(s, nil, "")
}

func (suite *CheckTxIdUsageTestSuite) TestCheckTxIdUsage() {
	ctx := context.Background()

	req := &notificationpb.CheckTxIdUsageRequest{
		TxId: txId,
	}
	resp := &notificationpb.CheckTxIdUsageResponse{}
	err := suite.hdl.CheckTxIdUsage(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *CheckTxIdUsageTestSuite) TestTxIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &notificationpb.CheckTxIdUsageRequest{
		TxId: txIdIsNull,
	}
	resp := &notificationpb.CheckTxIdUsageResponse{}
	err := suite.hdl.CheckTxIdUsage(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *CheckTxIdUsageTestSuite) TearDownSuite() {

}

func TestCheckTxIdUsageTestSuite(t *testing.T) {
	suite.Run(t, new(CheckTxIdUsageTestSuite))
}
