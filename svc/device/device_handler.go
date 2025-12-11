package device

import (
	"github.com/jinmukeji/huimaibao-service/svc/device/domain"
)

type DeviceAPIHandler struct {
	deviceStore domain.DeviceRepository
}

func (svc *DeviceAPIHandler) Name() string {
	const name = "DeviceAPI"
	return name
}

func NewDeviceAPIHandler(deviceStore domain.DeviceRepository) *DeviceAPIHandler {
	return &DeviceAPIHandler{
		deviceStore: deviceStore,
	}
}
