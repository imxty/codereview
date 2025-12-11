package app

import (
	"context"
	gerr "errors"

	"github.com/jinmukeji/go-pkg/v2/crypto/hash"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) SubmitPulseTest(ctx context.Context, req *pb.SubmitPulseTestRequest, rsp *pb.SubmitPulseTestResponse) error {
	err := validateSubmitPulseTestRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 验证签名
	sha256Hash := hash.HexString(hash.SHA256(req.GetPayload().GetSampleData().GetData()))
	if sha256Hash != req.GetPayload().GetSampleData().GetSignature() {
		return errors.Error(api.ErrSignatureError, "invalid signature")
	}
	// 提交测量数据
	submitptReq := &reportpb.SubmitPulseTestRequest{
		TenantId:     req.GetTenantId(),
		CustomerId:   req.GetCustomerId(),
		StaffId:      req.GetStaffId(),
		TempCustomer: toReportTempCustomer(req.GetTempCustomer()),
		Payload: &reportpb.SamplePayload{
			Hand: toSrvHand(req.GetPayload().GetHand()),
			SampleDevice: &reportpb.SampleDevice{
				DeviceMac:    req.GetPayload().GetSampleDevice().GetDeviceMac(),
				DeviceModel:  req.GetPayload().GetSampleDevice().GetDeviceModel(),
				DeviceParams: req.GetPayload().GetSampleDevice().GetDeviceParams(),
				DeviceSn:     req.GetPayload().GetSampleDevice().GetDeviceSn(),
				SampleRate:   req.GetPayload().GetSampleDevice().GetSampleRate(),
			},
			SampleData: &reportpb.SampleData{
				Codec:       req.GetPayload().GetSampleData().GetCodec(),
				CodecParams: req.GetPayload().GetSampleData().GetCodecParams(),
				Data:        req.GetPayload().GetSampleData().GetData(),
				Signature:   req.GetPayload().GetSampleData().GetSignature(),
				Spo:         req.GetPayload().GetSampleData().GetSpo(),
			},
			SamplingStartTime: req.GetPayload().GetSamplingStartTime(),
			SamplingStopTime:  req.GetPayload().GetSamplingStopTime(),
		},
	}
	if req.GetTempCustomer() != nil {
		submitptReq.TempCustomer = &reportpb.TempCustomer{
			// 年龄
			Age: req.GetTempCustomer().GetAge(),
			// 性别
			Gender: toSvcGender(req.GetTempCustomer().GetGender()),
		}
	}
	// 地理位置信息
	if req.GetPayload().GetGeoLocation() != nil {
		location := req.GetPayload().GetGeoLocation()
		submitptReq.Payload.GeoLocation = &reportpb.GeoLocation{
			Latitude:         location.GetLatitude(),
			Longitude:        location.GetLongitude(),
			Altitude:         location.GetAltitude(),
			Accuracy:         location.GetAccuracy(),
			AltitudeAccuracy: location.GetAltitudeAccuracy(),
			Heading:          location.GetHeading(),
			Speed:            location.GetSpeed(),
		}
	}
	// 发送请求
	submitptRsp, err := s.reportAPI.SubmitPulseTest(ctx, submitptReq)
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}
	rsp.ReportId = submitptRsp.GetReportId()
	return nil
}

// 验证request
func validateSubmitPulseTestRequest(req *pb.SubmitPulseTestRequest) error {
	if req.GetPayload() == nil {
		return gerr.New("payload should not be empty")
	}
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetTempCustomer() != nil {
		if req.GetTempCustomer().GetGender() == pb.Gender_GENDER_UNSET || req.GetTempCustomer().GetGender() == pb.Gender_GENDER_INVALID {
			return gerr.New("invalid gender")
		}
		if req.GetTempCustomer().GetAge() == "" {
			return gerr.New("age should not be empty")
		}
	}
	if req.GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	payload := req.GetPayload()
	if payload.GetSampleDevice().GetDeviceMac() == "" {
		return gerr.New("device_mac should not be empty")
	}
	if payload.GetSampleDevice().GetDeviceModel() == "" {
		return gerr.New("device_model should not be empty")
	}
	return nil
}
