package customer

import (
	"context"
	"errors"
	"time"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/huimaibao-service/svc/customer/domain"
	"gorm.io/gorm"
)

// CustomerStore 实现 Domain 层 CustomerStore 接口
type CustomerStore struct {
	// 数据库连接
	*dbutils.Connection
}

// 实现 Domain 行为
var _ domain.CustomerRepository = (*CustomerStore)(nil)

func NewCustomerStore(management *dbutils.Connection) *CustomerStore {
	return &CustomerStore{
		management,
	}
}

// BatchGetAppCustomersLastStatus 获取客户最新测量情况
func (c *CustomerStore) BatchGetAppCustomersLastStatus(ctx context.Context, tenant_id string, key string, overdue_mask int32, offset, size int32) ([]domain.CustomerLastStatusIntf, int64, error) {
	db := c.GetConnection(ctx)
	var results []CustomerLastStatus
	var total int64

	whereClause := `c.tenant_id = ? AND (c.nickname LIKE ?  OR c.phone = ?) AND c.deleted_at IS NULL AND cls.report_id != '' AND DATEDIFF(NOW(), cls.updated_at) >= ?`
	args := []interface{}{tenant_id, "%" + key + "%", key, overdue_mask}

	countSql := `SELECT COUNT(*) 
	FROM customer c 
	LEFT JOIN customer_last_status cls
	ON c.customer_id = cls.customer_id
	WHERE ` + whereClause
	err := db.Raw(countSql, args...).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	dataSql := `SELECT c.customer_id, c.tenant_id, cls.report_id, cls.rev, cls.updated_at
	FROM customer c
	LEFT JOIN customer_last_status cls
	ON c.customer_id = cls.customer_id
	WHERE ` + whereClause + ` LIMIT ? OFFSET ?`
	args = append(args, size, offset)
	err = db.Raw(dataSql, args...).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	dresults := make([]domain.CustomerLastStatusIntf, len(results))
	for k, v := range results {
		dresults[k] = v.ToDomainCustomerLastStatus()
	}

	return dresults, total, nil
}

// SearchAppTenantCustomers 搜索 App 常客数据
func (c *CustomerStore) SearchAppTenantCustomers(ctx context.Context, tenant_id, key string, size, offset int32) ([]domain.CustomerIntf, int32, error) {
	db := c.GetConnection(ctx)
	var cs []Customer
	var count int64
	likeKey := "%" + key + "%"
	err := db.Model(&Customer{}).Where(`tenant_id = ? AND (phone LIKE ? OR nickname LIKE ?)`, tenant_id, likeKey, likeKey).Limit(int(size)).Offset(int(offset)).Find(&cs).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Model(&Customer{}).Where(`tenant_id = ? AND (phone LIKE ? OR nickname LIKE ?)`, tenant_id, likeKey, likeKey).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	dcs := make([]domain.CustomerIntf, len(cs))
	for k, v := range cs {
		dcs[k] = v.ToDomainCustomer()
	}

	return dcs, int32(count), nil
}

// CreateCustomerLastStatus 创建客户最新测量状态
func (c *CustomerStore) CreateCustomerLastStatus(ctx context.Context, status domain.CustomerLastStatusIntf) error {
	db := c.GetConnection(ctx)
	var newRecord CustomerLastStatus
	newRecord.FromDomainCustomerLastStatus(status)
	return db.Model(&CustomerLastStatus{}).Create(&newRecord).Error
}

// UpdateCustomerLastStatus 更新客户最新测量状态
func (c *CustomerStore) UpdateCustomerLastStatus(ctx context.Context, customer_id, report_id string) error {
	db := c.GetConnection(ctx)
	return db.Model(&CustomerLastStatus{}).Where("customer_id = ?", customer_id).Update("report_id", report_id).Error
}

