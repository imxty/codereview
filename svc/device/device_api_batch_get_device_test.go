package device

import (
	"context"
	"testing"

	devicepb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/device"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BatchGetDeviceTestSuite struct {
	suite.Suite
	hdl *DeviceAPIHandler
}

func (suite *BatchGetDeviceTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, deviceFile)
	s := store.NewDeviceStore(conn)
	suite.hdl = NewDeviceAPIHandler(s)
}

func (suite *BatchGetDeviceTestSuite) TestBatchGetDevice() {
	ctx := context.Background()

	req := &devicepb.ListDevicesRequest{
		Mac: []string{""},
	}
	resp := &devicepb.ListDevicesResponse{}
	err := suite.hdl.ListDevices(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

// TestMacIsNull
func (suite *BatchGetDeviceTestSuite) TestMacIsNull() {
	t := suite.T()
	ctx := context.Background()

	req := &devicepb.ListDevicesRequest{
		Mac: nil,
	}
	resp := &devicepb.ListDevicesResponse{}
	err := suite.hdl.ListDevices(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *BatchGetDeviceTestSuite) TearDownSuite() {

}

func TestBatchGetDeviceTestSuite(t *testing.T) {
	suite.Run(t, new(BatchGetDeviceTestSuite))
}
