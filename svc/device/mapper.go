package device

import (
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/device/v1"
	"github.com/jinmukeji/huimaibao-service/svc/device/domain"
)

// toProtoDevice
func toProtoDevice(d domain.DeviceIntf) *pb.Device {
	if d == nil {
		return nil
	}
	return &pb.Device{
		// 设备ID
		DeviceId: d.GetDeviceID(),
		// 设备mac地址
		Mac: d.GetMac(),
		// 设备型号
		Model: d.GetModel(),
	}
}
