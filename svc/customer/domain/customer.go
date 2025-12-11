package domain

import "time"

type Customer struct {
	CustomerID     string
	OrganizationID string
	TenantID       string
	StaffID        string
	Nickname       string
	Initial        string
	Gender         int32
	Phone          string
	Birthday       time.Time
	Height         int32
	Weight         int32
	Pmh            string
	Remark         string
	Rev            int32
}

const (
	// 非法的性别
	ErrTipInvalidGender int32 = 0
	// 未设置性别
	UnsetGender int32 = 1
	// 男
	Male int32 = 2
	// 女
	Female int32 = 3
)

type CustomerMapper interface {
	ToDomainCustomerMapper
	FromDomainCustomerMapper
}

type ToDomainCustomerMapper interface {
	ToDomainCustomer() CustomerIntf
}

type FromDomainCustomerMapper interface {
	FromDomainCustomer(CustomerIntf)
}

type CustomerIntf interface {
	GetCustomerID() string
	GetOrganizationID() string
	GetTenantID() string
	GetStaffID() string
	GetNickname() string
	GetInitial() string
	GetGender() int32
	GetPhone() string
	GetBirthday() time.Time
	GetHeight() int32
	GetWeight() int32
	GetPmh() string
	GetRemark() string
	GetRev() int32
}

var _ CustomerIntf = (*Customer)(nil)

func (c *Customer) GetCustomerID() string {
	if c == nil {
		return ""
	}
	return c.CustomerID
}

func (c *Customer) GetOrganizationID() string {
	if c == nil {
		return ""
	}
	return c.OrganizationID
}

func (c *Customer) GetTenantID() string {
	if c == nil {
		return ""
	}
	return c.TenantID
}

func (c *Customer) GetStaffID() string {
	if c == nil {
		return ""
	}
	return c.StaffID
}

func (c *Customer) GetNickname() string {
	if c == nil {
		return ""
	}
	return c.Nickname
}

func (c *Customer) GetInitial() string {
	if c == nil {
		return ""
	}
	return c.Initial
}

func (c *Customer) GetGender() int32 {
	if c == nil {
		return 0
	}
	return c.Gender
}

func (c *Customer) GetPhone() string {
	if c == nil {
		return ""
	}
	return c.Phone
}

func (c *Customer) GetBirthday() time.Time {
	if c == nil {
		return time.Time{}
	}
	return c.Birthday
}

func (c *Customer) GetHeight() int32 {
	if c == nil {
		return 0
	}
	return c.Height
}

func (c *Customer) GetWeight() int32 {
	if c == nil {
		return 0
	}
	return c.Weight
}

func (c *Customer) GetPmh() string {
	if c == nil {
		return ""
	}
	return c.Pmh
}

func (c *Customer) GetRemark() string {
	if c == nil {
		return ""
	}
	return c.Remark
}

func (c *Customer) GetRev() int32 {
	if c == nil {
		return 0
	}
	return c.Rev
}
