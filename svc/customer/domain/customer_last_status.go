package domain

import "time"

func (c *CustomerLastStatus) GetCustomerID() string {
	if c == nil {
		return ""
	}
	return c.CustomerID
}

func (c *CustomerLastStatus) GetTenantID() string {
	if c == nil {
		return ""
	}
	return c.TenantID
}

func (c *CustomerLastStatus) GetReportID() string {
	if c == nil {
		return ""
	}
	return c.ReportID
}

func (c *CustomerLastStatus) GetRev() int32 {
	if c == nil {
		return 0
	}
	return c.Rev
}

func (c *CustomerLastStatus) GetUpdatedAt() time.Time {
	if c == nil {
		return time.Time{}
	}
	return c.UpdatedAt
}

type CustomerLastStatus struct {
	CustomerID string
	TenantID   string
	ReportID   string
	Rev        int32
	UpdatedAt  time.Time
}

type CustomerLastStatusMapper interface {
	ToDomainCustomerLastStatusMapper
	FromDomainCustomerLastStatusMapper
}

type ToDomainCustomerLastStatusMapper interface {
	ToDomainCustomerLastStatus() CustomerLastStatusIntf
}

type FromDomainCustomerLastStatusMapper interface {
	FromDomainCustomerLastStatus(CustomerLastStatusIntf)
}

type CustomerLastStatusIntf interface {
	GetCustomerID() string
	GetTenantID() string
	GetReportID() string
	GetRev() int32
	GetUpdatedAt() time.Time
}

var _ CustomerLastStatusIntf = (*CustomerLastStatus)(nil)
