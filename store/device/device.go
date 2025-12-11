package device

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/device/domain"
	"gorm.io/gorm"
)

type Device struct {
	DeviceID    string          `gorm:"primary_key;column:device_id"`
	Mac         string          `gorm:"column:mac"`
	Sn          string          `gorm:"column:sn"`
	Model       string          `gorm:"column:model"`
	IsAvailable bool            `gorm:"column:is_available"`
	Remark      string          `gorm:"column:remark"`
	CreatedAt   time.Time       // 创建时间
	UpdatedAt   time.Time       // 更新时间
	DeletedAt   *gorm.DeletedAt // 删除时间
}

func (d Device) TableName() string {
	return "device"
}

// domain->db
func (de *Device) FromDomainDevice(d domain.DeviceIntf) {
	if de == nil || d == nil {
		return
	}

	de.DeviceID = d.GetDeviceID()
	de.Mac = d.GetMac()
	de.Sn = d.GetSn()
	de.Model = d.GetModel()
	de.IsAvailable = d.GetIsAvailable()
	de.Remark = d.GetRemark()
}

// db->domain
func (d *Device) ToDomainDevice() domain.DeviceIntf {
	if d == nil {
		return nil
	}

	p := domain.Device{
		DeviceID:    d.DeviceID,
		Mac:         d.Mac,
		Sn:          d.Sn,
		Model:       d.Model,
		IsAvailable: d.IsAvailable,
		Remark:      d.Remark,
	}
	return &p
}
