package domain

type Device struct {
	DeviceID    string
	Mac         string
	Sn          string
	Model       string
	IsAvailable bool
	Remark      string
}

type DeviceMapper interface {
	ToDomainDeviceMapper
	FromDomainDeviceMapper
}

type ToDomainDeviceMapper interface {
	ToDomainDevice() DeviceIntf
}

type FromDomainDeviceMapper interface {
	FromDomainDevice(DeviceIntf)
}

type DeviceIntf interface {
	GetDeviceID() string
	GetMac() string
	GetSn() string
	GetModel() string
	GetIsAvailable() bool
	GetRemark() string
}

var _ DeviceIntf = (*Device)(nil)

func (d *Device) GetDeviceID() string {
	if d == nil {
		return ""
	}
	return d.DeviceID
}

func (d *Device) GetMac() string {
	if d == nil {
		return ""
	}
	return d.Mac
}

func (d *Device) GetSn() string {
	if d == nil {
		return ""
	}
	return d.Sn
}

func (d *Device) GetModel() string {
	if d == nil {
		return ""
	}
	return d.Model
}

func (d *Device) GetIsAvailable() bool {
	if d == nil {
		return false
	}
	return d.IsAvailable
}

func (d *Device) GetRemark() string {
	if d == nil {
		return ""
	}
	return d.Remark
}
