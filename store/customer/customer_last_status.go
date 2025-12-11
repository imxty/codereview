package customer

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/customer/domain"
	"gorm.io/gorm"
)

type CustomerLastStatus struct {
	CustomerID string          `gorm:"primary_key;column:customer_id"`
	TenantID   string          `gorm:"column:tenant_id"`
	ReportID   string          `gorm:"column:report_id"`
	Rev        int32           `gorm:"column:rev"`
	CreatedAt  time.Time       // 创建时间
	UpdatedAt  time.Time       // 更新时间
	DeletedAt  *gorm.DeletedAt // 删除时间
}

func (c CustomerLastStatus) TableName() string {
	return "customer_last_status"
}

// domain->db
func (c *CustomerLastStatus) FromDomainCustomerLastStatus(d domain.CustomerLastStatusIntf) {
	if c == nil || d == nil {
		return
	}

	c.CustomerID = d.GetCustomerID()
	c.TenantID = d.GetTenantID()
	c.ReportID = d.GetReportID()
	c.Rev = d.GetRev()
	c.UpdatedAt = d.GetUpdatedAt()
}

// db->domain
func (c *CustomerLastStatus) ToDomainCustomerLastStatus() domain.CustomerLastStatusIntf {
	if c == nil {
		return nil
	}

	p := domain.CustomerLastStatus{
		CustomerID: c.CustomerID,
		TenantID:   c.TenantID,
		ReportID:   c.ReportID,
		Rev:        c.Rev,
		UpdatedAt:  c.UpdatedAt,
	}
	return &p
}
