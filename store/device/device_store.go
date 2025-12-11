package device

import (
	"context"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/huimaibao-service/svc/device/domain"
)

// DeviceStore 实现 Domain 层 DeviceStore 接口
type DeviceStore struct {
	// 数据库连接
	*dbutils.Connection
}

func NewDeviceStore(management *dbutils.Connection) *DeviceStore {
	return &DeviceStore{
		management,
	}
}

// 实现 Domain 行为
var _ domain.DeviceRepository = (*DeviceStore)(nil)

// BatchGetDevices 批量获取设备
func (d *DeviceStore) BatchGetDevices(ctx context.Context, mac []string) ([]domain.DeviceIntf, error) {
	db := d.GetConnection(ctx)
	var devices []Device
	// 查询设备
	err := db.Model(&Device{}).Where("mac IN (?)", mac).Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	// 转化
	dDevices := make([]domain.DeviceIntf, len(devices))
	for k, v := range devices {
		dDevices[k] = v.ToDomainDevice()
	}
	return dDevices, nil
}
