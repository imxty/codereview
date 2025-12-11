package domain

import (
	"context"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
)

type DeviceRepository interface {
	dbutils.Tx
	// BatchGetDevices 批量获取设备
	BatchGetDevices(ctx context.Context, mac []string) ([]DeviceIntf, error)
}
