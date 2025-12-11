package report

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jinmukeji/plat-pkg/v4/micro/meta"
	"github.com/rs/xid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/golang/protobuf/ptypes/timestamp"
	md5 "github.com/jinmukeji/go-pkg/v2/crypto/hash"
	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	devicepb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	customermock "github.com/jinmukeji/huimaibao-service/svc/customer/mock"
	devicemock "github.com/jinmukeji/huimaibao-service/svc/device/mock"
	calcmock "github.com/jinmukeji/huimaibao-service/svc/report/mock"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	usermock "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	calcpb "github.com/jinmukeji/platform-proto/gen/go/platform/report/v1"
	pt "github.com/jinmukeji/ptcodec"
	ptf "github.com/jinmukeji/ptcodec/ptfile"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SubmitPulseTestTestSuite struct {
	suite.Suite
	hdl      *ReportAPIHandler
	customer *customermock.CustomerAPIService
	user     *usermock.UserAPIService
	device   *devicemock.DeviceAPIService
	calc     *calcmock.ReportAPIClient
}

func (suite *SubmitPulseTestTestSuite) SetupTest() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	customer := &customermock.CustomerAPIService{}
	suite.customer = customer
	device := &devicemock.DeviceAPIService{}
	suite.device = device
	calc := &calcmock.ReportAPIClient{}
	suite.calc = calc
	user := &usermock.UserAPIService{}
	suite.user = user
	suite.hdl = NewReportAPIHandler(s, calc, customer, device, user, nil, "", "")
}

func (suite *SubmitPulseTestTestSuite) SubmitGeneralPulseTest(tenantId string, customerId string, age string, gender reportpb.Gender, staffId string, hand reportpb.Hand, sampleRate uint32, deviceMac string, deviceSn string, deviceModel string, codec string, data []byte, signature string, startTime *timestamp.Timestamp, stopTime *timestamp.Timestamp, enableStatistics bool) (*reportpb.SubmitPulseTestResponse, error) {
	ctx := context.Background()
	suite.customer.On("GetCustomer", mock.Anything, mock.Anything).Return(
		&customerpb.GetCustomerResponse{
			Customer: &customerpb.Customer{
				CustomerId: customerId,
			},
		}, nil).After(utils.RpcLatency())
	suite.device.On("ListDevices", mock.Anything, mock.Anything).Return(
		&devicepb.ListDevicesResponse{
			Devices: []*devicepb.Device{
				&devicepb.Device{
					DeviceId: deviceId,
				},
			},
		}, nil).After(utils.RpcLatency())
	suite.calc.On("SubmitPulseTest", mock.Anything, mock.Anything).Return(
		&calcpb.SubmitPulseTestResponse{
			ReportId: reportIdIsNew,
		}, nil).After(utils.RpcLatency())
	suite.user.On("GetTenantEntity", mock.Anything, mock.Anything, mock.Anything).
		Return(&userpb.GetTenantEntityResponse{
			Entity: &userpb.Entity{
				TenantId:     tenantId,
				TenantStatus: userpb.TenantStatus_TENANT_STATUS_USING,
			},
		}, nil).After(utils.RpcLatency())
	req := &reportpb.SubmitPulseTestRequest{
		TenantId:   tenantId,
		CustomerId: customerId,
		TempCustomer: &reportpb.TempCustomer{
			Age:    age,
			Gender: gender,
		},
		StaffId: staffId,
		Payload: &reportpb.SamplePayload{
			Hand: hand,
			SampleDevice: &reportpb.SampleDevice{
				SampleRate:   sampleRate,
				DeviceModel:  deviceModel,
				DeviceMac:    deviceMac,
				DeviceSn:     deviceSn,
				DeviceParams: make(map[string]string),
			},
			SampleData: &reportpb.SampleData{
				Codec:       codec,
				CodecParams: make(map[string]string),
				Data:        data,
				Signature:   signature,
			},
			SamplingStartTime: startTime,
			SamplingStopTime:  stopTime,

			GeoLocation: &reportpb.GeoLocation{
				Latitude:         100,
				Longitude:        100,
				Altitude:         nil,
				Accuracy:         0,
				AltitudeAccuracy: nil,
				Heading:          nil,
				Speed:            nil,
			},
		},
	}
	ctx = meta.ContextWithCid(ctx, xid.New().String())
	resp := &reportpb.SubmitPulseTestResponse{}
	err := suite.hdl.SubmitPulseTest(ctx, req, resp)
	return resp, err
}