// GetTenantAddedCustomerCount 获取商户新增常客数量
func (c *CustomerStore) GetTenantAddedCustomerCount(ctx context.Context, tenant_id string, start_time time.Time, end_time time.Time) (int64, error) {
	db := c.GetConnection(ctx)
	var count int64
	err := db.Model(&Customer{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenant_id, start_time, end_time).Count(&count).Error
	return count, err
}

// GetOrganizationAddedCustomerCount 获取组织新增常客
func (c *CustomerStore) GetOrganizationAddedCustomerCount(ctx context.Context, organization_id string, start_time, end_time time.Time) (int64, error) {
	db := c.GetConnection(ctx)
	var count int64
	err := db.Model(&Customer{}).Where("organization_id = ? AND created_at >= ? AND created_at < ?", organization_id, start_time, end_time).Count(&count).Error
	return count, err
}

// GetOrganizationOverdueCount 获取组织待复查人数
func (c *CustomerStore) GetOrganizationOverdueCount(ctx context.Context, organization_id string, days int32) (int64, error) {
	db := c.GetConnection(ctx)
	var count int64
	err := db.Raw(`SELECT COUNT(*) AS overdue_customers
	FROM customer c
	LEFT JOIN customer_last_status cls ON c.customer_id = cls.customer_id
	WHERE c.organization_id = ?
	AND cls.updated_at < DATE_SUB(NOW(), INTERVAL ? DAY)
	AND c.deleted_at IS NULL`, organization_id, days).Count(&count).Error
	return count, err
}

// GetTenantOverdueCount 获取商户待复查人数
func (c *CustomerStore) GetTenantOverdueCount(ctx context.Context, tenant_id string, days int32) (int64, error) {
	db := c.GetConnection(ctx)
	var count int64
	err := db.Raw(`SELECT COUNT(*)
	FROM customer_last_status
	WHERE tenant_id = ?
	AND updated_at < DATE_SUB(NOW(), INTERVAL ? DAY)
	AND deleted_at IS NULL`, tenant_id, days).Count(&count).Error
	return count, err
}

// BatchGetTenantsAddedCustomerCount 批量获取商户新增常客数量
func (c *CustomerStore) BatchGetTenantsAddedCustomerCount(ctx context.Context, tenant_ids []string, start_time, end_time time.Time) (map[string]int64, error) {
	db := c.GetConnection(ctx)
	var results []struct {
		TenantID string `gorm:"column:tenant_id"`
		Count    int64  `gorm:"column:added_customers_count"`
	}
	err := db.Raw(`SELECT tenant_id, COUNT(*) AS added_customers_count
	FROM customer
	WHERE tenant_id IN (?)
	AND created_at >= ?
	AND created_at < ?
	AND deleted_at IS NULL
	GROUP BY tenant_id`, tenant_ids, start_time, end_time).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	addedMap := make(map[string]int64)
	for _, result := range results {
		addedMap[result.TenantID] = result.Count
	}

	return addedMap, nil
}

// BatchGetCustomersLastStatus 获取客户最新测量情况
func (c *CustomerStore) BatchGetCustomersLastStatus(ctx context.Context, tenant_ids []string, contact_name, contact_phone string, overdue_mask int32, offset, size int32) ([]domain.CustomerLastStatusIntf, int64, error) {
	db := c.GetConnection(ctx)
	var results []CustomerLastStatus
	var total int64

	whereClause := `c.tenant_id IN (?) AND c.nickname LIKE ? AND c.deleted_at IS NULL AND cls.report_id != '' AND DATEDIFF(NOW(), cls.updated_at) >= ?`
	args := []interface{}{tenant_ids, "%" + contact_name + "%", overdue_mask}

	if contact_phone != "" {
		whereClause += " AND c.phone = ?"
		args = append(args, contact_phone)
	}

	countSql := `SELECT COUNT(*) 
	FROM customer c 
	LEFT JOIN customer_last_status cls
	ON c.customer_id = cls.customer_id
	WHERE ` + whereClause
	err := db.Raw(countSql, args...).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	dataSql := `SELECT c.customer_id, c.tenant_id, cls.report_id, cls.rev, cls.updated_at
	FROM customer c
	LEFT JOIN customer_last_status cls
	ON c.customer_id = cls.customer_id
	WHERE ` + whereClause + ` LIMIT ? OFFSET ?`
	args = append(args, size, offset)
	err = db.Raw(dataSql, args...).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	dresults := make([]domain.CustomerLastStatusIntf, len(results))
	for k, v := range results {
		dresults[k] = v.ToDomainCustomerLastStatus()
	}

	return dresults, total, nil
}

// BatchGetCustomers 批量获取常客信息
func (c *CustomerStore) BatchGetCustomers(ctx context.Context, customer_ids []string, include_deleted bool) ([]domain.CustomerIntf, error) {
	db := c.GetConnection(ctx)
	var cs []Customer
	query := `SELECT * FROM customer where customer_id IN (?)`
	if !include_deleted {
		query += " AND deleted_at IS NULL"
	}
	err := db.Raw(query, customer_ids).Scan(&cs).Error
	if err != nil {
		return nil, err
	}

	dcs := make([]domain.CustomerIntf, len(cs))
	for k, v := range cs {
		dcs[k] = v.ToDomainCustomer()
	}

	return dcs, nil
}

// BatchGetTenantMonthlyAddedCustomerCount 获取商户每月常客新增数
func (c *CustomerStore) BatchGetTenantMonthlyAddedCustomerCount(ctx context.Context, tenant_ids []string, start_time, end_time time.Time, offset, size int32) ([]domain.CustomerTenantRankMonthly, int64, error) {
	db := c.GetConnection(ctx)
	var count int64
	var results []domain.CustomerTenantRankMonthly

	query := `SELECT 
		r1.tenant_id,
		DATE_FORMAT(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), '%Y-%m') AS month,
		COUNT(r1.customer_id) AS added_count,
		-- 环比：当前月和前一个月的比较
		IFNULL((
			COUNT(r1.customer_id) - 
			(SELECT COUNT(r2.customer_id)
			FROM customer r2
			WHERE r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 1 MONTH), '%Y-%m')
			)
		) / NULLIF(
			(SELECT COUNT(r2.customer_id)
			FROM customer r2
			WHERE r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 1 MONTH), '%Y-%m')
			), 0) * 100, 2147483647) AS month_on_month,

		-- 同比：当前月和前一年的比较
		IFNULL((
			COUNT(r1.customer_id) - 
			(SELECT COUNT(r2.customer_id)
			FROM customer r2
			WHERE r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 12 MONTH), '%Y-%m')
			)
		) / NULLIF(
			(SELECT COUNT(r2.customer_id)
			FROM customer r2
			WHERE r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 12 MONTH), '%Y-%m')
			), 0) * 100, 2147483647) AS year_on_year

	FROM 
		customer r1
	WHERE 
		r1.tenant_id IN (?) 
		AND r1.created_at between ? and ?
	GROUP BY 
		r1.tenant_id, month
	ORDER BY 
		r1.tenant_id, month
	LIMIT ?
	OFFSET ?`
	err := db.Raw(query, tenant_ids, start_time, end_time, size, offset).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Raw(`SELECT COUNT(*) FROM
	(SELECT 
		r1.tenant_id,
		DATE_FORMAT(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), '%Y-%m') AS month,
		COUNT(r1.customer_id) AS added_count,
		-- 环比：当前月和前一个月的比较
		IFNULL((
			COUNT(r1.customer_id) - 
			(SELECT COUNT(r2.customer_id)
			FROM customer r2
			WHERE r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 1 MONTH), '%Y-%m')
			)
		) / NULLIF(
			(SELECT COUNT(r2.customer_id)
			FROM customer r2
			WHERE r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 1 MONTH), '%Y-%m')
			), 0) * 100, 2147483647) AS month_on_month,

		-- 同比：当前月和前一年的比较
		IFNULL((
			COUNT(r1.customer_id) - 
			(SELECT COUNT(r2.customer_id)
			FROM customer r2
			WHERE r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 12 MONTH), '%Y-%m')
			)
		) / NULLIF(
			(SELECT COUNT(r2.customer_id)
			FROM customer r2
			WHERE r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 12 MONTH), '%Y-%m')
			), 0) * 100, 2147483647) AS year_on_year

	FROM 
		customer r1
	WHERE 
		r1.tenant_id IN (?) 
		AND r1.created_at between ? and ?
	GROUP BY 
		r1.tenant_id, month
	ORDER BY 
		r1.tenant_id, month) a`, tenant_ids, start_time, end_time).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

// SearchTenantOverdueCustomers 搜索商户待复查客户
func (c *CustomerStore) SearchTenantOverdueCustomers(ctx context.Context, tenant_id, key string, overdue, offset, size int32) ([]domain.CustomerIntf, int64, error) {
	db := c.GetConnection(ctx)
	var count int64
	var results []Customer
	err := db.Raw(`SELECT c.*
	FROM customer c
	LEFT JOIN
		customer_last_status cls ON c.customer_id = cls.customer_id
	WHERE
	c.tenant_id = ?
	AND DATEDIFF(NOW(), cls.updated_at) > ?
	AND ? = '' OR c.nickname LIKE ? OR c.phone = ?
	AND c.deleted_at IS NULL
	ORDER BY cls.updated_at DESC
	LIMIT ? OFFSET ?`, tenant_id, overdue, key, key, key, size, offset).Scan(&results).Error
	if err != nil {
		return nil, count, err
	}

	err = db.Raw(`SELECT COUNT(*) FROM(
		SELECT c.*
		FROM customer c
		LEFT JOIN
			customer_last_status cls ON c.customer_id = cls.customer_id
		WHERE
		c.tenant_id = ?
		AND DATEDIFF(NOW(), cls.updated_at) > ?
		AND c.deleted_at IS NULL
		AND ? = '' OR c.nickname LIKE ? OR c.phone = ?
	) a`, tenant_id, overdue, key, key, key).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	dcs := make([]domain.CustomerIntf, len(results))
	for k, v := range results {
		dcs[k] = v.ToDomainCustomer()
	}

	return dcs, count, nil
}

// ListCustomerLastStatus 获取常客最新测量数据
func (c *CustomerStore) ListCustomerLastStatus(ctx context.Context, customer_ids []string) ([]domain.CustomerLastStatusIntf, error) {
	db := c.GetConnection(ctx)
	var cs []CustomerLastStatus
	querySql := `SELECT * FROM customer_last_status WHERE customer_id IN (?) AND deleted_at IS NULL`
	err := db.Raw(querySql, customer_ids).Find(&cs).Error
	if err != nil {
		return nil, err
	}

	dcs := make([]domain.CustomerLastStatusIntf, len(cs))
	for k, v := range cs {
		dcs[k] = v.ToDomainCustomerLastStatus()
	}

	return dcs, nil
}

// SearchCustomerIDByPhoneAndName 根据常客手机号和姓名搜索常客 ID
func (c *CustomerStore) SearchCustomerIDByPhoneAndName(ctx context.Context, phone, name string) (string, error) {
	db := c.GetConnection(ctx)
	var result string
	err := db.Raw(`SELECT customer_id
	FROM customer
	WHERE (? = '' OR phone = ?)
	AND (? = '' OR nickname LIKE ?)
	AND deleted_at IS NULL
	LIMIT 1`, phone, phone, name, "%"+name+"%").Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
	}
	return result, err
}

// ListCustomers 获取常客列表 (按首字母排序).
func (c *CustomerStore) ListCustomers(ctx context.Context, tenantId string, offset, size int) ([]domain.CustomerIntf, error) {
	var customers []Customer
	// 查询所有常客信息
	err := c.GetConnection(ctx).Raw(`
    SELECT 
        customer.customer_id,
        customer.tenant_id,
        customer.staff_id,
        customer.nickname,
        customer.initial,
        customer.gender,
        customer.phone,
        customer.birthday,
        customer.height,
        customer.weight,
        customer.pmh,
        customer.remark,
        customer.rev,
        customer.customer_id,
        customer.created_at,
        customer.updated_at,
        customer.deleted_at
    FROM
        customer
    WHERE
        tenant_id = ?
            AND customer.deleted_at IS NULL
    ORDER BY CASE
        WHEN customer.initial = '#' THEN 2
        ELSE 1
    END , customer.initial , CONVERT( customer.nickname USING GBK) ASC , customer.created_at , customer.customer_id
    LIMIT ? OFFSET ?
    `, tenantId, size, offset).Find(&customers).Error
	if err != nil {
		return nil, err
	}
	// 转化
	dCustomers := make([]domain.CustomerIntf, len(customers))
	for k, v := range customers {
		dCustomers[k] = v.ToDomainCustomer()
	}
	return dCustomers, nil
}

// AddCustomer 添加常客.
func (c *CustomerStore) AddCustomer(ctx context.Context, tenantId string, customer domain.CustomerIntf) error {
	// 创建db常客
	var cs Customer
	cs.FromDomainCustomer(customer)
	// 数据库创建常客
	return c.GetConnection(ctx).Create(&cs).Error
}

// GetCustomer 获取常客信息.
func (c *CustomerStore) GetCustomer(ctx context.Context, tenantId, customerId string) (domain.CustomerIntf, error) {
	var customer Customer
	// 查询常客信息
	err := c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id = ? AND customer_id = ?", tenantId, customerId).First(&customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	// 转化
	dCustomer := customer.ToDomainCustomer()

	return dCustomer, nil
}

// SearchCustomers 搜索常客.
func (c *CustomerStore) SearchCustomers(ctx context.Context, tenantId string, keywords string, offset, size int) ([]domain.CustomerIntf, error) {
	var customers []Customer
	// keywords为空 返回空
	if keywords == "" {
		// 查询所有常客信息
		err := c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id = ?", tenantId).Limit(size).Offset(offset).Find(&customers).Error
		if err != nil {
			return nil, err
		}
		// 转化
		dCustomers := make([]domain.CustomerIntf, len(customers))
		for k, v := range customers {
			dCustomers[k] = v.ToDomainCustomer()
		}
		return dCustomers, nil
	}

	// 模糊查询关键词
	phoneKeywords := keywords + "%"
	nicknameKeywords := "%" + keywords + "%"
	// 查询所有常客信息
	err := c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id = ? AND (phone LIKE ? OR nickname LIKE ?) ", tenantId, phoneKeywords, nicknameKeywords).Limit(size).Offset(offset).Find(&customers).Error
	if err != nil {
		return nil, err
	}
	// 转化
	dCustomers := make([]domain.CustomerIntf, len(customers))
	for k, v := range customers {
		dCustomers[k] = v.ToDomainCustomer()
	}
	return dCustomers, nil
}

// UpdateCustomer 更新常客档案
func (c *CustomerStore) UpdateCustomer(ctx context.Context, tenantId, customerId string, rev int32, updates map[string]interface{}) error {
	// 更新字段 版本号加1
	updates["rev"] = rev + 1
	err := c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id = ? AND customer_id = ? AND rev = ?", tenantId, customerId, rev).
		Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteCustomer 删除常客档案
func (c *CustomerStore) DeleteCustomer(ctx context.Context, tenantId, customerId string, rev int32) error {
	return c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id = ? AND customer_id = ? AND rev = ?", tenantId, customerId, rev).
		Updates(map[string]interface{}{
			"deleted_at": time.Now().UTC(),
			"rev":        rev + 1,
		}).Error
}

// CheckCustomerPhoneExist 查询常客手机是否存在
func (c *CustomerStore) CheckCustomerPhoneExist(ctx context.Context, tenantId, phone, areaCode string) (bool, error) {
	// 查询手机号
	var customer Customer
	err := c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id = ? AND phone = ?", tenantId, phone).First(&customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

// GetCustomerStatistics 获取常客统计(当日新增常客数量,累计常客数量,error)
func (c *CustomerStore) GetCustomerStatistics(ctx context.Context, tenantId string, startTime, endTime time.Time) (int64, int64, error) {
	// 统计
	var totalCount, toadyAdd int64
	err := c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id = ?", tenantId).Count(&totalCount).Error
	if err != nil {
		return 0, 0, err
	}
	err = c.GetConnection(ctx).Model(&Customer{}).Where(`tenant_id = ? AND created_at BETWEEN ? AND ?`, tenantId, startTime, endTime).Count(&toadyAdd).Error
	if err != nil {
		return 0, 0, err
	}
	return totalCount, toadyAdd, nil
}

// GetOrganizationCustomerStatistics 获取常客统计(当日新增常客数量,累计常客数量,error)
func (c *CustomerStore) GetOrganizationCustomerStatistics(ctx context.Context, tenantIds []string, startTime, endTime time.Time) (int64, int64, error) {
	// 统计
	var totalCount, toadyAdd int64
	err := c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id IN (?)", tenantIds).Count(&totalCount).Error
	if err != nil {
		return 0, 0, err
	}
	err = c.GetConnection(ctx).Model(&Customer{}).Where(`tenant_id IN (?) AND created_at BETWEEN ? AND ?`, tenantIds, startTime, endTime).Count(&toadyAdd).Error
	if err != nil {
		return 0, 0, err
	}
	return totalCount, toadyAdd, nil
}

// BatchGetTenantCustomerStatistics 获取商户的常客统计,返回为[商户id]常客数量
func (c *CustomerStore) BatchGetTenantCustomerStatistics(ctx context.Context, tenantIds []string, startTime, endTime time.Time) (map[string]int64, error) {
	db := c.GetConnection(ctx)
	var results []struct {
		TenantID      string `gorm:"column:tenant_id"`
		CustomerCount int64  `gorm:"column:customer_count"`
	}
	err := db.Raw(`SELECT tenant_id, COUNT(*) AS customer_count
	FROM customer
	WHERE tenant_id IN (?)
	AND created_at BETWEEN ? AND ?
	GROUP BY tenant_id`, tenantIds, startTime, endTime).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	drs := make(map[string]int64)
	for _, each := range results {
		drs[each.TenantID] = each.CustomerCount
	}

	return drs, nil
}

// GetTodayStaffCustomerCount 获取今日员工客户统计
func (c *CustomerStore) GetTodayStaffCustomerCount(ctx context.Context, tenantID string, staffIds []string, startTime, endTime time.Time) (map[string]int32, error) {
	db := c.GetConnection(ctx)

	var results []struct {
		StaffID         string `gorm:"column:staff_id"`
		TodayAddedCount int32  `gorm:"column:today_added_count"`
	}

	err := db.Raw(`SELECT
		staff_id,
		SUM(CASE WHEN DATE(created_at) = CURDATE() THEN 1 ELSE 0 END) AS today_added_count
		FROM customer
		WHERE tenant_id = ?
		AND staff_id IN (?)
		AND created_at BETWEEN ? AND ?
		AND deleted_at IS NULL
		GROUP BY staff_id`, tenantID, staffIds, startTime, endTime).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	rs := make(map[string]int32)
	for _, each := range results {
		rs[each.StaffID] = each.TodayAddedCount
	}

	return rs, nil
}

// GetTotalStaffCustomerCount 获取全部员工客户数量
func (c *CustomerStore) GetTotalStaffCustomerCount(ctx context.Context, tenantID string, staffIds []string) (map[string]int32, error) {
	db := c.GetConnection(ctx)
	var results []struct {
		StaffID         string `gorm:"column:staff_id"`
		TotalAddedCount int32  `gorm:"column:total_added_count"`
	}

	err := db.Raw(`SELECT
		staff_id,
		COUNT(*) AS total_added_count
		FROM customer
		WHERE tenant_id = ?
		AND staff_id IN (?)
		AND deleted_at IS NULL
		GROUP BY staff_id`, tenantID, staffIds).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	rs := make(map[string]int32)
	for _, each := range results {
		rs[each.StaffID] = each.TotalAddedCount
	}

	return rs, nil
}

// GetCustomerCount 获取常客总数
func (c *CustomerStore) GetCustomerCount(ctx context.Context, tenantID string) (int32, error) {
	// 查询所有常客数量
	var count int64
	err := c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// ListCustomerTotalCount 获取常客总数
func (c *CustomerStore) ListCustomerTotalCount(ctx context.Context, tenantID []string) (int32, error) {
	// 查询所有常客数量
	var count int64
	err := c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id IN (?)", tenantID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// ListTodayCustomerCount 获取当日常客总数
func (c *CustomerStore) ListTodayCustomerCount(ctx context.Context, tenantID []string, startTime, endTime time.Time) (int32, error) {
	// 查询所有常客数量
	var count int64
	err := c.GetConnection(ctx).Model(&Customer{}).Where("tenant_id IN (?) and created_at between ? and ?", tenantID, startTime, endTime).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// ListCustomerCount 获取商户常客数量
func (c *CustomerStore) ListCustomerCount(ctx context.Context, tenantID []string, startTime, endTime time.Time) ([]domain.CustomerTenantRankIntf, error) {
	db := c.GetConnection(ctx)
	var rs []struct {
		TenantID      string `gorm:"column:tenant_id"`
		CustomerCount int64  `gorm:"column:customer_count"`
	}
	err := db.Raw(`SELECT tenant_id, COUNT(*) AS customer_count
	FROM customer
	WHERE tenant_id IN (?)
	AND created_at BETWEEN ? AND ?
	GROUP BY tenant_id`, tenantID, startTime, endTime).Scan(&rs).Error
	if err != nil {
		return nil, err
	}

	drs := make([]domain.CustomerTenantRankIntf, len(rs))
	for k, v := range rs {
		drs[k] = &domain.CustomerTenantRank{
			TenantID:      v.TenantID,
			CustomerCount: v.CustomerCount,
		}
	}

	return drs, nil
}
