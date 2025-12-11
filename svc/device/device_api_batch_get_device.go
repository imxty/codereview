package device

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *DeviceAPIHandler) ListDevices(ctx context.Context, req *pb.ListDevicesRequest, rsp *pb.ListDevicesResponse) error {

	// 1.验证request
	err := validateListDevicesRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 2.查询设备
	devices, err := s.deviceStore.BatchGetDevices(ctx, req.GetMac())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回设备
	appDevices := make([]*pb.Device, len(devices))
	for k, v := range devices {
		appDevices[k] = toProtoDevice(v)
	}

	rsp.Devices = appDevices

	return nil
}

// 验证request
func validateListDevicesRequest(req *pb.ListDevicesRequest) error {
	if req.GetMac() == nil {
		return gerr.New("devices should not be nil")
	}
	return nil
}
