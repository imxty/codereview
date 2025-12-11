package report

import (
	"context"
	"encoding/base64"
	"encoding/json"
	gerr "errors"
	"fmt"
	"strconv"

	customerpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	devicepb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/jinmukeji/plat-pkg/v4/micro/meta"
	platformpb "github.com/jinmukeji/platform-proto/gen/go/platform/report/v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type address struct {
	CountryCode int    `json:"country_code"`
	Country     string `json:"country"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	FullAddress string `json:"full_address"`
	Ip          string `json:"ip"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
}

const (
	ContextAddress = "address"
)

const (
	// ErrInitReportFailed 初始化报告失败
	ErrInitReportFailed = 5022
	// ErrTenantDisable
	ErrTenantDisable = 5304
	// ErrPlatformError
	ErrPlatformError = 5051
)

// SubmitPulseTest 提交测量数据
func (s *ReportAPIHandler) SubmitPulseTest(ctx context.Context, req *pb.SubmitPulseTestRequest, rsp *pb.SubmitPulseTestResponse) error {

	// 1.验证request
	err := validateSubmitPulseTestRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商户信息
	getRsp, err := s.userAPI.GetTenantEntity(ctx, &userpb.GetTenantEntityRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return err
	}

	// 判断商户状态
	if getRsp.GetEntity().GetTenantStatus() != userpb.TenantStatus_TENANT_STATUS_USING {
		// 停用无法测量
		return errors.Errorf(ErrTenantDisable, "tenant[%s] has been disabled. can not submit pulse test", req.GetTenantId())
	}

	// 2.构建测量主体信息
	var age string
	var subject *platformpb.SubjectProfile
	var opt *platformpb.PulseTestOptions
	movingAverage := false
	if req.GetCustomerId() != "" {
		// 获取常客信息
		getCustomerRsp, err := s.customerAPI.GetCustomer(ctx, &customerpb.GetCustomerRequest{
			TenantId:   req.GetTenantId(),
			CustomerId: req.GetCustomerId(),
		})
		if err != nil {
			return errors.Error(codes.InvalidOperation, err.Error())
		}
		customer := getCustomerRsp.GetCustomer()
		cusAge := customer.GetAge()
		if cusAge >= 120 {
			cusAge = 120
		}
		// 构建主体信息
		subject = &platformpb.SubjectProfile{
			// 常客ID
			SubjectId: customer.GetCustomerId(),
			// 性别
			Gender: toPlatformGender(customer.GetGender()),
			// 身高，单位 cm
			Height: uint32(customer.GetHeight()),
			// 体重，单位 kg
			Weight: uint32(customer.GetWeight()),
			// 年龄
			Age: uint32(cusAge),
		}
		opt = &platformpb.PulseTestOptions{
			// 是否是散客用户测量
			// 如果是散客用户，则忽略 subject_id 与 moving_average 的值.
			// 散客用户产生的数据不纳入统计功能分析.
			IsGuestSubject: false,
			// 是否启用历史记录，启用后的可以在历史记录中查询到
			EnableHistory: true,
			// 是否启用统计，启用统计会参与周报月报统计
			EnableTrendingStatistics: true,
		}
		movingAverage = true
		age = strconv.Itoa(int(customer.GetAge()))
	} else {
		tempCustomer := req.GetTempCustomer()
		subject = &platformpb.SubjectProfile{
			// 性别
			Gender: toPlatformGenderFromReport(tempCustomer.GetGender()),
			// 年龄
			Age: uint32(mapTempCustomerAge(tempCustomer.GetAge())),
			// 测量时身高，单位 cm,必填
			Height: 170,
			// 测量时体重，单位 kg,必填
			Weight: 50,
		}
		opt = &platformpb.PulseTestOptions{
			// 是否是散客用户测量
			// 如果是散客用户，则忽略 subject_id 与 moving_average 的值.
			// 散客用户产生的数据不纳入统计功能分析.
			IsGuestSubject: true,
			// 是否启用历史记录，启用后的可以在历史记录中查询到
			EnableHistory: true,
			// 是否启用统计，启用统计会参与周报月报统计
			EnableTrendingStatistics: true,
		}
		age = tempCustomer.GetAge()
	}

	// 获取设备
	getDeviceRsp, err := s.deviceAPI.ListDevices(ctx, &devicepb.ListDevicesRequest{
		Mac: []string{
			req.GetPayload().GetSampleDevice().GetDeviceMac(),
		},
	})
	if err != nil {
		return err
	}
	if getDeviceRsp.GetDevices() == nil || len(getDeviceRsp.GetDevices()) != 1 {
		return errors.Error(codes.InvalidRequest, "fail to find device")
	}
	// 获取设备
	device := getDeviceRsp.GetDevices()[0]

	// 提交到平台
	// 构建平台提交数据请求
	submitReq, err := s.buildPlatformSubmitRequest(ctx, req, subject, opt, movingAverage, device)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 构建OutgoingContext
	ctx = buildPlatformContext(ctx)
	submitRsp, err := s.platformClient.SubmitPulseTest(ctx, submitReq)
	if err != nil {
		return errors.Errorf(ErrPlatformError, "Platform Error:[%s]", err.Error())
	}

	// 6. 初始化报告信息
	report, err := s.initReport(ctx, req.GetPayload().GetSampleData().GetSpo(), req.GetPayload().GetHand(), req.GetTenantId(), submitRsp.GetReportId(), req.GetStaffId(), age, req.GetCustomerId(), device.GetDeviceId(), getRsp.GetEntity().GetOrganizationId(), subject)
	if err != nil {
		return errors.Error(ErrInitReportFailed, err.Error())
	}

	// 7. 构建报告地址信息
	addr := meta.MustGet(ctx, ContextAddress)
	if addr != "" {
		// base64 decode
		res, err := base64.StdEncoding.DecodeString(addr)
		if err != nil {
			return errors.Errorf(codes.InvalidRequest, "invalid address[%s %s]", addr, err.Error())
		}
		ad := new(address)
		err = json.Unmarshal(res, &ad)
		if err != nil {
			return errors.Errorf(ErrInitReportFailed, "unmarshal report address failed[%s:%s]", addr, err.Error())
		}
		// 如果ip为空则获取context
		ip := ad.Ip
		if ip == "" {
			ip = meta.MustGet(ctx, "X-Forwarded-For")
		}
		rpAddr := &domain.ReportAddress{
			ReportID:    submitRsp.GetReportId(),
			Mac:         req.GetPayload().GetSampleDevice().GetDeviceMac(),
			CountryCode: ad.CountryCode,
			Country:     ad.Country,
			Province:    ad.Province,
			City:        ad.City,
			FullAddress: ad.FullAddress,
			District:    ad.District,
			Latitude:    ad.Latitude,
			Longitude:   ad.Longitude,
			Ip:          ip,
		}
		err = s.reportStore.CreateReportAddress(ctx, rpAddr)
		if err != nil {
			return errors.Errorf(ErrInitReportFailed, "create report address failed[%s]", err.Error())
		}
	}

	// 返回报告ID
	rsp.ReportId = report.GetReportID()
	return nil
}

