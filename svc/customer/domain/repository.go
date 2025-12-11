package domain

import (
	"time"

	"context"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
)

type CustomerRepository interface {
	dbutils.Tx

	// ListCustomers 获取常客列表 (按首字母排序).
	ListCustomers(ctx context.Context, tenantId string, offset, size int) ([]CustomerIntf, error)
	// AddCustomer 添加常客.
	AddCustomer(ctx context.Context, tenantId string, customer CustomerIntf) error
	// GetCustomer 获取常客信息.
	GetCustomer(ctx context.Context, tenantId, customerId string) (CustomerIntf, error)
	// SearchCustomers 搜索常客.
	SearchCustomers(ctx context.Context, tenantId string, keywords string, offset, size int) ([]CustomerIntf, error)
	// UpdateCustomer 更新常客档案
	UpdateCustomer(ctx context.Context, tenantId, customerId string, rev int32, updates map[string]interface{}) error
	// DeleteCustomer 删除常客档案
	DeleteCustomer(ctx context.Context, tenantId, customerId string, rev int32) error
	// CheckCustomerPhoneExist 查询常客手机是否存在
	CheckCustomerPhoneExist(ctx context.Context, tenantId, phone, areaCode string) (bool, error)

	// GetCustomerStatistics 获取常客统计(当日新增常客数量,累计常客数量,error)
	GetCustomerStatistics(ctx context.Context, tenantId string, startTime, endTime time.Time) (int64, int64, error)
	// GetOrganizationCustomerStatistics 获取常客统计(当日新增常客数量,累计常客数量,error)
	GetOrganizationCustomerStatistics(ctx context.Context, tenantIds []string, startTime, endTime time.Time) (int64, int64, error)
	// BatchGetTenantCustomerStatistics 获取商户的常客统计,返回为[商户id]常客数量
	BatchGetTenantCustomerStatistics(ctx context.Context, tenantIds []string, startTime, endTime time.Time) (map[string]int64, error)
	// GetTodayStaffCustomerCount 获取今日员工客户数量
	GetTodayStaffCustomerCount(ctx context.Context, tenantID string, staffIds []string, startTime, endTime time.Time) (map[string]int32, error)
	// GetTotalStaffCustomerCount 获取全部员工客户数量
	GetTotalStaffCustomerCount(ctx context.Context, tenantID string, staffIds []string) (map[string]int32, error)
	// GetCustomerCount 获取常客总数
	GetCustomerCount(ctx context.Context, tenantID string) (int32, error)
	// ListCustomerTotalCount 获取常客总数
	ListCustomerTotalCount(ctx context.Context, tenantID []string) (int32, error)
	// ListTodayCustomerCount 获取当日常客总数
	ListTodayCustomerCount(ctx context.Context, tenantID []string, startTime, endTime time.Time) (int32, error)
	// ListCustomerCount 获取商户常客数量
	ListCustomerCount(ctx context.Context, tenantID []string, startTime, endTime time.Time) ([]CustomerTenantRankIntf, error)
	// GetTenantAddedCustomerCount 获取商户新增常客
	GetTenantAddedCustomerCount(ctx context.Context, tenant_id string, start_time, end_time time.Time) (int64, error)
	// GetOrganizationAddedCustomerCount 获取组织新增常客
	GetOrganizationAddedCustomerCount(ctx context.Context, organization_id string, start_time, end_time time.Time) (int64, error)
	// GetOrganizationOverdueCount 获取组织待复查人数
	GetOrganizationOverdueCount(ctx context.Context, organization_id string, days int32) (int64, error)
	// GetTenantOverdueCount 获取商户待复查人数
	GetTenantOverdueCount(ctx context.Context, tenant_id string, days int32) (int64, error)
	// BatchGetTenantsAddedCustomerCount 批量获取商户新增常客数量
	BatchGetTenantsAddedCustomerCount(ctx context.Context, tenant_ids []string, start_time, end_time time.Time) (map[string]int64, error)
	// BatchGetCustomersLastStatus 获取客户最新测量情况
	BatchGetCustomersLastStatus(ctx context.Context, tenant_ids []string, contact_name, contact_phone string, overdue_mask int32, offset, size int32) ([]CustomerLastStatusIntf, int64, error)
	// BatchGetCustomers 批量获取常客信息
	BatchGetCustomers(ctx context.Context, customer_ids []string, include_deleted bool) ([]CustomerIntf, error)
	// BatchGetTenantMonthlyAddedCustomerCount 获取商户每月常客新增数
	BatchGetTenantMonthlyAddedCustomerCount(ctx context.Context, tenant_ids []string, start_time, end_time time.Time, offset, limit int32) ([]CustomerTenantRankMonthly, int64, error)
	// SearchTenantOverdueCustomers 搜索商户待复查客户
	SearchTenantOverdueCustomers(ctx context.Context, tenant_id, key string, overdue, offset, size int32) ([]CustomerIntf, int64, error)
	// ListCustomerLastStatus 获取常客最新测量数据
	ListCustomerLastStatus(ctx context.Context, customer_ids []string) ([]CustomerLastStatusIntf, error)
	// SearchCustomerIDByPhoneAndName 根据常客手机号和姓名搜索常客 ID
	SearchCustomerIDByPhoneAndName(ctx context.Context, phone, name string) (string, error)
	// CreateCustomerLastStatus 创建客户最新测量状态
	CreateCustomerLastStatus(ctx context.Context, status CustomerLastStatusIntf) error
	// UpdateCustomerLastStatus 更新客户最新测量状态
	UpdateCustomerLastStatus(ctx context.Context, customer_id, report_id string) error
	// SearchAppTenantCustomers 搜索 App 常客数据
	SearchAppTenantCustomers(ctx context.Context, tenant_id, key string, size, offset int32) ([]CustomerIntf, int32, error)
	// BatchGetAppCustomersLastStatus 获取客户最新测量情况
	BatchGetAppCustomersLastStatus(ctx context.Context, tenant_id string, key string, overdue_mask int32, offset, size int32) ([]CustomerLastStatusIntf, int64, error)
}
