package domain

type CustomerStaffRank struct {
	StaffID       string
	TodayCustomer int64
	TotalCustomer int64
}

type CustomerStaffRankMapper interface {
	ToDomainCustomerStaffRankMapper
	FromDomainCustomerStaffRankMapper
}

type ToDomainCustomerStaffRankMapper interface {
	ToDomainCustomerStaffRank() CustomerStaffRankIntf
}

type FromDomainCustomerStaffRankMapper interface {
	FromDomainCustomerStaffRank(CustomerStaffRankIntf)
}

type CustomerStaffRankIntf interface {
	GetStaffID() string
	GetTodayCustomer() int64
	GetTotalCustomer() int64
}

var _ CustomerStaffRankIntf = (*CustomerStaffRank)(nil)

func (c *CustomerStaffRank) GetStaffID() string {
	if c == nil {
		return ""
	}
	return c.StaffID
}

func (c *CustomerStaffRank) GetTodayCustomer() int64 {
	if c == nil {
		return 0
	}
	return c.TodayCustomer
}

func (c *CustomerStaffRank) GetTotalCustomer() int64 {
	if c == nil {
		return 0
	}
	return c.TotalCustomer
}
