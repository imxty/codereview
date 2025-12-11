package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/rs/xid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

// UserStore 实现 Domain 层 UserStore 接口
type UserStore struct {
	// 数据库连接
	*dbutils.Connection
}

// 实现 Domain 行为
var _ domain.UserRepository = (*UserStore)(nil)

func NewUserStore(management *dbutils.Connection) *UserStore {
	return &UserStore{
		management,
	}
}

// UpdateOrganizationContactPhone 更新组织联系人电话
func (u *UserStore) UpdateOrganizationContactPhone(ctx context.Context, organization_id, organization_contact_phone string) error {
	db := u.GetConnection(ctx)
	return db.Model(&User{}).Where("organization_id = ?", organization_id).Update("organization_contact_phone", organization_contact_phone).Error
}

// ListPrivilegeGroupsPagination 获取权限组列表
func (u *UserStore) ListPrivilegeGroupsPagination(ctx context.Context, size, offset int32) ([]domain.PrivilegePolicyIntf, int64, error) {
	db := u.GetConnection(ctx)
	var pgs []*PrivilegePolicy
	var total int64

	err := db.Model(&PrivilegePolicy{}).Offset(int(offset)).Limit(int(size)).Find(&pgs).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Model(&PrivilegePolicy{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	dbgs := make([]domain.PrivilegePolicyIntf, len(pgs))
	for k, v := range pgs {
		dbgs[k] = v.ToDomainPrivilegePolicy()
	}
	return dbgs, total, nil
}

// DeleteDataPrivileges 删除数据权限
func (u *UserStore) DeleteDataPrivileges(ctx context.Context, ids []string) error {
	db := u.GetConnection(ctx)
	return db.Exec(`DELETE FROM privilege_policy_data where data_privilege_id in (?)`, ids).Error
}

// BatchGetDataPrivileges 获取数据权限信息
func (u *UserStore) BatchGetDataPrivileges(ctx context.Context, privilege_id string) ([]domain.PrivilegePolicyDataIntf, error) {
	db := u.GetConnection(ctx)
	var datas []PrivilegePolicyData
	err := db.Model(&PrivilegePolicyData{}).Where("privilege_id = ?", privilege_id).Find(&datas).Error
	if err != nil {
		return nil, err
	}

	results := make([]domain.PrivilegePolicyDataIntf, len(datas))
	for k, v := range datas {
		results[k] = v.ToDomainPrivilegePolicyData()
	}

	return results, nil
}

// GetUserByPrivilegeID 通过权限组 ID 获取用户
func (u *UserStore) GetUserByPrivilegeID(ctx context.Context, privilege_id string) (domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	var user User
	err := db.Model(&User{}).Where("privilege_id = ?", privilege_id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user.ToDomainUser(), nil
}

// DeletePrivilegeMenus 删除目录
func (u *UserStore) DeletePrivilegeMenus(ctx context.Context, privilege_id string) error {
	db := u.GetConnection(ctx)
	return db.Exec(`DELETE FROM privilege_menu where privilege_id = ?`, privilege_id).Error
}

// DeletePrivilegePages 删除页面权限
func (u *UserStore) DeletePrivilegePages(ctx context.Context, privilege_id string) error {
	db := u.GetConnection(ctx)
	return db.Exec(`DELETE FROM privilege_policy_page where privilege_id = ?`, privilege_id).Error
}

// UpdatePrivilegeGroup 修改权限组信息
func (u *UserStore) UpdatePrivilegeGroup(ctx context.Context, privilege_id, privilege_name, remark string) error {
	db := u.GetConnection(ctx)
	return db.Model(&PrivilegePolicy{}).Where("privilege_id = ?", privilege_id).Updates(map[string]interface{}{
		"privilege_name": privilege_name,
		"remark":         remark,
	}).Error
}

// DeletePrivilegeGroup 通过 ID 删除权限组
func (u *UserStore) DeletePrivilegeGroup(ctx context.Context, privilege_id string) error {
	db := u.GetConnection(ctx)
	return db.Where("privilege_id = ?", privilege_id).Delete(&PrivilegePolicy{}).Error
}

// GetTreatmentTenant 获取被分配方案的商户
func (u *UserStore) GetTreatmentTenant(ctx context.Context, treatment_id string) (domain.TenantTreatmentIntf, error) {
	db := u.GetConnection(ctx)
	var tt TenantTreatment
	err := db.Model(&TenantTreatment{}).Where("treatment_id = ?", treatment_id).First(&tt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return tt.ToDomainTenantTreatment(), nil
}

// GetMenusByPaths 通过路径获取目录
func (u *UserStore) GetMenusByPaths(ctx context.Context, paths []string) ([]domain.MenuIntf, error) {
	db := u.GetConnection(ctx)
	var ms []Menu
	err := db.Model(&Menu{}).Where("path IN (?)", paths).Find(&ms).Error
	if err != nil {
		return nil, err
	}

	dms := make([]domain.MenuIntf, len(ms))
	for k, v := range ms {
		dms[k] = v.ToDomainMenu()
	}

	return dms, nil
}

// CreateMenus 创建目录
func (u *UserStore) CreateMenus(ctx context.Context, privilege_id string, menu_ids []string) error {
	ms := make([]PrivilegeMenu, len(menu_ids))

	for k, v := range menu_ids {
		ms[k].PrivilegeID = privilege_id
		ms[k].MenuID = v
		ms[k].Rev = 0
		ms[k].CreatedAt = time.Now()
		ms[k].UpdatedAt = time.Now()
	}

	return u.GetConnection(ctx).Model(&PrivilegeMenu{}).Create(&ms).Error
}

// ListAllOrganizationTenantIdsIncludeDeleted 获取组织下所有商户 ID（包括已删除）
func (u *UserStore) ListAllOrganizationTenantIdsIncludeDeleted(ctx context.Context, oids []string) ([]string, error) {
	var tids []string
	err := u.GetConnection(ctx).Raw("select tenant_id from organization_tenant where organization_id IN (?)", oids).Scan(&tids).Error
	if err != nil {
		return nil, err
	}
	return tids, nil
}

// GetOrganizationSubscriptionsMonthly 获取组织月份商户续费信息
func (u *UserStore) GetOrganizationSubscriptionsMonthly(ctx context.Context, organization_id string, start_time, end_time time.Time, size, offset int32) ([]domain.TenantSubscriptionMonthlyIntf, int32, error) {
	db := u.GetConnection(ctx)
	var total int64
	var results []struct {
		Date        string `gorm:"column:period_month"`
		TenantCount int32  `gorm:"column:renewal_count"`
		YearCount   int32  `gorm:"column:total_years"`
	}
	err := db.Raw(`SELECT DATE_FORMAT(created_at, '%Y-%m-01') AS period_month,
	COUNT(DISTINCT tenant_id) AS renewal_count,
	SUM(years) AS total_years
	FROM subscription_period
	WHERE organization_id = ?
	AND created_at BETWEEN ? AND ?
	AND deleted_at IS NULL
	GROUP BY period_month
	ORDER BY period_month
	LIMIT ? OFFSET ?`, organization_id, start_time, end_time, size, offset).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Raw(`SELECT COUNT(*) FROM
	(SELECT DATE_FORMAT(created_at, '%Y-%m-01') AS period_month,
	COUNT(DISTINCT tenant_id) AS renewal_count,
	SUM(years) AS total_years
	FROM subscription_period
	WHERE organization_id = ?
	AND created_at BETWEEN ? AND ?
	AND deleted_at IS NULL
	GROUP BY period_month
	ORDER BY period_month) a`, organization_id, start_time, end_time).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	dres := make([]domain.TenantSubscriptionMonthlyIntf, len(results))
	for k, v := range results {
		dTime, err := time.Parse("2006-01-02", v.Date)
		if err != nil {
			return nil, 0, err
		}
		dres[k] = &domain.TenantSubscriptionMonthly{
			Date:        dTime,
			TenantCount: v.TenantCount,
			YearCount:   v.YearCount,
		}
	}

	return dres, int32(total), nil
}

// RecoverStaff 恢复已删除的用户
func (u *UserStore) RecoverStaff(ctx context.Context, staffId, tenantId string, rev int32) error {
	// 修改tenant_id
	return u.GetConnection(ctx).Model(&User{}).Where("role_type = ? and user_id = ? and is_activated = 0 and rev = ?", domain.RoleTypeStaff, staffId, rev).Updates(map[string]interface{}{
		"rev":          rev + 1,
		"is_activated": 1,
		"tenant_id":    tenantId,
	}).Error
}

// GetDataPrivilegePagination 分页获取数据权限
func (u *UserStore) GetDataPrivilegePagination(ctx context.Context, privilege_id string, size, offset int32) ([]string, int32, error) {
	db := u.GetConnection(ctx)
	var ds []string
	var count int64
	err := db.Raw(`SELECT organization_id
	FROM privilege_policy_data
	WHERE privilege_id = ?
	AND deleted_at IS NULL
	LIMIT ?
	OFFSET ?`, privilege_id, size, offset).Scan(&ds).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Raw(`SELECT COUNT(*)
	FROM privilege_policy_data
	WHERE privilege_id = ?
	AND deleted_at IS NULL`, privilege_id).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return ds, int32(count), nil
}

// GetTenantTreatments 获取商户方案
func (u *UserStore) GetTenantTreatments(ctx context.Context, tids []string) ([]domain.TenantTreatmentIntf, error) {
	db := u.GetConnection(ctx)
	var tts []TenantTreatment
	err := db.Model(&TenantTreatment{}).Where("tenant_id IN (?)", tids).Find(&tts).Error
	if err != nil {
		return nil, err
	}

	results := make([]domain.TenantTreatmentIntf, len(tts))
	for k, v := range tts {
		results[k] = v.ToDomainTenantTreatment()
	}

	return results, nil
}

// UpdateSafePhone 更新User中tenant安全手机号
func (u *UserStore) UpdateTenantUserSafePhone(ctx context.Context, tid string, safePhone string, rev int32) error {
	// 更新手机号
	return u.GetConnection(ctx).Model(&User{}).Where("tenant_id = ? and role_type = ? and rev = ?", tid, domain.RoleTypeTenant, rev).Updates(map[string]interface{}{
		"phone": safePhone,
		"rev":   rev + 1,
	}).Error
}

// UpdateTenantSafePhone 更新tenantEntity安全手机号
func (u *UserStore) UpdateTenantSafePhone(ctx context.Context, tid string, safePhone string, rev int32) error {
	// 更新手机号
	return u.GetConnection(ctx).Model(&TenantEntity{}).Where("tenant_id = ? and rev = ?", tid, rev).Updates(map[string]interface{}{
		"safe_phone": safePhone,
		"rev":        rev + 1,
	}).Error
}

// ResetAppStaffPassword 重置 APP 登录密码
func (u *UserStore) ResetAppStaffPassword(ctx context.Context, user_id string, new_password string) error {
	db := u.GetConnection(ctx)
	return db.Model(&User{}).Where("user_id = ? AND role_type = ?", user_id, domain.RoleTypeStaff).Update("hashed_password", new_password).Error
}

// CancelTreatmentToTenants 取消方案分配
func (u *UserStore) CancelTreatmentToTenants(ctx context.Context, tenant_ids []string) error {
	db := u.GetConnection(ctx)

	return db.Where("tenant_id IN (?)", tenant_ids).Delete(&TenantTreatment{}).Error
}

// SubmitFeedback 提交反馈
func (u *UserStore) SubmitFeedback(ctx context.Context, feedback domain.FeedbackIntf) error {
	db := u.GetConnection(ctx)
	fd := new(Feedback)
	fd.FromDomainFeedback(feedback)
	return db.Model(&Feedback{}).Create(&fd).Error
}

// UpdateStaffInfo 更新员工信息
func (u *UserStore) UpdateStaffInfo(ctx context.Context, staff_id, staff_name, staff_phone, privilege_id, remark string) error {
	db := u.GetConnection(ctx)
	return db.Model(&User{}).Where("user_id = ? AND role_type >= ?", staff_id, domain.RoleTypeJmAdmin).Updates(map[string]interface{}{
		"username":     staff_name,
		"nickname":     staff_name,
		"phone":        staff_phone,
		"privilege_id": privilege_id,
		"remark":       remark,
	}).Error
}

// BatchAddTenantTimelines 更新所有已过期的时间线
func (u *UserStore) BatchAddTenantTimelines(ctx context.Context, tids []string, end_time time.Time, years int32) error {
	db := u.GetConnection(ctx)
	return db.Model(&SubscriptionTimeline{}).Where("tenant_id in (?)", tids).Updates(map[string]interface{}{
		"end_time": end_time,
		"years":    gorm.Expr("years + ?", years),
	}).Error
}

// BatchUpdateTenantTimelines 增加所有未过期的时间线
func (u *UserStore) BatchUpdateTenantTimelines(ctx context.Context, tids []string, add_time int32, years int32) error {
	db := u.GetConnection(ctx)
	return db.Model(&SubscriptionTimeline{}).Where("tenant_id in (?)", tids).Updates(map[string]interface{}{
		"end_time": gorm.Expr("DATE_ADD(end_time, INTERVAL ? YEAR)", add_time),
		"years":    gorm.Expr("years + ?", years),
	}).Error
}

// BatchCreateTenantSubscriptionPeriod 批量创建订阅周期
func (u *UserStore) BatchCreateTenantSubscriptionPeriod(ctx context.Context, periods []domain.SubscriptionPeriodIntf) error {
	db := u.GetConnection(ctx)

	tls := make([]*SubscriptionPeriod, len(periods))
	for k, v := range periods {
		us := new(SubscriptionPeriod)
		us.FromDomainSubscriptionPeriod(v)
		tls[k] = us
	}

	return db.Model(&SubscriptionPeriod{}).Create(&tls).Error
}

// BatchGetLatestTenantSubscriptionTimelinesDetails 分页获取订阅时间线
func (u *UserStore) BatchGetLatestTenantSubscriptionTimelinesDetails(ctx context.Context, tenantIds []string) ([]domain.SubscriptionTimelineIntf, int64, error) {
	db := u.GetConnection(ctx)
	var count int64
	ts := []SubscriptionTimeline{}
	err := db.Model(&SubscriptionTimeline{}).Where("tenant_id in (?)", tenantIds).Find(&ts).Error
	if err != nil {
		return nil, count, err
	}

	err = db.Model(&SubscriptionTimeline{}).Where("tenant_id in (?)", tenantIds).Count(&count).Error
	if err != nil {
		return nil, count, err
	}

	dts := make([]domain.SubscriptionTimelineIntf, len(ts))
	for k, v := range ts {
		dts[k] = v.ToDomainSubscriptionTimeline()
	}

	return dts, count, nil
}

// BatchGetOrganizationTenants 批量获取组织商户关系
func (u *UserStore) BatchGetOrganizationTenants(ctx context.Context, tids []string) ([]domain.OrganizationTenantIntf, error) {
	db := u.GetConnection(ctx)
	ots := []OrganizationTenant{}
	err := db.Model(&OrganizationTenant{}).Where("tenant_id in (?)", tids).Find(&ots).Error
	if err != nil {
		return nil, err
	}

	dots := make([]domain.OrganizationTenantIntf, len(ots))
	for k, v := range ots {
		dots[k] = v.ToDomainOrganizationTenant()
	}

	return dots, nil
}

// BatchGetTenantsSubscriptionPeriodsDetails 查询商户订阅详情
func (u *UserStore) BatchGetTenantsSubscriptionPeriodsDetails(ctx context.Context, tenant_ids []string, contact_name, contact_phone string, start_time time.Time, end_time time.Time, offset int32, size int32) ([]domain.SubscriptionPeriodIntf, int64, int32, error) {
	db := u.GetConnection(ctx)
	ps := []SubscriptionPeriod{}
	var count int64
	var total_years int32

	// 基础查询条件
	baseQuery := db.Model(&SubscriptionPeriod{}).Where("tenant_id in (?)", tenant_ids).Where("created_at BETWEEN ? AND ?", start_time, end_time)

	if contact_name != "" {
		baseQuery = baseQuery.Where("contact_name LIKE ?", "%"+contact_name+"%")
	}

	if contact_phone != "" {
		baseQuery = baseQuery.Where("contact_phone = ?", contact_phone)
	}

	// 统计总数（不受分页影响）
	err := baseQuery.Session(&gorm.Session{}).Count(&count).Error
	if err != nil {
		return nil, count, total_years, err
	}

	// 统计总订阅年限（不受分页影响）
	err = baseQuery.Session(&gorm.Session{}).Select("COALESCE(SUM(years), 0)").Scan(&total_years).Error
	if err != nil {
		return nil, count, total_years, err
	}

	// 查询分页数据
	err = baseQuery.Offset(int(offset)).Limit(int(size)).Find(&ps).Error
	if err != nil {
		return nil, count, total_years, err
	}

	dps := make([]domain.SubscriptionPeriodIntf, len(ps))
	for k, v := range ps {
		dps[k] = v.ToDomainSubscriptionPeriod()
	}

	return dps, count, total_years, nil
}

// BatchGetTenantsSubscriptionPeriods 批量获取商户订阅周期
func (u *UserStore) BatchGetTenantsSubscriptionPeriods(ctx context.Context, tenant_ids []string, start_time time.Time, end_time time.Time, offset int32, size int32) ([]domain.SubscriptionPeriodIntf, int64, error) {
	db := u.GetConnection(ctx)
	ps := []SubscriptionPeriod{}
	var count int64
	err := db.Model(&SubscriptionPeriod{}).Where("tenant_id in (?)", tenant_ids).Where("created_at BETWEEN ? AND ?", start_time, end_time).Offset(int(offset)).Limit(int(size)).Find(&ps).Error
	if err != nil {
		return nil, count, err
	}

	err = db.Model(&SubscriptionPeriod{}).Where("tenant_id in (?)", tenant_ids).Where("created_at BETWEEN ? AND ?", start_time, end_time).Count(&count).Error
	if err != nil {
		return nil, count, err
	}

	dps := make([]domain.SubscriptionPeriodIntf, len(ps))
	for k, v := range ps {
		dps[k] = v.ToDomainSubscriptionPeriod()
	}

	return dps, count, nil
}

// GetMenus 根据权限组 ID 获取菜单
func (u *UserStore) GetMenus(ctx context.Context, privilege_id string) ([]domain.MenuIntf, error) {
	db := u.GetConnection(ctx)
	ms := []Menu{}
	err := db.Raw(`SELECT m.*
	FROM menu m
	JOIN privilege_menu pm ON m.menu_id = pm.menu_id
	WHERE pm.privilege_id = ?
	AND m.deleted_at IS NULL
	AND pm.deleted_at IS NULL`, privilege_id).Scan(&ms).Error
	if err != nil {
		return nil, err
	}

	dms := make([]domain.MenuIntf, len(ms))
	for k, v := range ms {
		dms[k] = v.ToDomainMenu()
	}

	return dms, nil
}

// GetPagePrivilegeByID 通过权限组 ID 获取页面权限信息
func (u *UserStore) GetPagePrivilegeByID(ctx context.Context, privilege_id string) ([]domain.PrivilegePolicyPageIntf, error) {
	db := u.GetConnection(ctx)
	ps := []PrivilegePolicyPage{}
	err := db.Model(&PrivilegePolicyPage{}).Where("privilege_id = ?", privilege_id).Find(&ps).Error
	if err != nil {
		return nil, err
	}

	dps := make([]domain.PrivilegePolicyPageIntf, len(ps))
	for k, v := range ps {
		dps[k] = v.ToDomainPrivilegePolicyPage()
	}

	return dps, nil
}

// GetPrivilegeGroupsByName 根据名称获取权限组信息
func (u *UserStore) GetPrivilegeGroupsByName(ctx context.Context, privilege_name string) ([]domain.PrivilegePolicyIntf, error) {
	db := u.GetConnection(ctx)
	ps := []PrivilegePolicy{}
	err := db.Model(&PrivilegePolicy{}).Where("privilege_name like ?", "%"+privilege_name+"%").Find(&ps).Error
	if err != nil {
		return nil, err
	}

	dps := make([]domain.PrivilegePolicyIntf, len(ps))
	for k, v := range ps {
		dps[k] = v.ToDomainPrivilegePolicy()
	}

	return dps, nil
}

// ListTreatmentsByTenantIDs 根据商户 ID 获取方案
func (u *UserStore) ListTreatmentsByTenantIDs(ctx context.Context, tids []string) ([]domain.TenantTreatmentIntf, error) {
	db := u.GetConnection(ctx)
	ts := []TenantTreatment{}
	err := db.Model(&TenantTreatment{}).Where("tenant_id in (?)", tids).Find(&ts).Error
	if err != nil {
		return nil, err
	}

	dts := make([]domain.TenantTreatmentIntf, len(ts))
	for k, v := range ts {
		dts[k] = v.ToDomainTenantTreatment()
	}

	return dts, nil
}

// SearchTenantByNameAndOrganizationID 根据组织 ID 和商户名模糊搜索商户 ID
func (u *UserStore) SearchTenantByNameAndOrganizationID(ctx context.Context, organization_id string, tenant_name string) ([]string, error) {
	db := u.GetConnection(ctx)
	ts := []string{}
	err := db.Raw(`SELECT ot.tenant_id
	FROM organization_tenant ot
	JOIN tenant_entity te ON ot.tenant_id = te.tenant_id
	WHERE ot.organization_id = ?
	AND te.store_name LIKE ?
	AND ot.deleted_at IS NULL
	AND te.deleted_at IS NULL`, organization_id, "%"+tenant_name+"%").Scan(&ts).Error
	if err != nil {
		return nil, err
	}

	return ts, nil
}

// SearchTenantsByName 通过商户名模糊搜索商户
func (u *UserStore) SearchTenantsByName(ctx context.Context, tids []string, tenant_name string, contact_name string, contact_phone string, status, offset, size int32) ([]domain.TenantEntityIntf, int32, error) {
	db := u.GetConnection(ctx)
	ts := []TenantEntity{}
	var total int64
	var sql string
	// 正常使用
	if status == int32(userv1.TenantStatus_TENANT_STATUS_USING) {
		sql = `SELECT te.*
		FROM tenant_entity te
		JOIN subscription_timeline st ON te.tenant_id = st.tenant_id
		WHERE st.end_time > NOW()
		AND te.tenant_id IN (?)
		AND te.store_name LIKE ?
		AND te.contact_name like ?
		AND te.deleted_at IS NULL
		AND st.deleted_at IS NULL`
	}
	// 待续期
	if status == int32(userv1.TenantStatus_TENANT_STATUS_PENDING) {
		sql = `SELECT te.*
		FROM tenant_entity te
		JOIN subscription_timeline st ON te.tenant_id = st.tenant_id
		WHERE st.end_time <= NOW()
		AND te.tenant_id IN (?)
		AND te.store_name LIKE ?
		AND te.contact_name like ?
		AND te.deleted_at IS NULL
		AND st.deleted_at IS NULL`
	}
	// 全部
	if status == int32(userv1.TenantStatus_TENANT_STATUS_UNSET) {
		sql = `SELECT *
		FROM tenant_entity
		WHERE tenant_id IN (?)
		AND store_name like ?
		AND contact_name like ?
		AND deleted_at IS NULL`
	}
	if contact_phone != "" {
		sql += " AND contact_phone = ?"
	}
	total_sql := sql
	sql += " LIMIT ? OFFSET ?"

	var args []interface{}
	args = append(args, tids, "%"+tenant_name+"%", "%"+contact_name+"%")
	if contact_phone != "" {
		args = append(args, contact_phone)
	}

	err := db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM (%s) e", total_sql), args...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	args = append(args, size, offset)

	err = db.Raw(sql, args...).Scan(&ts).Error
	if err != nil {
		return nil, 0, err
	}

	dts := make([]domain.TenantEntityIntf, len(ts))
	for k, v := range ts {
		dts[k] = v.ToDomainTenantEntity()
	}

	return dts, int32(total), nil
}

// SearchTenantsByNameAndTime 通过商户名和时间模糊搜索商户
func (u *UserStore) SearchTenantsByNameAndTime(ctx context.Context, tids []string, tenant_name string, contact_name string, contact_phone string, status, offset, size int32, start_time, end_time *timestamppb.Timestamp) ([]domain.TenantEntityIntf, int32, error) {
	db := u.GetConnection(ctx)
	ts := []TenantEntity{}
	var total int64
	var sql string
	// 正常使用
	if status == int32(userv1.TenantStatus_TENANT_STATUS_USING) {
		sql = `SELECT te.*
		FROM tenant_entity te
		JOIN subscription_timeline st ON te.tenant_id = st.tenant_id
		WHERE st.end_time > NOW()
		AND te.tenant_id IN (?)
		AND te.store_name LIKE ?
		AND te.contact_name like ?
		AND te.deleted_at IS NULL
		AND st.deleted_at IS NULL`
	}
	// 待续期
	if status == int32(userv1.TenantStatus_TENANT_STATUS_PENDING) {
		sql = `SELECT te.*
		FROM tenant_entity te
		LEFT JOIN subscription_timeline st ON te.tenant_id = st.tenant_id
		WHERE (st.end_time <= NOW() OR st.subscription_timeline_id IS NULL)
		AND te.tenant_id IN (?)
		AND te.safe_phone != ""
		AND te.store_name LIKE ?
		AND te.contact_name like ?
		AND te.deleted_at IS NULL
		AND (st.deleted_at IS NULL OR st.subscription_timeline_id IS NULL)`
	}
	// 全部
	if status == int32(userv1.TenantStatus_TENANT_STATUS_UNSET) {
		sql = `SELECT te.*
		FROM tenant_entity te
		LEFT JOIN subscription_timeline st ON te.tenant_id = st.tenant_id
		WHERE te.tenant_id IN (?)
		AND te.store_name like ?
		AND te.contact_name like ?
		AND te.deleted_at IS NULL
		AND (st.deleted_at IS NULL OR st.subscription_timeline_id IS NULL)`
	}
	// 判断过期时间
	if start_time != nil {
		sql += " AND st.end_time >= ?"
	}
	if end_time != nil {
		sql += " AND st.end_time <= ?"
	}
	if contact_phone != "" {
		sql += " AND contact_phone = ?"
	}
	total_sql := sql
	sql += " LIMIT ? OFFSET ?"

	var args []interface{}
	args = append(args, tids, "%"+tenant_name+"%", "%"+contact_name+"%")
	if contact_phone != "" {
		args = append(args, contact_phone)
	}

	if start_time != nil {
		args = append(args, start_time.AsTime())
	}

	if end_time != nil {
		args = append(args, end_time.AsTime())
	}

	err := db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM (%s) e", total_sql), args...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	args = append(args, size, offset)

	err = db.Raw(sql, args...).Scan(&ts).Error
	if err != nil {
		return nil, 0, err
	}

	dts := make([]domain.TenantEntityIntf, len(ts))
	for k, v := range ts {
		dts[k] = v.ToDomainTenantEntity()
	}

	return dts, int32(total), nil
}

// UpdateTenantRevisionFailReason 更新审核失败信息
func (u *UserStore) UpdateTenantRevisionFailReason(ctx context.Context, revision_id string, fail_reason string) error {
	db := u.GetConnection(ctx)
	return db.Model(&TenantEntityRevision{}).Where("tenant_entity_revision_id = ?", revision_id).Update("fail_reason", fail_reason).Error
}

// ResetPassword 重置密码
func (u *UserStore) ResetPassword(ctx context.Context, user_id string, newHashedPassword string) error {
	db := u.GetConnection(ctx)
	return db.Model(&User{}).Where("user_id = ?", user_id).Updates(map[string]interface{}{
		"hashed_password": newHashedPassword,
	}).Error
}

// GetTenantByName 根据名称获取商户信息
func (u *UserStore) GetTenantByName(ctx context.Context, name string) (domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	us := &User{}
	err := db.Model(&User{}).Where("nickname = ?", name).First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return us.ToDomainUser(), nil
}

// BatchGetOrganizationNamesByIDs 通过 ID 批量获取组织信息
func (u *UserStore) BatchGetOrganizationNamesByIDs(ctx context.Context, organization_ids []string) ([]domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	var us []*User
	err := db.Model(&User{}).Where("role_type = ? and organization_id in (?)", domain.RoleTypeOrganization, organization_ids).Find(&us).Error
	if err != nil {
		return nil, err
	}

	// 转化
	dus := make([]domain.UserIntf, len(us))
	for k, v := range us {
		dus[k] = v.ToDomainUser()
	}
	return dus, nil
}

// GetDataPrivilegeByID 获取权限组数据权限
func (u *UserStore) GetDataPrivilegeByID(ctx context.Context, privilege_id string) ([]string, error) {
	db := u.GetConnection(ctx)
	var oids []string
	err := db.Raw(`SELECT organization_id
	FROM privilege_policy_data 
	WHERE privilege_id = ?
	AND deleted_at IS NULL`, privilege_id).Scan(&oids).Error
	if err != nil {
		return nil, err
	}
	return oids, nil
}

// GetDataPrivilegeByOrganizationID 通过用户 ID 和组织 ID 获取权限
func (u *UserStore) GetDataPrivilegeByOrganizationID(ctx context.Context, privilege_id string, organization_id string) (domain.PrivilegePolicyDataIntf, error) {
	db := u.GetConnection(ctx)
	var ps []*PrivilegePolicyData
	err := db.Model(&PrivilegePolicyData{}).Where("privilege_id = ? and organization_id = ?", privilege_id, organization_id).Find(&ps).Error
	if err != nil {
		return nil, err
	}
	// 查询到了数据
	if len(ps) > 0 {
		return ps[0].ToDomainPrivilegePolicyData(), nil
	}
	// 没有数据
	return nil, nil
}

// SearchStaffsByNameAndStatus 根据名字和状态搜索员工数据
func (u *UserStore) SearchStaffsByNameAndStatus(ctx context.Context, tenant_id string, staff_name string, is_activated *wrapperspb.BoolValue) ([]domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	nameQuery := "%" + staff_name + "%"
	var us []*User
	rawSql := "tenant_id = ? and nickname like ? AND deleted_at IS NULL"
	if is_activated != nil {
		if is_activated.GetValue() {
			rawSql += " and is_activated = 1"
		}
		if !is_activated.GetValue() {
			rawSql += " and is_activated = 0"
		}
	}
	err := db.Model(&User{}).Where(rawSql, tenant_id, nameQuery).Find(&us).Error
	if err != nil {
		return nil, err
	}

	dus := make([]domain.UserIntf, len(us))
	for k, v := range us {
		dus[k] = v.ToDomainUser()
	}
	return dus, nil
}

func (u *UserStore) UpdateTreatmentToTenant(ctx context.Context, organization_id string, tenant_ids []string, treatment_id string) error {
	db := u.GetConnection(ctx)
	return db.Model(TenantTreatment{}).Where("tenant_id IN (?)", tenant_ids).Update("treatment_id", treatment_id).Error
}

// SubmitTreatmentToTenant 分配方案
func (u *UserStore) SubmitTreatmentToTenant(ctx context.Context, organization_id string, tenant_ids []string, treatment_id string) error {
	db := u.GetConnection(ctx)
	temp := make([]TenantTreatment, len(tenant_ids))
	for k, v := range tenant_ids {
		temp[k].TenantTreatmentID = xid.New().String()
		temp[k].TenantID = v
		temp[k].TreatmentID = treatment_id
	}
	return db.Model(&TenantTreatment{}).Create(&temp).Error
}

// UpdateOrganizationPhone 更新组织手机号
func (u *UserStore) UpdateOrganizationPhone(ctx context.Context, organization_id string, new_phone string) error {
	db := u.GetConnection(ctx)
	return db.Model(&User{}).Where("organization_id = ? and role_type = ?", organization_id, domain.RoleTypeOrganization).Updates(map[string]interface{}{
		"phone": new_phone,
	}).Error
}

// UpdateTenantConstitutionSwitch 修改体质辩证开关状态
func (u *UserStore) UpdateTenantConstitutionSwitch(ctx context.Context, tenant_id string, constitution_switch_status bool) error {
	db := u.GetConnection(ctx)
	return db.Model(&TenantEntity{}).Where("tenant_id = ?", tenant_id).Updates(map[string]interface{}{
		"constitution_status": constitution_switch_status,
	}).Error
}

// UpdateTenantReview 修改复查天数
func (u *UserStore) UpdateTenantReview(ctx context.Context, tenant_id string, days int32) error {
	db := u.GetConnection(ctx)
	return db.Model(&TenantEntity{}).Where("tenant_id = ?", tenant_id).Updates(map[string]interface{}{
		"overdue": days,
	}).Error
}

// GetExpiringTenantIds 查询即将过期的商户 ID
func (u *UserStore) GetExpiringTenantIds(ctx context.Context, organization_id string) ([]string, error) {
	db := u.GetConnection(ctx)
	var tids []string
	err := db.Raw(`SELECT st.tenant_id
	FROM subscription_timeline st
	JOIN organization_tenant ot ON st.tenant_id = ot.tenant_id
	WHERE ot.organization_id = ?
	AND st.deleted_at IS NULL
	AND ot.deleted_at IS NULL
	AND st.end_time BETWEEN DATE_FORMAT(NOW(), '%Y-%m-01 00:00:00')
		AND LAST_DAY(NOW()) + INTERVAL 1 DAY - INTERVAL 1 SECOND
	`, organization_id).Scan(&tids).Error
	if err != nil {
		return nil, err
	}
	return tids, nil
}

// GetTodayAddedTenantIds 获取今日新增商户 ID
func (u *UserStore) GetTodayAddedTenantIds(ctx context.Context, organization_id string) ([]string, error) {
	db := u.GetConnection(ctx)
	var tids []string
	err := db.Raw(`SELECT tenant_id
	FROM organization_tenant
	WHERE organization_id = ?
	AND deleted_at IS NULL
	AND DATE(created_at) = CURDATE()`, organization_id).Scan(&tids).Error
	if err != nil {
		return nil, err
	}
	return tids, nil
}

// GetUsingTenantIds 获取正在使用中的商户
func (u *UserStore) GetUsingTenantIds(ctx context.Context, organization_id string) ([]string, error) {
	db := u.GetConnection(ctx)
	var tids []string
	err := db.Raw(`SELECT ot.tenant_id
	FROM
    	organization_tenant ot
	JOIN
		subscription_timeline st
	ON
		ot.tenant_id = st.tenant_id
	WHERE
		ot.deleted_at IS NULL
		AND st.deleted_at IS NULL
		AND ot.organization_id = ? 
		AND ot.is_activated = 1
		AND st.end_time > NOW();`, organization_id).Scan(&tids).Error
	if err != nil {
		return nil, err
	}
	return tids, nil
}

// ListPrivilegeGroups 获取所有权限组
func (u *UserStore) ListPrivilegeGroups(ctx context.Context) ([]domain.PrivilegePolicyIntf, error) {
	db := u.GetConnection(ctx)
	var pgs []*PrivilegePolicy
	err := db.Model(&PrivilegePolicy{}).Find(&pgs).Error
	if err != nil {
		return nil, err
	}

	dbgs := make([]domain.PrivilegePolicyIntf, len(pgs))
	for k, v := range pgs {
		dbgs[k] = v.ToDomainPrivilegePolicy()
	}
	return dbgs, nil
}

// SearchOrganizationByName 通过组织名查询组织信息
func (u *UserStore) SearchOrganizationByName(ctx context.Context, organization_name string) ([]domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	var us []*User
	query := "%" + organization_name + "%"
	err := db.Model(&User{}).Where("role_type = ? AND nickname like ?", domain.RoleTypeOrganization, query).Find(&us).Error
	if err != nil {
		return nil, err
	}

	dus := make([]domain.UserIntf, len(us))
	for k, v := range us {
		dus[k] = v.ToDomainUser()
	}
	return dus, nil
}

// SearchSystemUser 查询系统用户
func (u *UserStore) SearchSystemUser(ctx context.Context, staff_name string, contact_phone string, privilege_group_ids []string, is_activated *wrapperspb.BoolValue, offset, size int32) ([]domain.UserIntf, int32, error) {
	db := u.GetConnection(ctx)
	nameQuery := "%" + staff_name + "%"
	phoneQuery := "%" + contact_phone + "%"
	querySQL := "username like ? and phone like ? and privilege_id in ?"
	if is_activated != nil {
		if is_activated.GetValue() {
			querySQL += " and is_activated = 1"
		}
		if !is_activated.GetValue() {
			querySQL += " and is_activated = 0"
		}
	}

	var count int64
	var us []*User
	err := db.Model(&User{}).Where(querySQL, nameQuery, phoneQuery, privilege_group_ids).Offset(int(offset)).Limit(int(size)).Find(&us).Error
	if err != nil {
		return nil, int32(count), err
	}
	err = db.Model(&User{}).Where(querySQL, nameQuery, phoneQuery, privilege_group_ids).Count(&count).Error
	if err != nil {
		return nil, int32(count), err
	}

	dus := make([]domain.UserIntf, len(us))
	for k, v := range us {
		dus[k] = v.ToDomainUser()
	}
	return dus, int32(count), err
}

// SearchTenantIdsByName 通过商户名模糊查询商户 ID
func (u *UserStore) SearchTenantIdsByName(ctx context.Context, tenant_name string) ([]string, error) {
	db := u.GetConnection(ctx)

	query := "%" + tenant_name + "%"
	var ids []string
	err := db.Raw(`SELECT tenant_id
	FROM tenant_entity
	WHERE store_name like ?
	AND deleted_at IS NULL`, query).Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// BatchGetTenantNamesByIds 通过商户 ID 批量获取商户
func (u *UserStore) BatchGetTenantNamesByIds(ctx context.Context, tenant_ids []string) ([]domain.TenantEntityIntf, error) {
	db := u.GetConnection(ctx)
	var tenants []*TenantEntity
	err := db.Model(&TenantEntity{}).Where("tenant_id in (?)", tenant_ids).Find(&tenants).Error
	if err != nil {
		return nil, err
	}

	results := make([]domain.TenantEntityIntf, len(tenants))
	for k, v := range tenants {
		results[k] = v.ToDomainTenantEntity()
	}

	return results, nil
}

// CreateDataPrivileges 创建数据权限
func (u *UserStore) CreateDataPrivileges(ctx context.Context, privilegeId string, organizationIds []string) error {
	datas := make([]PrivilegePolicyData, len(organizationIds))
	for k, v := range organizationIds {
		datas[k].DataPrivilegeID = xid.New().String()
		datas[k].PrivilegeID = privilegeId
		datas[k].OrganizationID = v
		datas[k].Rev = 0
	}
	db := u.GetConnection(ctx)
	return db.Model(&PrivilegePolicyData{}).Create(&datas).Error
}

// CreatePagePrivileges 创建页面权限
func (u *UserStore) CreatePagePrivileges(ctx context.Context, privilege_id string, preset_ids []string) error {
	db := u.GetConnection(ctx)
	pages := make([]PrivilegePolicyPage, len(preset_ids))
	for k, v := range preset_ids {
		pages[k].PrivilegeID = privilege_id
		pages[k].PresetID = v
	}
	return db.Model(&PrivilegePolicyPage{}).Create(&pages).Error
}

// CreatePrivilegeGroup 创建权限组
func (u *UserStore) CreatePrivilegeGroup(ctx context.Context, privilegeId string, privilegeName string, remark string) error {
	db := u.GetConnection(ctx)
	group := PrivilegePolicy{
		PrivilegeID:   privilegeId,
		PrivilegeName: privilegeName,
		Remark:        remark,
		Rev:           0,
	}
	return db.Model(&PrivilegePolicy{}).Create(&group).Error
}

// CreateSystemUser 创建系统用户
func (u *UserStore) CreateSystemUser(ctx context.Context, systemUser domain.UserIntf) error {
	db := u.GetConnection(ctx)
	user := new(User)
	user.FromDomainUser(systemUser)
	return db.Model(&User{}).Create(&user).Error
}

// DeleteSystemUser 删除系统用户（员工离职）
func (u *UserStore) DeleteSystemUser(ctx context.Context, user_id string) error {
	return u.GetConnection(ctx).Model(&User{}).Where("user_id = ?", user_id).Updates(map[string]interface{}{
		"is_activated": 0,
	}).Error
}

// DeleteTenantUser 删除商户用户
func (u *UserStore) DeleteTenantUser(ctx context.Context, tenant_id string) error {
	db := u.GetConnection(ctx)
	return db.Where("tenant_id = ?", tenant_id).Delete(&User{}).Error
}

// GetPrivilegeGroupById 根据 ID 获取权限组
func (u *UserStore) GetPrivilegeGroupById(ctx context.Context, privilege_id string) (domain.PrivilegePolicyIntf, error) {
	db := u.GetConnection(ctx)
	var group PrivilegePolicy
	err := db.Model(&PrivilegePolicy{}).Where("privilege_id = ?", privilege_id).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return group.ToDomainPrivilegePolicy(), nil
}

// BatchCreateTenantSubscriptionTimeline 批量创建商户订阅时间线
func (u *UserStore) BatchCreateTenantSubscriptionTimeline(ctx context.Context, timelines []domain.SubscriptionTimelineIntf) error {
	db := u.GetConnection(ctx)
	tls := make([]*SubscriptionTimeline, len(timelines))
	for k, v := range timelines {
		us := new(SubscriptionTimeline)
		us.FromDomainSubscriptionTimeline(v)
		tls[k] = us
	}
	return db.Model(&SubscriptionTimeline{}).Create(&tls).Error
}

// BatchGetLatestTenantSubscriptionTimelines 批量获取最新的订阅信息
func (u *UserStore) BatchGetLatestTenantSubscriptionTimelines(ctx context.Context, tenantIds []string) ([]domain.SubscriptionTimelineIntf, error) {
	db := u.GetConnection(ctx)
	var us []SubscriptionTimeline
	err := db.Model(&SubscriptionTimeline{}).Where("tenant_id in (?)", tenantIds).Find(&us).Error
	if err != nil {
		return nil, err
	}
	dus := make([]domain.SubscriptionTimelineIntf, len(us))
	for k, v := range us {
		dus[k] = v.ToDomainSubscriptionTimeline()
	}
	return dus, nil
}

// BatchGetStaffs 批量获取员工
func (u *UserStore) BatchGetStaffs(ctx context.Context, staffIds []string) ([]domain.UserIntf, error) {
	var staffs []User
	err := u.GetConnection(ctx).Model(&User{}).Where("user_id IN (?) AND role_type = ?", staffIds, domain.RoleTypeStaff).Find(&staffs).Error
	if err != nil {
		return nil, err
	}
	dSts := make([]domain.UserIntf, len(staffs))
	for k, v := range staffs {
		dSts[k] = v.ToDomainUser()
	}
	return dSts, nil
}

// ChangeOrganizationTenantReviewStatus 改变组织商户审核状态
func (u *UserStore) ChangeOrganizationTenantReviewStatus(ctx context.Context, tenantId string, status int32, rev int32) error {
	db := u.GetConnection(ctx)
	return db.Model(&OrganizationTenant{}).Where("tenant_id = ? and rev = ?", tenantId, rev).Updates(map[string]interface{}{
		"rev":           rev + 1,
		"review_status": status,
	}).Error
}

// CheckOrganizationIDExist 检测组织 ID 是否存在
func (u *UserStore) CheckOrganizationIDExist(ctx context.Context, id string) (bool, error) {
	db := u.GetConnection(ctx)
	o := &User{}
	err := db.Model(&User{}).Where("organization_id = ?", id).First(&o).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateOrganizationTenant 创建组织商户信息
func (u *UserStore) CreateOrganizationTenant(ctx context.Context, ot domain.OrganizationTenantIntf) error {
	db := u.GetConnection(ctx)
	dt := new(OrganizationTenant)
	dt.FromDomainOrganizationTenant(ot)
	return db.Model(&OrganizationTenant{}).Create(&dt).Error
}

// CreateStaff 创建员工
func (u *UserStore) CreateStaff(ctx context.Context, staff domain.UserIntf) error {
	var user User
	user.FromDomainUser(staff)
	return u.GetConnection(ctx).Model(&User{}).Create(&user).Error
}

// CreateSubscriptionTimeline 创建订阅
func (u *UserStore) CreateSubscriptionTimeline(ctx context.Context, timelines []domain.SubscriptionTimelineIntf) error {

	tls := make([]*SubscriptionTimeline, len(timelines))
	for k, v := range timelines {
		us := new(SubscriptionTimeline)
		us.FromDomainSubscriptionTimeline(v)
		tls[k] = us
	}
	return u.GetConnection(ctx).Model(&SubscriptionTimeline{}).Create(&tls).Error
}

// CreateTenantEntity 创建商户
func (u *UserStore) CreateTenantEntity(ctx context.Context, te domain.TenantEntityIntf) error {
	var pte TenantEntity
	pte.FromDomainTenantEntity(te)
	return u.GetConnection(ctx).Model(&TenantEntity{}).Create(&pte).Error
}

// CreateTenantEntityRevision 创建商户副本
func (u *UserStore) CreateTenantEntityRevision(ctx context.Context, revision domain.TenantEntityRevisionIntf) error {
	var vision TenantEntityRevision
	vision.FromDomainTenantEntityRevision(revision)
	return u.GetConnection(ctx).Model(&TenantEntityRevision{}).Create(&vision).Error
}

// CreateTenantSubscriptionTimeline 创建商户时间线
func (u *UserStore) CreateTenantSubscriptionTimeline(ctx context.Context, timeline domain.SubscriptionTimelineIntf) error {
	var ts SubscriptionTimeline
	ts.FromDomainSubscriptionTimeline(timeline)
	return u.GetConnection(ctx).Model(&SubscriptionTimeline{}).Create(ts).Error
}

// CreateUser 创建用户
func (u *UserStore) CreateUser(ctx context.Context, user domain.UserIntf) error {
	db := u.GetConnection(ctx)
	us := new(User)
	us.FromDomainUser(user)
	return db.Model(&User{}).Create(&us).Error
}

// DeleteOrganizationTenant 删除商户组织关系
func (u *UserStore) DeleteOrganizationTenant(ctx context.Context, organizationID string, tenantID string) error {
	db := u.GetConnection(ctx)
	return db.Where("organization_id = ? and tenant_id = ?", organizationID, tenantID).Delete(&OrganizationTenant{}).Error
}

// DeleteStaff 删除员工
func (u *UserStore) DeleteStaff(ctx context.Context, staffId string) error {
	return u.GetConnection(ctx).Model(&User{}).Where("user_id = ?", staffId).Updates(map[string]interface{}{
		"is_activated": 0,
	}).Error
}

// GetActivateStaffByPhone 通过手机号获取已经激活员工
func (u *UserStore) GetActivateStaffByPhone(ctx context.Context, phone string) (domain.UserIntf, error) {
	var staff User
	err := u.GetConnection(ctx).Model(&User{}).Where("phone = ? and role_type = ? and is_activated = 1", phone, domain.RoleTypeStaff).First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return staff.ToDomainUser(), nil
}

// GetLatestTenantEntityRevision 获取最新的 revision
func (u *UserStore) GetLatestTenantEntityRevision(ctx context.Context, tenantId string) (domain.TenantEntityRevisionIntf, error) {
	var re TenantEntityRevision
	err := u.GetConnection(ctx).Model(&TenantEntityRevision{}).Where("tenant_id = ?", tenantId).Order("created_at desc").First(&re).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return re.ToDomainTenantEntityRevision(), nil
}

// GetLatestTenantSubscriptionTimeline 获取最新的商户时间线
func (u *UserStore) GetLatestTenantSubscriptionTimeline(ctx context.Context, organizationId string, tenantId string) (domain.SubscriptionTimelineIntf, error) {
	db := u.GetConnection(ctx)
	var us SubscriptionTimeline
	err := db.Model(&SubscriptionTimeline{}).Where("tenant_id = ?", tenantId).Order("created_at desc").First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return us.ToDomainSubscriptionTimeline(), nil
}

// GetOrganizationByOrganizationId 通过组织 ID 查询组织信息
func (u *UserStore) GetOrganizationByOrganizationId(ctx context.Context, organizationID string) (domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	us := &User{}
	err := db.Model(&User{}).Where("organization_id = ?", organizationID).First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return us.ToDomainUser(), nil
}

// GetOrganizationByPhone 通过手机号查询组织信息
func (u *UserStore) GetOrganizationByPhone(ctx context.Context, phone string) (domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	us := &User{}
	err := db.Model(&User{}).Where("phone = ? and role_type = ?", phone, domain.RoleTypeOrganization).First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return us.ToDomainUser(), nil
}

// GetOrganizationByUsername 通过用户名查询组织信息
func (u *UserStore) GetOrganizationByUsername(ctx context.Context, username string) (domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	us := &User{}
	err := db.Model(&User{}).Where("username = ? and role_type = ?", username, domain.RoleTypeOrganization).First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return us.ToDomainUser(), nil
}

// GetOrganizationByVagueUsername 通过模糊用户名查询组织信息
func (u *UserStore) GetOrganizationByVagueUsername(ctx context.Context, username string) ([]domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	var us []*User
	vagueName := "%" + username + "%"
	err := db.Model(&User{}).Where("username like ? and role_type = ?", vagueName, domain.RoleTypeOrganization).Find(&us).Error
	if err != nil {
		return nil, err
	}

	// 转化
	dus := make([]domain.UserIntf, len(us))
	for k, v := range us {
		dus[k] = v.ToDomainUser()
	}
	return dus, nil
}

// GetOrganizationTenant 获取组织商户信息
func (u *UserStore) GetOrganizationTenant(ctx context.Context, organizationID string, tenantID string) (domain.OrganizationTenantIntf, error) {
	db := u.GetConnection(ctx)
	ot := new(OrganizationTenant)
	err := db.Model(&OrganizationTenant{}).Where("organization_id = ? and tenant_id = ?", organizationID, tenantID).First(&ot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ot.ToDomainOrganizationTenant(), nil
}

// GetOrganizationTenantByTenantId 通过商户 ID 获取组织商户信息
func (u *UserStore) GetOrganizationTenantByTenantId(ctx context.Context, tenantID string) (domain.OrganizationTenantIntf, error) {
	var ot OrganizationTenant
	err := u.GetConnection(ctx).Model(&OrganizationTenant{}).Where("tenant_id = ?", tenantID).First(&ot).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ot.ToDomainOrganizationTenant(), nil
}

// GetStaffByID 根据员工 ID 获取信息
func (u *UserStore) GetStaffByID(ctx context.Context, staffId string) (domain.UserIntf, error) {
	var staff User
	err := u.GetConnection(ctx).Model(&User{}).Where("user_id = ?", staffId).First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return staff.ToDomainUser(), nil
}

// GetStaffByPhone 通过手机号查询员工
func (u *UserStore) GetStaffByPhone(ctx context.Context, areaCode string, phone string) (domain.UserIntf, error) {
	// 获取用户的数据
	var us User
	err := u.GetConnection(ctx).Model(&User{}).Where("phone = ? AND role_type = ? AND is_activated = 1", phone, domain.RoleTypeStaff).First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	// 返回用户
	return us.ToDomainUser(), nil
}

// GetSubscriptionTimeline 获取订阅时间线
func (u *UserStore) GetSubscriptionTimeline(ctx context.Context, tenantId string) (domain.SubscriptionTimelineIntf, error) {
	var timeline SubscriptionTimeline
	err := u.GetConnection(ctx).Model(&SubscriptionTimeline{}).Where("tenant_id = ?", tenantId).First(&timeline).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return timeline.ToDomainSubscriptionTimeline(), nil
}

// GetSystemUserByID 通过用户 ID 获取系统用户
func (u *UserStore) GetSystemUserByID(ctx context.Context, systemUserID string) (domain.UserIntf, error) {
	var us User
	err := u.GetConnection(ctx).Model(&User{}).Where("user_id = ? and role_type >= ?", systemUserID, domain.RoleTypeJmAdmin).First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return us.ToDomainUser(), nil
}

// GetSystemUserByPhone 通过手机号获取系统用户
func (u *UserStore) GetSystemUserByPhone(ctx context.Context, phone string) (domain.UserIntf, error) {
	var us User
	err := u.GetConnection(ctx).Model(&User{}).Where("phone = ? and role_type >= ? and is_activated = 1", phone, domain.RoleTypeJmAdmin).First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return us.ToDomainUser(), nil
}

// GetTenantEntity 获取租户实体信息
func (u *UserStore) GetTenantEntity(ctx context.Context, tid string) (domain.TenantEntityIntf, error) {
	db := u.GetConnection(ctx)
	// 获取租户实体信息
	var entity TenantEntity
	err := db.Model(&TenantEntity{}).Where("tenant_id = ?", tid).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return entity.ToDomainTenantEntity(), nil
}

// GetTenantEntityRevision 获取最新的 revision
func (u *UserStore) GetTenantEntityRevision(ctx context.Context, revisionId string) (domain.TenantEntityRevisionIntf, error) {
	// 获取商户信息审核副本
	var re TenantEntityRevision
	err := u.GetConnection(ctx).Model(&TenantEntityRevision{}).Where("tenant_entity_revision_id = ?", revisionId).First(&re).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return re.ToDomainTenantEntityRevision(), nil
}

// GetTenantStaffByPhone 通过手机号获取商户员工
func (u *UserStore) GetTenantStaffByPhone(ctx context.Context, tenantId string, phone string) (domain.UserIntf, error) {
	var staff User
	err := u.GetConnection(ctx).Model(&User{}).Where("tenant_id = ? and phone = ? and role_type = ?", tenantId, phone, domain.RoleTypeStaff).First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return staff.ToDomainUser(), nil
}

// GetTenantUser 获取租户 User 信息
func (u *UserStore) GetTenantUser(ctx context.Context, tid string) (domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	us := new(User)
	err := db.Model(&User{}).Where("tenant_id = ? AND role_type = ?", tid, domain.RoleTypeTenant).First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return us.ToDomainUser(), nil
}

// GetTenantUserByPhone 通过手机号获取租户 User 信息
func (u *UserStore) GetTenantUserByPhone(ctx context.Context, phone string) (domain.UserIntf, error) {
	db := u.GetConnection(ctx)
	us := new(User)
	err := db.Model(&User{}).Where("phone = ? and role_type = ?", phone, domain.RoleTypeTenant).First(&us).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return us.ToDomainUser(), nil
}

// LeaveTenant 员工离开租户
func (u *UserStore) LeaveTenant(ctx context.Context, staffId string, tenantId string, rev int32) error {
	// 员工离开租户
	return u.GetConnection(ctx).Model(&User{}).Where("tenant_id = ? and user_id = ? and rev = ?", tenantId, staffId, rev).Updates(map[string]interface{}{
		"is_activated": 0,
		"rev":          rev + 1,
	}).Error
}

// ListActivatedStaffs 获取活跃员工列表
func (u *UserStore) ListActivatedStaffs(ctx context.Context, tid string) ([]domain.UserIntf, error) {
	var staffs []User
	err := u.GetConnection(ctx).Model(&User{}).Where("tenant_id = ? and role_type = ? and is_activated = 1", tid, domain.RoleTypeStaff).Find(&staffs).Error
	if err != nil {
		return nil, err
	}
	dStaffs := make([]domain.UserIntf, len(staffs))
	for k, v := range staffs {
		dStaffs[k] = v.ToDomainUser()
	}
	return dStaffs, nil
}

// ListAllOrganizationTenantIds 获取组织下所有商户 ID
func (u *UserStore) ListAllOrganizationTenantIds(ctx context.Context, oid string) ([]string, error) {
	var tids []string
	err := u.GetConnection(ctx).Raw("select tenant_id from organization_tenant where organization_id = ? AND deleted_at IS NULL", oid).Scan(&tids).Error
	if err != nil {
		return nil, err
	}
	return tids, nil
}

// ListAllOrganizations 获取所有组织信息
func (u *UserStore) ListAllOrganizations(ctx context.Context) ([]domain.UserIntf, error) {
	var us []User
	// 查询组织
	err := u.GetConnection(ctx).Model(&User{}).Where("role_type = ?", domain.RoleTypeOrganization).Find(&us).Error
	if err != nil {
		return nil, err
	}

	dus := make([]domain.UserIntf, len(us))
	for k, v := range us {
		dus[k] = v.ToDomainUser()
	}
	return dus, nil
}

// ListAllStaffs 获取所有的员工包括删除的
func (u *UserStore) ListAllStaffs(ctx context.Context, tenantId string) ([]domain.UserIntf, error) {
	var users []User
	// 查询所有员工
	err := u.GetConnection(ctx).Model(&User{}).Where("tenant_id = ? and role_type = ?", tenantId, domain.RoleTypeStaff).Scan(&users).Error
	if err != nil {
		return nil, err
	}
	dUser := make([]domain.UserIntf, len(users))
	for k, v := range users {
		dUser[k] = v.ToDomainUser()
	}
	return dUser, nil
}

// ListOrganizationAuthTenants 获取所有认证组织商户
func (u *UserStore) ListOrganizationAuthTenants(ctx context.Context, organizationID string) ([]domain.OrganizationTenantIntf, error) {
	db := u.GetConnection(ctx)
	ots := []*OrganizationTenant{}
	// 查询
	err := db.Model(&OrganizationTenant{}).Where("organization_id = ? and is_activated = ?", organizationID, 1).Find(&ots).Error
	if err != nil {
		return nil, err
	}
	dots := make([]domain.OrganizationTenantIntf, len(ots))
	for k, v := range ots {
		dots[k] = v.ToDomainOrganizationTenant()
	}
	return dots, nil
}

// ListOrganizationExistTenants 获取所有未删除的商户 ID
func (u *UserStore) ListOrganizationExistTenants(ctx context.Context, organizationID string) ([]domain.OrganizationTenantIntf, error) {
	db := u.GetConnection(ctx)
	ots := []*OrganizationTenant{}
	// 查询
	err := db.Model(&OrganizationTenant{}).Where("organization_id = ?", organizationID).Find(&ots).Error
	if err != nil {
		return nil, err
	}
	dots := make([]domain.OrganizationTenantIntf, len(ots))
	for k, v := range ots {
		dots[k] = v.ToDomainOrganizationTenant()
	}
	return dots, nil
}

// ListOrganizationTenants 获取组织商户列表信息
func (u *UserStore) ListOrganizationTenants(ctx context.Context, organizationID string, tenantID []string) ([]domain.OrganizationTenantIntf, error) {
	db := u.GetConnection(ctx)

	ots := []OrganizationTenant{}
	err := db.Model(&OrganizationTenant{}).Where("organization_id = ? and tenant_id in (?)", organizationID, tenantID).Find(&ots).Error
	if err != nil {
		return nil, err
	}
	dots := make([]domain.OrganizationTenantIntf, len(ots))
	for k, v := range ots {
		dots[k] = v.ToDomainOrganizationTenant()
	}
	return dots, nil
}

// ListOrganizationTenantsWithoutTenantId 获取组织商户列表信息
func (u *UserStore) ListOrganizationTenantsWithoutTenantId(ctx context.Context, organizationID string, offset int32, size int32) ([]domain.OrganizationTenantIntf, int64, error) {
	db := u.GetConnection(ctx)
	var count int64
	ots := []OrganizationTenant{}
	err := db.Model(&OrganizationTenant{}).Where("organization_id = ?", organizationID).
		Offset(int(offset)).
		Limit(int(size)).
		Find(&ots).
		Error
	if err != nil {
		return nil, count, err
	}
	err = db.Model(&OrganizationTenant{}).Where("organization_id = ?", organizationID).
		Count(&count).
		Error
	if err != nil {
		return nil, count, err
	}

	dots := make([]domain.OrganizationTenantIntf, len(ots))
	for k, v := range ots {
		dots[k] = v.ToDomainOrganizationTenant()
	}
	return dots, count, nil
}

// ListOrganizations 获取组织
func (u *UserStore) ListOrganizations(ctx context.Context, oids []string) ([]domain.UserIntf, error) {
	var us []User
	// 查询组织
	err := u.GetConnection(ctx).Model(&User{}).Where("organization_id in (?) and role_type = ?", oids, domain.RoleTypeOrganization).Find(&us).Error
	if err != nil {
		return nil, err
	}

	dus := make([]domain.UserIntf, len(us))
	for k, v := range us {
		dus[k] = v.ToDomainUser()
	}
	return dus, nil
}

// ListStaffs 获取员工列表
func (u *UserStore) ListStaffs(ctx context.Context, tid string) ([]domain.UserIntf, error) {
	var staffs []User
	err := u.GetConnection(ctx).Model(&User{}).Where("tenant_id = ? and role_type = ?", tid, domain.RoleTypeStaff).Find(&staffs).Error
	if err != nil {
		return nil, err
	}
	dStaffs := make([]domain.UserIntf, len(staffs))
	for k, v := range staffs {
		dStaffs[k] = v.ToDomainUser()
	}
	return dStaffs, nil
}

// ListSubscriptionTimeline 获取订阅时间线
func (u *UserStore) ListSubscriptionTimeline(ctx context.Context, tenantIds []string) ([]domain.SubscriptionTimelineIntf, error) {
	var timelines []SubscriptionTimeline
	err := u.GetConnection(ctx).Model(&SubscriptionTimeline{}).Where("tenant_id IN (?)", tenantIds).Find(&timelines).Error
	if err != nil {
		return nil, err
	}
	dTs := make([]domain.SubscriptionTimelineIntf, len(timelines))
	for k, v := range timelines {
		dTs[k] = v.ToDomainSubscriptionTimeline()
	}
	return dTs, nil
}

// ListSystemOperators 获取系统操作员
func (u *UserStore) ListSystemOperators(ctx context.Context) ([]domain.UserIntf, error) {
	var users []User
	err := u.GetConnection(ctx).Model(&User{}).Where("role_type >= ?", domain.RoleTypeJmAdmin).Find(&users).Error
	if err != nil {
		return nil, err
	}
	du := make([]domain.UserIntf, len(users))
	for k, v := range users {
		du[k] = v.ToDomainUser()
	}
	return du, nil
}

// ListSystemUsers 获取系统用户列表
func (u *UserStore) ListSystemUsers(ctx context.Context, userIds []string) ([]domain.UserIntf, error) {
	var users []User
	err := u.GetConnection(ctx).Model(&User{}).Where("user_id IN (?) AND role_type >= ?", userIds, domain.RoleTypeJmAdmin).Find(&users).Error
	if err != nil {
		return nil, err
	}
	du := make([]domain.UserIntf, len(users))
	for k, v := range users {
		du[k] = v.ToDomainUser()
	}
	return du, nil
}

// ListTenantEntity 获取商户信息列表
func (u *UserStore) ListTenantEntity(ctx context.Context, tids []string) ([]domain.TenantEntityIntf, error) {
	db := u.GetConnection(ctx)
	// 获取租户实体信息
	var entity []TenantEntity
	err := db.Model(&TenantEntity{}).Where("tenant_id IN (?)", tids).Find(&entity).Error
	if err != nil {
		return nil, err
	}
	de := make([]domain.TenantEntityIntf, len(entity))
	for k, v := range entity {
		de[k] = v.ToDomainTenantEntity()
	}
	return de, nil
}

// ListTenants 获取商户列表信息
func (u *UserStore) ListTenants(ctx context.Context, tenantID []string) ([]domain.OrganizationTenantIntf, error) {
	db := u.GetConnection(ctx)

	ots := []OrganizationTenant{}
	err := db.Model(&OrganizationTenant{}).Where("tenant_id in (?)", tenantID).Find(&ots).Error
	if err != nil {
		return nil, err
	}
	dots := make([]domain.OrganizationTenantIntf, len(ots))
	for k, v := range ots {
		dots[k] = v.ToDomainOrganizationTenant()
	}
	return dots, nil
}

// ListTenantsActivationTimeline 获取商户激活时间线
func (u *UserStore) ListTenantsActivationTimeline(ctx context.Context) ([]domain.SubscriptionTimelineIntf, error) {
	var ats []SubscriptionTimeline
	err := u.GetConnection(ctx).Model(&SubscriptionTimeline{}).Order("tenant_id,start_time asc").Find(&ats).Error
	if err != nil {
		return nil, err
	}
	dats := make([]domain.SubscriptionTimelineIntf, len(ats))
	for k, v := range ats {
		dats[k] = v.ToDomainSubscriptionTimeline()
	}
	return dats, nil
}

// ListTenantsContainsDeleted 根据 ID 获取商户信息包括已删除的
func (u *UserStore) ListTenantsContainsDeleted(ctx context.Context, tenantID []string) ([]domain.OrganizationTenantIntf, error) {
	db := u.GetConnection(ctx)

	ots := []OrganizationTenant{}
	err := db.Raw(`SELECT * FROM
		organization_tenant
		WHERE
		tenant_id in (?)
		AND deleted_at IS NULL`, tenantID).Find(&ots).Error
	if err != nil {
		return nil, err
	}
	dots := make([]domain.OrganizationTenantIntf, len(ots))
	for k, v := range ots {
		dots[k] = v.ToDomainOrganizationTenant()
	}
	return dots, nil
}

// ListTenantsStatus 获取商户列表状态
func (u *UserStore) ListTenantsStatus(ctx context.Context, tenantID []string) ([]domain.OrganizationTenantIntf, error) {
	db := u.GetConnection(ctx)

	ots := []OrganizationTenant{}
	err := db.Model(&OrganizationTenant{}).Where("tenant_id in (?)", tenantID).Find(&ots).Error
	if err != nil {
		return nil, err
	}
	dots := make([]domain.OrganizationTenantIntf, len(ots))
	for k, v := range ots {
		dots[k] = v.ToDomainOrganizationTenant()
	}
	return dots, nil
}

// ModifyReportSharingStatus 修改报告共享状态
func (u *UserStore) ModifyReportSharingStatus(ctx context.Context, tid string, status bool, rev int32) error {
	db := u.GetConnection(ctx)
	// 修改商品推荐状态
	return db.Model(&TenantEntity{}).Where("tenant_id = ? and rev = ?", tid, rev).Updates(map[string]interface{}{
		"report_sharing_status": status,
		"rev":                   rev + 1,
	}).Error
}

// RenewSubscriptionTimeline 续费订阅（未过期）
func (u *UserStore) RenewSubscriptionTimeline(ctx context.Context, tenantIds []string, renewDays int32) error {
	return u.GetConnection(ctx).Model(&SubscriptionTimeline{}).Exec(`
        UPDATE 
        subscription_timeline 
        SET 
        end_time = DATE_ADD(end_time, INTERVAL ? DAY)  
        WHERE 
        tenant_id IN (?)
    `, renewDays, tenantIds).Error
}

// UpdateOrganizationPassword 更新组织密码
func (u *UserStore) UpdateOrganizationPassword(ctx context.Context, oid string, hashedPassword string, rev int32) error {
	return u.GetConnection(ctx).Model(&User{}).Where("organization_id = ? and role_type = ? and rev = ?", oid, domain.RoleTypeOrganization, rev).Updates(map[string]interface{}{
		"rev":             rev + 1,
		"hashed_password": hashedPassword,
	}).Error
}

// UpdateOrganizationTenantStatus 更新商户审核状态
func (u *UserStore) UpdateOrganizationTenantStatus(ctx context.Context, tenantId string, reviewStatus int32, rev int32) error {
	if reviewStatus == domain.TenantReviewStatusReviewSuccess {
		return u.GetConnection(ctx).Model(&OrganizationTenant{}).Where("tenant_id = ? and rev = ?", tenantId, rev).Updates(map[string]interface{}{
			"is_activated":  true,
			"review_status": reviewStatus,
			"rev":           rev + 1,
		}).Error
	}
	return u.GetConnection(ctx).Model(&OrganizationTenant{}).Where("tenant_id = ? and rev = ?", tenantId, rev).Updates(map[string]interface{}{
		"review_status": reviewStatus,
		"rev":           rev + 1,
	}).Error
}

// UpdateOrganizationUserPassword 更新组织密码
func (u *UserStore) UpdateOrganizationUserPassword(ctx context.Context, organizationId string, newHashedPassword string, rev int32) error {
	db := u.GetConnection(ctx)
	return db.Model(&User{}).Where("organization_id = ? and rev = ? and role_type = ?", organizationId, rev, domain.RoleTypeOrganization).Updates(map[string]interface{}{
		"hashed_password": newHashedPassword,
		"rev":             rev + 1,
	}).Error
}

// UpdateStaff implements 更新员工信息
func (u *UserStore) UpdateStaff(ctx context.Context, staffId string, phone string, name string, hashedPassword string, rev int32) error {
	// 构建更新map
	updateMap := map[string]interface{}{
		"phone":    phone,
		"nickname": name,
		"rev":      rev + 1,
	}
	if hashedPassword != "" {
		updateMap["hashed_password"] = hashedPassword
	}
	// 更新员工信息
	return u.GetConnection(ctx).Model(&User{}).Where("user_id = ? and rev = ?", staffId, rev).Updates(updateMap).Error
}

// UpdateStaffNickname 更新员工的昵称
func (u *UserStore) UpdateStaffNickname(ctx context.Context, staffId string, newName string, rev int32) error {
	return u.GetConnection(ctx).Model(&User{}).Where("user_id = ? and role_type = ? and rev = ?", staffId, domain.RoleTypeStaff, rev).
		Update("nickname", newName).Error
}

// UpdateStaffPassword 更新员工密码
func (u *UserStore) UpdateStaffPassword(ctx context.Context, staffId string, hashedPassword string, rev int32) error {
	// 更新密码
	return u.GetConnection(ctx).Model(&User{}).Where("user_id = ? and rev = ?", staffId, rev).Updates(map[string]interface{}{
		"hashed_password": hashedPassword,
		"rev":             rev + 1,
	}).Error
}

// UpdateTenantEntity 更新租户实体
func (u *UserStore) UpdateTenantEntity(ctx context.Context, revision domain.TenantEntityRevisionIntf, rev int32) error {
	return u.GetConnection(ctx).Model(&TenantEntity{}).Where("tenant_id = ?", revision.GetTenantID()).Updates(map[string]interface{}{
		"logo_url":             revision.GetLogoUrl(),
		"safe_phone":           revision.GetSafePhone(),
		"store_name":           revision.GetStoreName(),
		"province":             revision.GetProvince(),
		"city":                 revision.GetCity(),
		"district":             revision.GetDistrict(),
		"street":               revision.GetStreet(),
		"contact_name":         revision.GetContactName(),
		"contact_phone":        revision.GetContactPhone(),
		"social_credit_code":   revision.GetSocialCreditCode(),
		"business_license_url": revision.GetBusinessLicenseUrl(),
		"rev":                  rev + 1,
	}).Error
}

// UpdateTenantPassword 更新商户密码
func (u *UserStore) UpdateTenantPassword(ctx context.Context, tid string, hashedPassword string, rev int32) error {
	// 更新租户密码
	return u.GetConnection(ctx).Model(&User{}).Where("tenant_id = ? and role_type = ? and rev = ?", tid, domain.RoleTypeTenant, rev).Updates(map[string]interface{}{
		"rev":             rev + 1,
		"hashed_password": hashedPassword,
	}).Error
}

// UpdateTenantUser 更新商户 user
func (u *UserStore) UpdateTenantUser(ctx context.Context, tid string, updateMap map[string]interface{}) error {
	return u.GetConnection(ctx).Model(&User{}).Where("tenant_id = ? and role_type = ?", tid, domain.RoleTypeTenant).Updates(updateMap).Error
}

// UpdateTenantUserPassword 更新租户密码
func (u *UserStore) UpdateTenantUserPassword(ctx context.Context, tenantId string, newHashedPassword string, rev int32) error {
	db := u.GetConnection(ctx)
	return db.Model(&User{}).Where("tenant_id = ? and rev = ? and role_type = ?", tenantId, rev, domain.RoleTypeTenant).Updates(map[string]interface{}{
		"hashed_password": newHashedPassword,
		"rev":             rev + 1,
	}).Error
}

// CancelRevisionIssue 取消issue
func (u *UserStore) CancelRevisionIssue(ctx context.Context, revisionId string) error {
	return u.GetConnection(ctx).Exec("update review_issue set is_cancel = 1 where target_rev_id = ? and target_type = 0", revisionId).Error
}

// ListFaqs 获取所有faq
func (u *UserStore) ListFaqs(ctx context.Context) ([]domain.FaqIntf, error) {
	db := u.GetConnection(ctx)
	faqs := []Faq{}
	err := db.Model(&Faq{}).Find(&faqs).Error
	if err != nil {
		return nil, err
	}
	dfaqs := make([]domain.FaqIntf, len(faqs))
	for k, v := range faqs {
		dfaqs[k] = v.ToDomainFaq()
	}
	return dfaqs, nil
}