// TestCustomerSubmitPulseTest 常客 （TempCustomer不传）
func (suite *SubmitPulseTestTestSuite) TestCustomerSubmitPulseTest() {
	gender := reportpb.Gender_GENDER_INVALID
	hand := reportpb.Hand_HAND_RIGHT
	data := ReadPayload(RawData)
	seconds := len(data) / fps
	now := time.Now()
	startTime := timestamppb.New(now)
	stopTime := timestamppb.New(now.Add(time.Second * time.Duration(seconds)))
	signature := string(md5.MD5(data))
	enableStatistics := true

	reply, err := suite.SubmitGeneralPulseTest(tenantId, customerId, ageIsNull, gender, staffId, hand, fps, deviceMac, deviceSn, deviceModel, codec, data, signature, startTime, stopTime, enableStatistics)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(reply)
}

// TestTempCustomerSubmitPulseTest 散客（custormerId不传）
func (suite *SubmitPulseTestTestSuite) TestTempCustomerSubmitPulseTest() {
	gender := reportpb.Gender_GENDER_INVALID
	hand := reportpb.Hand_HAND_RIGHT
	data := ReadPayload(RawData)
	seconds := len(data) / fps
	now := time.Now()
	startTime := timestamppb.New(now)
	stopTime := timestamppb.New(now.Add(time.Second * time.Duration(seconds)))
	signature := string(md5.MD5(data))
	enableStatistics := true

	reply, err := suite.SubmitGeneralPulseTest(tenantId, customerId, ageIsNull, gender, staffId, hand, fps, deviceMac, deviceSn, deviceModel, codec, data, signature, startTime, stopTime, enableStatistics)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(reply)
}

