package customer

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/customer/domain"
	"gorm.io/gorm"
)

type Customer struct {
	CustomerID     string          `gorm:"primary_key;column:customer_id"`
	OrganizationID string          `gorm:"column:organization_id"`
	TenantID       string          `gorm:"column:tenant_id"`
	StaffID        string          `gorm:"column:staff_id"`
	Nickname       string          `gorm:"column:nickname"`
	Initial        string          `gorm:"column:initial"`
	Gender         int32           `gorm:"column:gender"`
	Phone          string          `gorm:"column:phone"`
	Dirthday       time.Time       `gorm:"column:birthday"`
	Height         int32           `gorm:"column:height"`
	Weight         int32           `gorm:"column:weight"`
	Pmh            string          `gorm:"column:pmh"`
	Remark         string          `gorm:"column:remark"`
	Rev            int32           `gorm:"column:rev"`
	CreatedAt      time.Time       // 创建时间
	UpdatedAt      time.Time       // 更新时间
	DeletedAt      *gorm.DeletedAt // 删除时间
}

func (c Customer) TableName() string {
	return "customer"
}

// domain->db
func (c *Customer) FromDomainCustomer(d domain.CustomerIntf) {
	if c == nil || d == nil {
		return
	}

	c.CustomerID = d.GetCustomerID()
	c.OrganizationID = d.GetOrganizationID()
	c.TenantID = d.GetTenantID()
	c.StaffID = d.GetStaffID()
	c.Nickname = d.GetNickname()
	c.Initial = d.GetInitial()
	c.Gender = d.GetGender()
	c.Phone = d.GetPhone()
	c.Dirthday = d.GetBirthday()
	c.Height = d.GetHeight()
	c.Weight = d.GetWeight()
	c.Pmh = d.GetPmh()
	c.Remark = d.GetRemark()
	c.Rev = d.GetRev()
}

// db->domain
func (c *Customer) ToDomainCustomer() domain.CustomerIntf {
	if c == nil {
		return nil
	}

	p := domain.Customer{
		CustomerID:     c.CustomerID,
		OrganizationID: c.OrganizationID,
		TenantID:       c.TenantID,
		StaffID:        c.StaffID,
		Nickname:       c.Nickname,
		Initial:        c.Initial,
		Gender:         c.Gender,
		Phone:          c.Phone,
		Birthday:       c.Dirthday,
		Height:         c.Height,
		Weight:         c.Weight,
		Pmh:            c.Pmh,
		Remark:         c.Remark,
		Rev:            c.Rev,
	}
	return &p
}