// 验证request
func validateSubmitPulseTestRequest(req *pb.SubmitPulseTestRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetCustomerId() == "" && req.GetTempCustomer() == nil {
		return gerr.New("customer_id and temp_customer should not be nil either")
	}
	if req.GetPayload() == nil {
		return gerr.New("payload should not be nil")
	}
	payload := req.GetPayload()
	if payload.GetSampleDevice().GetDeviceMac() == "" {
		return gerr.New("device_mac should not be empty")
	}
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	return nil
}

// buildPlatformSubmitRequest
func (s *ReportAPIHandler) buildPlatformSubmitRequest(ctx context.Context, req *pb.SubmitPulseTestRequest, subject *platformpb.SubjectProfile, opt *platformpb.PulseTestOptions, movingAverage bool, device *devicepb.Device) (*platformpb.SubmitPulseTestRequest, error) {
	// 测量数据
	payload := req.GetPayload()
	geolocation := payload.GetGeoLocation()
	// 构建算法请求
	body, err := mapReportHandToPulseTest(payload.GetHand())
	if err != nil {
		return nil, fmt.Errorf("failed to map report body part to pulse test: %s", err.Error())
	}

	// 发起提交测量数据请求
	submitRequest := &platformpb.SubmitPulseTestRequest{
		AppId:   s.appId,
		Subject: subject,
		Options: opt,
		// 波形数据
		Payload: &platformpb.SamplePayload{
			Part: body,
			// 默认坐姿
			Posture: platformpb.PulseTestPosture_PULSE_TEST_POSTURE_SITTING,
			SampleDevice: &platformpb.SampleDevice{
				DeviceMac: device.GetMac(),
				// 默认外围设备形式采样设备
				DeviceType:  platformpb.DeviceType_DEVICE_TYPE_PERIPHERAL,
				DeviceModel: device.GetModel(),
			},
			SampleData: &platformpb.SampleData{
				Data:       payload.GetSampleData().GetData(),
				SampleRate: payload.GetSampleDevice().GetSampleRate(),
				Signature:  payload.GetSampleData().GetSignature(),
				StartTime:  payload.GetSamplingStartTime(),
				StopTime:   payload.GetSamplingStopTime(),
			},
		},
	}
	if req.GetPayload().GetSampleData().GetSpo() != 0 {
		submitRequest.Spo = &wrapperspb.Int32Value{
			Value: req.GetPayload().GetSampleData().GetSpo(),
		}
	}
	if geolocation != nil {
		submitRequest.Payload.GeoLocation = &platformpb.GeoLocation{
			Latitude: geolocation.GetLatitude(),
			// 经度: position's longitude in decimal degrees
			Longitude: geolocation.GetLongitude(),
			// 海拔: position's altitude in meters, relative to sea level
			Altitude: geolocation.GetAltitude(),
			// 经纬度的精度: accuracy of the latitude and longitude properties, expressed
			// in meters.
			Accuracy: geolocation.GetAccuracy(),
			// 海拔精度: accuracy of the altitude expressed in meters
			AltitudeAccuracy: geolocation.GetAltitudeAccuracy(),
			// 方向:  direction in which the device is traveling
			Heading: geolocation.GetHeading(),
			// 设备运动的速度: velocity of the device in meters per second
			Speed: geolocation.GetSpeed(),
		}
	}

	return submitRequest, nil
}

// initReport 初始化报告
func (s *ReportAPIHandler) initReport(ctx context.Context, spo int32, hand pb.Hand, tenantID, platformReportId, staffID, age, customerID, deviceId, organizationID string, subject *platformpb.SubjectProfile) (domain.ReportIntf, error) {
	// 检测是否是常客
	var isCustomer bool
	if customerID != "" {
		isCustomer = true
	}
	// 构建报告
	report := &domain.Report{
		ReportID:       platformReportId,
		OrganizationID: organizationID,
		TenantID:       tenantID,
		StaffID:        staffID,
		DeviceID:       deviceId,
		IsCustomer:     isCustomer,
		CustomerID:     customerID,
		Spo:            spo,
		Age:            age,
		Gender:         toDomainGender(subject.GetGender()),
		Hand:           toDomainHand(hand),
	}
	// 初始化报告
	err := s.reportStore.CreateReport(ctx, report)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// buildPlatformContext 构建平台context
func buildPlatformContext(ctx context.Context) context.Context {
	jwt := meta.MustGet(ctx, "Authorization")
	md := metadata.Pairs("Authorization", jwt)
	return metadata.NewOutgoingContext(ctx, md)
}