// TestTenantIdIsNull
func (suite *SubmitPulseTestTestSuite) TestTenantIdIsNull() {
	t := suite.T()
	gender := reportpb.Gender_GENDER_INVALID
	hand := reportpb.Hand_HAND_RIGHT
	data := ReadPayload(RawData)
	seconds := len(data) / fps
	now := time.Now()
	startTime := timestamppb.New(now)
	stopTime := timestamppb.New(now.Add(time.Second * time.Duration(seconds)))
	signature := string(md5.MD5(data))
	enableStatistics := true

	reply, err := suite.SubmitGeneralPulseTest(tenantIdIsNull, customerId, ageIsNull, gender, staffId, hand, fps, deviceMac, deviceSn, deviceModel, codec, data, signature, startTime, stopTime, enableStatistics)
	suite.Assert().NotNil(reply)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestStaffIdIsNull
func (suite *SubmitPulseTestTestSuite) TestStaffIdIsNull() {
	t := suite.T()
	gender := reportpb.Gender_GENDER_FEMALE
	hand := reportpb.Hand_HAND_RIGHT
	data := ReadPayload(RawData)
	seconds := len(data) / fps
	now := time.Now()
	startTime := timestamppb.New(now)
	stopTime := timestamppb.New(now.Add(time.Second * time.Duration(seconds)))
	signature := string(md5.MD5(data))
	enableStatistics := true

	reply, err := suite.SubmitGeneralPulseTest(tenantId, customerIdIsNull, age, gender, staffIdIsNull, hand, fps, deviceMac, deviceSn, deviceModel, codec, data, signature, startTime, stopTime, enableStatistics)
	suite.Assert().NotNil(reply)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestHandIsInvalid  无效Hand
func (suite *SubmitPulseTestTestSuite) TestHandIsInvalid() {
	t := suite.T()
	gender := reportpb.Gender_GENDER_INVALID
	hand := reportpb.Hand_HAND_INVALID
	data := ReadPayload(RawData)
	seconds := len(data) / fps
	now := time.Now()
	startTime := timestamppb.New(now)
	stopTime := timestamppb.New(now.Add(time.Second * time.Duration(seconds)))
	signature := string(md5.MD5(data))
	enableStatistics := true

	reply, err := suite.SubmitGeneralPulseTest(tenantId, customerId, ageIsNull, gender, staffId, hand, fps, deviceMac, deviceSn, deviceModel, codec, data, signature, startTime, stopTime, enableStatistics)
	suite.Assert().NotNil(reply)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestHandIsUnset  未设置手指
func (suite *SubmitPulseTestTestSuite) TestHandIsUnset() {
	t := suite.T()
	gender := reportpb.Gender_GENDER_INVALID
	hand := reportpb.Hand_HAND_UNSET
	data := ReadPayload(RawData)
	seconds := len(data) / fps
	now := time.Now()
	startTime := timestamppb.New(now)
	stopTime := timestamppb.New(now.Add(time.Second * time.Duration(seconds)))
	signature := string(md5.MD5(data))
	enableStatistics := true

	reply, err := suite.SubmitGeneralPulseTest(tenantId, customerId, ageIsNull, gender, staffId, hand, fps, deviceMac, deviceSn, deviceModel, codec, data, signature, startTime, stopTime, enableStatistics)
	suite.Assert().NotNil(reply)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestDeviceMacIsNull  设备mac地址为空
func (suite *SubmitPulseTestTestSuite) TestDeviceMacIsNull() {
	t := suite.T()
	gender := reportpb.Gender_GENDER_INVALID
	hand := reportpb.Hand_HAND_RIGHT
	data := ReadPayload(RawData)
	seconds := len(data) / fps
	now := time.Now()
	startTime := timestamppb.New(now)
	stopTime := timestamppb.New(now.Add(time.Second * time.Duration(seconds)))
	signature := string(md5.MD5(data))
	enableStatistics := true

	reply, err := suite.SubmitGeneralPulseTest(tenantId, customerId, ageIsNull, gender, staffId, hand, fps, deviceMacIsNull, deviceSn, deviceModel, codec, data, signature, startTime, stopTime, enableStatistics)
	suite.Assert().NotNil(reply)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

// TestPayloadIsErr   数据错误
func (suite *SubmitPulseTestTestSuite) TestPayloadIsErr() {
	gender := reportpb.Gender_GENDER_INVALID
	hand := reportpb.Hand_HAND_RIGHT
	data := ReadPayload(ErrData)
	seconds := len(data) / fps
	now := time.Now()
	startTime := timestamppb.New(now)
	stopTime := timestamppb.New(now.Add(time.Second * time.Duration(seconds)))
	signature := string(md5.MD5(data))
	enableStatistics := true

	reply, err := suite.SubmitGeneralPulseTest(tenantId, customerId, ageIsNull, gender, staffId, hand, fps, deviceMac, deviceSn, deviceModel, codec, data, signature, startTime, stopTime, enableStatistics)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(reply)
}

func (suite *SubmitPulseTestTestSuite) TearDownSuite() {
}

func TestSubmitPulseTestTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitPulseTestTestSuite))
}

func ReadPayload(name string) []byte {
	vals, err := ptf.ReadIntegers(name)
	if err != nil {
		fmt.Printf("Read file error: %s\n", err)
		return nil
	}

	codec := pt.NewCodec(pt.DeviceXMW23, make(map[string]string))
	b, err := codec.Encode(vals)
	if err != nil {
		fmt.Printf("Encoding error: %s\n", err)
		return nil
	}
	return b
}
