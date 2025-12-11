package domain

type CustomerTenantRankMonthly struct {
	TenantID     string  `gorm:"column:tenant_id"`
	Month        string  `gorm:"column:month"`
	AddedCount   int64   `gorm:"column:added_count"`
	MonthOnMonth float64 `gorm:"column:month_on_month"`
	YearOnYear   float64 `gorm:"year_on_year"`
}

type CustomerTenantRank struct {
	TenantID      string
	CustomerCount int64
}

type CustomerTenantRankMapper interface {
	ToDomainCustomerTenantRankMapper
	FromDomainCustomerTenantRankMapper
}

type ToDomainCustomerTenantRankMapper interface {
	ToDomainCustomerTenantRank() CustomerTenantRankIntf
}

type FromDomainCustomerTenantRankMapper interface {
	FromDomainCustomerTenantRank(CustomerTenantRankIntf)
}

type CustomerTenantRankIntf interface {
	GetTenantID() string
	GetCustomerCount() int64
}

var _ CustomerTenantRankIntf = (*CustomerTenantRank)(nil)

func (c *CustomerTenantRank) GetTenantID() string {
	if c == nil {
		return ""
	}
	return c.TenantID
}

func (c *CustomerTenantRank) GetCustomerCount() int64 {
	if c == nil {
		return 0
	}
	return c.CustomerCount
}
