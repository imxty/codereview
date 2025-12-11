package report

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	ptime "github.com/jinmukeji/huimaibao-service/pkg/time"
	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
)

// ReportStore 实现 Domain 层 ReportStore 接口
type ReportStore struct {
	// 数据库连接
	*dbutils.Connection
}

func NewReportStore(management *dbutils.Connection) *ReportStore {
	return &ReportStore{
		management,
	}
}

// 实现 Domain 行为
var _ domain.ReportRepository = (*ReportStore)(nil)

// ListCustomerLastReports 获取常客测量的最后limit份报告
func (r *ReportStore) ListCustomerLastReports(ctx context.Context, tenantID string, customerID string, limit int) ([]domain.ReportIntf, error) {
	db := r.GetConnection(ctx)
	// 查询报告
	var reports []Report
	err := db.Model(&Report{}).
		Where("customer_id = ? and is_stress_state = 0 and is_report_complete = 1 and tenant_id = ?", customerID, tenantID).
		Order("created_at desc").
		Limit(limit).
		Find(&reports).Error
	if err != nil {
		return nil, err
	}
	// 返回数据
	dReports := make([]domain.ReportIntf, len(reports))
	for k, v := range reports {
		dReports[k] = v.ToDomainReport()
	}
	return dReports, nil
}

// GetReportIDByToken 通过报告的token获取报告ID
func (r *ReportStore) GetReportIDByToken(ctx context.Context, tenantID string, token string) (domain.SharedReportIntf, error) {
	db := r.GetConnection(ctx)
	// 通过报告的token获取报告ID
	var sr SharedReport
	err := db.Model(&SharedReport{}).Where("token = ?", token).First(&sr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return sr.ToDomainSharedReport(), nil
}

// GetReportByID 通过报告ID获取报告
func (r *ReportStore) GetReportByID(ctx context.Context, tenantID string, reportID string) (domain.ReportIntf, error) {
	db := r.GetConnection(ctx)
	// 查找报告
	var report Report
	err := db.Model(&Report{}).Where("report_id = ?", reportID).First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return report.ToDomainReport(), nil
}

// GetReportByReportID 通过报告ID获取报告
func (r *ReportStore) GetReportByReportID(ctx context.Context, reportID string) (domain.ReportIntf, error) {
	db := r.GetConnection(ctx)
	// 查找报告
	var report Report
	err := db.Model(&Report{}).Where("report_id = ?", reportID).First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return report.ToDomainReport(), nil
}

// UpdateReportTongueFaceReportStatus 更新报告的舌面分析报告状态
func (r *ReportStore) UpdateReportTongueFaceReportStatus(ctx context.Context, faceUrl, tongueUrl, reportID, body string, status int32) error {
	db := r.GetConnection(ctx)
	// 更新报告
	err := db.Model(&Report{}).Where("report_id = ?", reportID).Updates(map[string]interface{}{
		"face_image_url":            faceUrl,
		"tongue_image_url":          tongueUrl,
		"tongue_face_report_status": status,
		"tongue_face_report":        body,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetCustomerReport 获取常客报告
func (r *ReportStore) GetCustomerReport(ctx context.Context, tenantID string, reportID, customerID string) (domain.ReportIntf, error) {
	db := r.GetConnection(ctx)
	// 查找报告
	var report Report
	err := db.Model(&Report{}).Where("report_id = ? and customer_id = ? and is_report_complete = 1 and tenant_id = ?", reportID, customerID, tenantID).First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return report.ToDomainReport(), nil
}

// CreateSharedReport 创建分享报告
func (r *ReportStore) CreateSharedReport(ctx context.Context, sharedReport domain.SharedReportIntf) error {
	db := r.GetConnection(ctx)
	// 创建分享报告
	var sr SharedReport
	sr.FromDomainSharedReport(sharedReport)
	err := db.Model(&SharedReport{}).Create(&sr).Error
	if err != nil {
		return err
	}
	return nil
}

// ListReports 获取常客/散报告列表
func (r *ReportStore) ListReports(ctx context.Context, tenantID string, customerID string, startTime, endTime time.Time, offset, size int32) ([]domain.ReportIntf, int64, error) {
	db := r.GetConnection(ctx)
	// 获取报告列表
	var count int64
	var reports []Report
	// 如果查询的是散客那么customer_id = ""
	err := db.Model(&Report{}).
		Where("created_at between ? and ? and customer_id = ? and is_report_complete = 1 and tenant_id = ?", startTime, endTime, customerID, tenantID).
		Order("created_at desc").
		Offset(int(offset)).
		Limit(int(size)).
		Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Model(&Report{}).Where("customer_id = ? and created_at between ? and ? and is_report_complete = 1 and tenant_id = ?", customerID, startTime, endTime, tenantID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	// 返回报告
	dReports := make([]domain.ReportIntf, len(reports))
	for k, v := range reports {
		dReports[k] = v.ToDomainReport()
	}
	return dReports, count, nil
}

// ListCustomerReports 获取常客报告列表
func (r *ReportStore) ListCustomerReports(ctx context.Context, tenantId, customerId string, startTime, endTime time.Time) ([]domain.ReportIntf, error) {
	db := r.GetConnection(ctx)
	// 获取报告列表
	var reports []Report
	err := db.Model(&Report{}).Where("customer_id = ? and created_at between ? and ? and is_report_complete = 1 and tenant_id = ?", customerId, startTime, endTime, tenantId).Order("created_at asc").Find(&reports).Error
	if err != nil {
		return nil, err
	}
	// 返回报告
	dReports := make([]domain.ReportIntf, len(reports))
	for k, v := range reports {
		dReports[k] = v.ToDomainReport()
	}
	return dReports, nil
}

// ListReportsWithStaff 获取员工常客/散报告列表
func (r *ReportStore) ListReportsWithStaff(ctx context.Context, tenantID string, customerID, staffId string, reportSharingStatus bool, startTime, endTime time.Time, offset, size int32) ([]domain.ReportIntf, int64, error) {
	db := r.GetConnection(ctx)
	// 获取报告列表
	var reports []Report
	var count int64
	// 如果共享报告开关关闭则需要根据staffID进行查询
	if !reportSharingStatus {
		// 如果查询的是散客那么customer_id = ""
		err := db.Model(&Report{}).
			Where("created_at between ? and ? and customer_id = ? and is_report_complete = 1 and tenant_id = ? and staff_id = ?", startTime, endTime, customerID, tenantID, staffId).
			Order("created_at desc").
			Offset(int(offset)).
			Limit(int(size)).
			Find(&reports).Error
		if err != nil {
			return nil, count, err
		}
		// 如果查询的是散客那么customer_id = ""
		err = db.Model(&Report{}).
			Where("created_at between ? and ? and customer_id = ? and is_report_complete = 1 and tenant_id = ? and staff_id = ?", startTime, endTime, customerID, tenantID, staffId).
			Count(&count).Error
		if err != nil {
			return nil, count, err
		}
	} else {
		// 如果查询的是散客那么customer_id = ""
		err := db.Model(&Report{}).
			Where("created_at between ? and ? and customer_id = ? and is_report_complete = 1 and tenant_id = ?", startTime, endTime, customerID, tenantID).
			Order("created_at desc").
			Offset(int(offset)).
			Limit(int(size)).
			Find(&reports).Error
		if err != nil {
			return nil, count, err
		}
		// 如果查询的是散客那么customer_id = ""
		err = db.Model(&Report{}).
			Where("created_at between ? and ? and customer_id = ? and is_report_complete = 1 and tenant_id = ?", startTime, endTime, customerID, tenantID).
			Count(&count).Error
		if err != nil {
			return nil, count, err
		}
	}

	// 返回报告
	dReports := make([]domain.ReportIntf, len(reports))
	for k, v := range reports {
		dReports[k] = v.ToDomainReport()
	}
	return dReports, count, nil
}

// ListReportsWithoutTime 获取常客/散客报告列表(没有时间区间)
func (r *ReportStore) ListReportsWithoutTime(ctx context.Context, tenantID string, customerID string, offset, size int32) ([]domain.ReportIntf, int64, error) {
	db := r.GetConnection(ctx)
	// 获取报告列表
	var count int64
	var reports []Report
	// 如果查询的是散客那么customer_id = ""
	err := db.Model(&Report{}).
		Where("customer_id = ? and is_report_complete = 1 and tenant_id = ?", customerID, tenantID).
		Order("created_at desc").
		Offset(int(offset)).
		Limit(int(size)).
		Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Model(&Report{}).Where("customer_id = ? and is_report_complete = 1 and tenant_id = ?", customerID, tenantID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	// 返回报告
	dReports := make([]domain.ReportIntf, len(reports))
	for k, v := range reports {
		dReports[k] = v.ToDomainReport()
	}
	return dReports, count, nil
}

// ModifyReportRemark 修改报告的普通备注
func (r *ReportStore) ModifyReportRemark(ctx context.Context, tenantID string, reportID, remark string, rev int32) error {
	db := r.GetConnection(ctx)
	// 修改备注
	return db.Model(&Report{}).Where("report_id = ? and rev = ?", reportID, rev).Updates(map[string]interface{}{
		"remarks": remark,
		"rev":     rev + 1,
	}).Error
}

// ModifyReportStaffRemark 修改报告的员工备注
func (r *ReportStore) ModifyReportStaffRemark(ctx context.Context, tenantID string, reportID, remark string, rev int32) error {
	db := r.GetConnection(ctx)
	// 修改备注
	return db.Model(&Report{}).Where("report_id = ? and rev = ?", reportID, rev).Updates(map[string]interface{}{
		"staff_remarks": remark,
		"rev":           rev + 1,
	}).Error
}

// AddCustomerReport 散客报告变为常客报告
func (r *ReportStore) AddCustomerReport(ctx context.Context, tenantID string, reportId, customerID, age string, rev int32) error {
	db := r.GetConnection(ctx)
	// 散客报告变为常客报告
	return db.Model(&Report{}).Where("report_id = ? and rev = ?", reportId, rev).Updates(map[string]interface{}{
		"is_customer": 1,
		"customer_id": customerID,
		"age":         age,
		"rev":         rev + 1,
	}).Error
}

// CreateReport 创建报告
func (r *ReportStore) CreateReport(ctx context.Context, report domain.ReportIntf) error {
	db := r.GetConnection(ctx)
	// 创建计算内容
	var re Report
	re.FromDomainReport(report)
	return db.Model(&Report{}).Create(&re).Error
}

// UpdateReportContent 更新报告的内容
func (r *ReportStore) UpdateReportContent(ctx context.Context, reportID string, tenantID string, heartRate int32, isStressState bool, factor *domain.FactorInterpretation, stressStateJson, dirtyDialectJson, diseasesJson, recommendJson, physiqueJson, measurementJson string, chart *domain.MeridianBarChart, dirtyDialectFlag, rev int32) error {
	db := r.GetConnection(ctx)
	isStress := 0
	if isStressState {
		isStress = 1
	}
	return db.Model(&Report{}).Where("report_id = ? and rev = ?", reportID, rev).Updates(map[string]interface{}{
		"c0":                   chart.C0,
		"c1":                   chart.C1,
		"c2":                   chart.C2,
		"c3":                   chart.C3,
		"c4":                   chart.C4,
		"c5":                   chart.C5,
		"c6":                   chart.C6,
		"c7":                   chart.C7,
		"dirty_dialectic":      dirtyDialectJson,
		"dirty_dialectic_flag": dirtyDialectFlag,
		"dietary_advice":       recommendJson,
		"stress_state":         stressStateJson,
		"measurement_judgment": measurementJson,
		"risks":                diseasesJson,
		"is_stress_state":      isStress,
		"physical_dialectics":  physiqueJson,
		"heart_rate":           heartRate,
		"f0":                   factor.F0,
		"f1":                   factor.F1,
		"f2":                   factor.F2,
		"f3":                   factor.F3,
		"f4":                   factor.F4,
		// FIX
		"rev":                rev + 1,
		"is_report_complete": 1,
	}).Error
}

// UpdateReportTreatment 更新报告方案
func (r *ReportStore) UpdateReportTreatment(ctx context.Context, tenantID string, reportID, tenantTreatmentRevId string, rev int32) error {
	db := r.GetConnection(ctx)
	// 更新报告方案
	return db.Model(&Report{}).Where("report_id = ? and rev = ? and tenant_id = ?", reportID, rev, tenantID).Updates(map[string]interface{}{
		"tenant_treatment_rev_id": tenantTreatmentRevId,
		"rev":                     rev + 1,
	}).Error
}

// GetOrganizationReportCount 获取报告统计(当日新增常客测量次数，当日新增散客测量次数，常客累计测量次数,散客累计测量次数)
func (r *ReportStore) GetOrganizationReportCount(ctx context.Context, organizationID string, startTime, endTime time.Time) (int64, int64, int64, int64, error) {
	db := r.GetConnection(ctx)
	// 统计数据
	var todayCustomerMeasurementCount, todayTempCustomerMeasurementCount, customerMeasurementTotalCount, tempCustomerMeasurementTotalCount int64
	// 当日新增常客测量次数
	err := db.Model(&Report{}).
		Where(`is_report_complete = 1 and is_customer = 1 and created_at between ? and ? and organization_id = ?`, startTime, endTime, organizationID).
		Count(&todayCustomerMeasurementCount).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	// 当日新增散客测量次数
	err = db.Model(&Report{}).
		Where(`is_report_complete = 1 and is_customer = 0 and created_at between ? and ? and organization_id = ?`, startTime, endTime, organizationID).
		Count(&todayTempCustomerMeasurementCount).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	// 当年统计开始时间
	yearStartTime := ptime.DayBeginUTCTime(endTime.AddDate(0, 0, -365), ptime.LocBeijing)
	// 常客当年累计测量次数
	err = db.Model(&Report{}).
		Where("is_report_complete = 1 and is_customer = 1 and created_at between ? and ? and organization_id = ?", yearStartTime, endTime, organizationID).
		Count(&customerMeasurementTotalCount).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	// 散客当年累计测量次数
	err = db.Model(&Report{}).
		Where("is_report_complete = 1 and is_customer = 0 and created_at between ? and ? and organization_id = ?", yearStartTime, endTime, organizationID).
		Count(&tempCustomerMeasurementTotalCount).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return todayCustomerMeasurementCount, todayTempCustomerMeasurementCount, customerMeasurementTotalCount, tempCustomerMeasurementTotalCount, nil
}

// GetTenantReportCount 获取报告统计(当日新增常客测量次数，当日新增散客测量次数，常客累计测量次数,散客累计测量次数)
func (r *ReportStore) GetTenantReportCount(ctx context.Context, tenantID string, startTime, endTime time.Time) (int64, int64, int64, int64, error) {
	db := r.GetConnection(ctx)
	// 统计数据
	var todayCustomerMeasurementCount, todayTempCustomerMeasurementCount, customerMeasurementTotalCount, tempCustomerMeasurementTotalCount int64
	// 当日新增常客测量次数
	err := db.Model(&Report{}).
		Where(`is_report_complete = 1 and is_customer = 1 and created_at between ? and ? and tenant_id = ?`, startTime, endTime, tenantID).
		Count(&todayCustomerMeasurementCount).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	// 当日新增散客测量次数
	err = db.Model(&Report{}).
		Where(`is_report_complete = 1 and is_customer = 0 and created_at between ? and ? and tenant_id = ?`, startTime, endTime, tenantID).
		Count(&todayTempCustomerMeasurementCount).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	// 当年统计开始时间
	yearStartTime := ptime.DayBeginUTCTime(endTime.AddDate(0, 0, -365), ptime.LocBeijing)
	// 常客累计测量次数
	err = db.Model(&Report{}).
		Where("is_report_complete = 1 and is_customer = 1 and created_at between ? and ? and tenant_id = ?", yearStartTime, endTime, tenantID).
		Count(&customerMeasurementTotalCount).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	// 散客累计测量次数
	err = db.Model(&Report{}).
		Where("is_report_complete = 1 and is_customer = 0 and created_at between ? and ? and tenant_id = ?", yearStartTime, endTime, tenantID).
		Count(&tempCustomerMeasurementTotalCount).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return todayCustomerMeasurementCount, todayTempCustomerMeasurementCount, customerMeasurementTotalCount, tempCustomerMeasurementTotalCount, nil
}

// GetStaffMeasurementCount 获取员工测量统计
func (r *ReportStore) GetStaffMeasurementCount(ctx context.Context, tenantID string, staffIds []string, startTime, endTime time.Time) ([]domain.StaffReportRankIntf, error) {
	db := r.GetConnection(ctx)
	var srs []StaffReportRank
	err := db.Raw(`Select 
        staff_id,
        count(*) AS total_measurement,
        SUM(created_at between ? and ?) AS today_measurement 
        FROM
        report
        WHERE
        tenant_id = ? and is_report_complete = 1
        GROUP BY 
        staff_id
        HAVING 
        staff_id IN (?)`, startTime, endTime, tenantID, staffIds).Scan(&srs).Error
	if err != nil {
		return nil, err
	}
	// 返回结果
	dSrs := make([]domain.StaffReportRankIntf, len(srs))
	for k, v := range srs {
		dSrs[k] = v.ToDomainStaffReportRank()
	}
	return dSrs, nil
}

// ListTenantsReportCount 获取商户报告数量
func (r *ReportStore) ListTenantsReportCount(ctx context.Context, tenantID []string, startTime, endTime time.Time) ([]domain.ReportRankIntf, error) {
	db := r.GetConnection(ctx)
	var srs []ReportRank
	err := db.Raw(`Select 
        tenant_id,
        count(*) AS report_count
        FROM
        report
        WHERE
        created_at between ? and ? and
        tenant_id IN (?) and is_report_complete = 1
        GROUP BY 
        tenant_id`, startTime, endTime, tenantID).Scan(&srs).Error
	if err != nil {
		return nil, err
	}
	// 返回结果
	dSrs := make([]domain.ReportRankIntf, len(srs))
	for k, v := range srs {
		dSrs[k] = v.ToDomainReportRank()
	}
	return dSrs, nil
}

// GetTenantReportsCount 获取时间范围内商户报告数量
func (r *ReportStore) GetTenantReportsCount(ctx context.Context, tenantID string, startTime, endTime time.Time) (int64, error) {
	var reportCount int64
	// 测量次数
	err := r.GetConnection(ctx).Model(&Report{}).
		Where("is_report_complete = 1 and created_at between ? and ? and tenant_id = ?", startTime, endTime, tenantID).
		Count(&reportCount).Error
	if err != nil {
		return 0, nil
	}
	return reportCount, nil
}

// GetOrganizationReportsCount 获取时间范围内商户报告数量
func (r *ReportStore) GetOrganizationReportsCount(ctx context.Context, organizationID string, startTime, endTime time.Time) (int64, error) {
	var reportCount int64
	// 测量次数
	err := r.GetConnection(ctx).Model(&Report{}).
		Where("is_report_complete = 1 and created_at between ? and ? and organization_id = ?", startTime, endTime, organizationID).
		Count(&reportCount).Error
	if err != nil {
		return 0, nil
	}
	return reportCount, nil
}

// SearchReports 查询报告
func (r *ReportStore) SearchReports(ctx context.Context, reportId string, customerType int, customerId string, tenantId []string, dirtyDialecticsFlag int32, organizationId string, startTime, endTime time.Time, offset, size int32) ([]domain.ReportIntf, int32, error) {
	var reports []Report
	var reportCount int64
	searchSql := "SELECT * FROM report where is_report_complete = 1"
	var searchParams []interface{}
	// 拼接sql
	if reportId != "" {
		searchSql = fmt.Sprintf("%s and report_id = ?", searchSql)
		searchParams = append(searchParams, reportId)
	}
	if customerType != CustomerTypeBoth {
		if customerType == CustomerTypeCustomer {
			searchSql = fmt.Sprintf("%s and is_customer = 1", searchSql)
		} else {
			searchSql = fmt.Sprintf("%s and is_customer = 0", searchSql)
		}
	}
	if customerId != "" {
		searchSql = fmt.Sprintf("%s and customer_id = ?", searchSql)
		searchParams = append(searchParams, customerId)
	}
	if len(tenantId) != 0 {
		searchSql = fmt.Sprintf("%s and tenant_id IN (?)", searchSql)
		searchParams = append(searchParams, tenantId)
	}
	if dirtyDialecticsFlag != 0 {
		searchSql = fmt.Sprintf("%s and (dirty_dialectic_flag & ? = ?)", searchSql)
		searchParams = append(searchParams, dirtyDialecticsFlag)
		searchParams = append(searchParams, dirtyDialecticsFlag)
	}
	if organizationId != "" {
		searchSql = fmt.Sprintf("%s and organization_id = ?", searchSql)
		searchParams = append(searchParams, organizationId)
	}
	searchSql = fmt.Sprintf("%s and created_at between ? and ?", searchSql)
	searchParams = append(searchParams, startTime, endTime)
	// 先查询总数量
	countSql := fmt.Sprintf("SELECT count(*) FROM (%s) a", searchSql)
	err := r.GetConnection(ctx).Raw(countSql, searchParams...).Count(&reportCount).Error
	if err != nil {
		return nil, 0, err
	}
	// 再分页查询报告
	searchSql = fmt.Sprintf("%s ORDER BY created_at DESC LIMIT ? OFFSET ?", searchSql)
	searchParams = append(searchParams, size, offset)
	err = r.GetConnection(ctx).Raw(searchSql, searchParams...).Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}
	// 返回报告
	dReports := make([]domain.ReportIntf, len(reports))
	for k, v := range reports {
		dReports[k] = v.ToDomainReport()
	}
	return dReports, int32(reportCount), nil
}

// 查询组织测量报告数量(常客测量数量，散客测量数量，error)
func (r *ReportStore) SearchOrganizationReportsCount(ctx context.Context, oid string, tids []string, startTime, endTime time.Time) (int64, int64, error) {
	var tempCount, customerCount int64
	if len(tids) != 0 {
		err := r.GetConnection(ctx).Model(&Report{}).
			Where("organization_id = ? and tenant_id IN (?) and is_report_complete = 1 and created_at between ? and ? and is_customer = 1", oid, tids, startTime, endTime).
			Count(&customerCount).Error
		if err != nil {
			return 0, 0, err
		}
		err = r.GetConnection(ctx).Model(&Report{}).
			Where("organization_id = ? and tenant_id IN (?) and is_report_complete = 1 and created_at between ? and ? and is_customer = 0", oid, tids, startTime, endTime).
			Count(&tempCount).Error
		if err != nil {
			return 0, 0, err
		}
		return customerCount, tempCount, nil
	}
	err := r.GetConnection(ctx).Model(&Report{}).
		Where("organization_id = ? and is_report_complete = 1 and created_at between ? and ? and is_customer = 1", oid, startTime, endTime).
		Count(&customerCount).Error
	if err != nil {
		return 0, 0, err
	}
	err = r.GetConnection(ctx).Model(&Report{}).
		Where("organization_id = ? and is_report_complete = 1 and created_at between ? and ? and is_customer = 0", oid, startTime, endTime).
		Count(&tempCount).Error
	if err != nil {
		return 0, 0, err
	}
	return customerCount, tempCount, nil
}

// SearchOrganizationReports
func (r *ReportStore) SearchOrganizationReports(ctx context.Context, oid string, tids []string, startTime, endTime time.Time, offset, size int32) ([]domain.ReportCountIntf, int64, error) {
	var rcs []ReportCount
	var count int64
	if len(tids) != 0 {
		err := r.GetConnection(ctx).Raw(`
		   SELECT 
			   organization_id,
			   tenant_id,
			   DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) AS date,
			   SUM(CASE WHEN is_customer = 1 THEN 1 ELSE 0 END) AS customer_count,
			   SUM(CASE WHEN is_customer = 0 THEN 1 ELSE 0 END) AS temp_count
		   FROM 
			   report where is_report_complete = 1 and organization_id = ? and tenant_id IN (?) and created_at between ? and ?
		   GROUP BY 
			   organization_id, tenant_id, date
		   ORDER BY 
			   tenant_id, date Limit ?,?
		   `, oid, tids, startTime, endTime, offset, size).Scan(&rcs).Error
		if err != nil {
			return nil, 0, err
		}
		drcs := make([]domain.ReportCountIntf, len(rcs))
		for k, v := range rcs {
			drcs[k] = v.ToDomainReportCount()
		}
		err = r.GetConnection(ctx).Raw(`
		   SELECT 
			   count(*) 
		   FROM
			   (
				   SELECT 
					   organization_id,
					   tenant_id,
					   DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) AS date,
					   SUM(CASE WHEN is_customer = 1 THEN 1 ELSE 0 END) AS customer_count,
					   SUM(CASE WHEN is_customer = 0 THEN 1 ELSE 0 END) AS temp_count
				   FROM 
					   report where is_report_complete = 1 and organization_id = ? and tenant_id IN (?) and created_at between ? and ?
				   GROUP BY 
					   tenant_id, date
				   ORDER BY 
					   tenant_id, date
			   ) a
		   `, oid, tids, startTime, endTime).Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
		return drcs, count, nil
	}
	err := r.GetConnection(ctx).Raw(`
		   SELECT 
			   organization_id,
			   tenant_id,
			   DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) AS date,
			   SUM(CASE WHEN is_customer = 1 THEN 1 ELSE 0 END) AS customer_count,
			   SUM(CASE WHEN is_customer = 0 THEN 1 ELSE 0 END) AS temp_count
		   FROM 
			   report where is_report_complete = 1 and organization_id = ? and created_at between ? and ?
		   GROUP BY 
			   organization_id, tenant_id, date
		   ORDER BY 
			   tenant_id, date Limit ?,?
		   `, oid, startTime, endTime, offset, size).Scan(&rcs).Error
	if err != nil {
		return nil, 0, err
	}
	drcs := make([]domain.ReportCountIntf, len(rcs))
	for k, v := range rcs {
		drcs[k] = v.ToDomainReportCount()
	}
	err = r.GetConnection(ctx).Raw(`
		   SELECT 
			   count(*) 
		   FROM
			   (
				   SELECT 
					   organization_id,
					   tenant_id,
					   DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) AS date,
					   SUM(CASE WHEN is_customer = 1 THEN 1 ELSE 0 END) AS customer_count,
					   SUM(CASE WHEN is_customer = 0 THEN 1 ELSE 0 END) AS temp_count
				   FROM 
					   report where is_report_complete = 1 and organization_id = ? and created_at between ? and ?
				   GROUP BY 
					   tenant_id, date
				   ORDER BY 
					   tenant_id, date
			   ) a
		   `, oid, startTime, endTime).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return drcs, count, nil
}

// 查询组织测量报告数量(常客测量数量，散客测量数量，error)
func (r *ReportStore) SearchStaffReportsCount(ctx context.Context, tid string, sid []string, startTime, endTime time.Time) (int64, int64, error) {
	var tempCount, customerCount int64
	err := r.GetConnection(ctx).Model(&Report{}).
		Where("tenant_id = ? and staff_id IN (?) and is_report_complete = 1 and created_at between ? and ? and is_customer = 1", tid, sid, startTime, endTime).
		Count(&customerCount).Error
	if err != nil {
		return 0, 0, err
	}
	err = r.GetConnection(ctx).Model(&Report{}).
		Where("tenant_id = ? and staff_id IN (?) and is_report_complete = 1 and created_at between ? and ? and is_customer = 0", tid, sid, startTime, endTime).
		Count(&tempCount).Error
	if err != nil {
		return 0, 0, err
	}
	return customerCount, tempCount, nil
}

// SearchStaffReports
func (r *ReportStore) SearchStaffReports(ctx context.Context, tid string, sid []string, startTime, endTime time.Time, offset, size int32) ([]domain.ReportStaffCountIntf, int64, error) {
	var rcs []ReportStaffCount
	var count int64
	err := r.GetConnection(ctx).Raw(`
        SELECT 
            organization_id,
            tenant_id,
            staff_id,
            DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) AS date,
            SUM(CASE WHEN is_customer = 1 THEN 1 ELSE 0 END) AS customer_count,
            SUM(CASE WHEN is_customer = 0 THEN 1 ELSE 0 END) AS temp_count
        FROM 
            report where is_report_complete = 1 and tenant_id = ? and staff_id IN (?) and created_at between ? and ?
        GROUP BY 
            staff_id, date
        ORDER BY 
            staff_id, date Limit ?,?
        `, tid, sid, startTime, endTime, offset, size).Scan(&rcs).Error
	if err != nil {
		return nil, 0, err
	}
	drcs := make([]domain.ReportStaffCountIntf, len(rcs))
	for k, v := range rcs {
		drcs[k] = v.ToDomainReportStaffCount()
	}
	err = r.GetConnection(ctx).Raw(`
        SELECT 
            count(*) 
        FROM
            (
                SELECT 
                    organization_id,
                    tenant_id,
                    staff_id,
                    DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) AS date,
                    SUM(CASE WHEN is_customer = 1 THEN 1 ELSE 0 END) AS customer_count,
                    SUM(CASE WHEN is_customer = 0 THEN 1 ELSE 0 END) AS temp_count
                FROM 
                    report where is_report_complete = 1 and tenant_id = ? and staff_id IN (?) and created_at between ? and ?
                GROUP BY 
                    staff_id, date
                ORDER BY 
                    staff_id, date
            ) a
        `, tid, sid, startTime, endTime).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return drcs, count, nil
}

// SearchTenantReportsDownload
func (r *ReportStore) SearchTenantReportsDownload(ctx context.Context, oid string, tid []string, startTime, endTime time.Time) ([]domain.ReportTenantCountIntf, int64, error) {
	var tcs []ReportTenantCount
	var counts int64

	// 查询所有满足条件的记录（不分页）
	err := r.GetConnection(ctx).Raw(`
		SELECT 
			r1.organization_id,
			r1.tenant_id,
			DATE_FORMAT(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), '%Y-%m') AS date,
			COUNT(r1.report_id) AS monthly_count,
			-- 环比
			IFNULL((
				COUNT(r1.report_id) - 
				(SELECT COUNT(r2.report_id)
				FROM report r2
				WHERE r2.organization_id = r1.organization_id
				AND r2.tenant_id = r1.tenant_id
				AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
					DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 1 MONTH), '%Y-%m')
				)
			) / NULLIF(
				(SELECT COUNT(r2.report_id)
				FROM report r2
				WHERE r2.organization_id = r1.organization_id
				AND r2.tenant_id = r1.tenant_id
				AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
					DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 1 MONTH), '%Y-%m')
				), 0) * 100, 2147483647) AS month_on_month,

			-- 同比
			IFNULL((
				COUNT(r1.report_id) - 
				(SELECT COUNT(r2.report_id)
				FROM report r2
				WHERE r2.organization_id = r1.organization_id
				AND r2.tenant_id = r1.tenant_id
				AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
					DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 12 MONTH), '%Y-%m')
				)
			) / NULLIF(
				(SELECT COUNT(r2.report_id)
				FROM report r2
				WHERE r2.organization_id = r1.organization_id
				AND r2.tenant_id = r1.tenant_id
				AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
					DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 12 MONTH), '%Y-%m')
				), 0) * 100, 2147483647) AS year_on_year

		FROM report r1
		WHERE r1.organization_id = ?
			AND r1.tenant_id IN (?)
			AND r1.is_report_complete = 1
			AND r1.created_at BETWEEN ? AND ?
		GROUP BY r1.organization_id, r1.tenant_id, date
		ORDER BY r1.tenant_id, date
	`, oid, tid, startTime, endTime).Scan(&tcs).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	counts = int64(len(tcs))

	// 转换为接口类型
	dtcs := make([]domain.ReportTenantCountIntf, len(tcs))
	for k, v := range tcs {
		dtcs[k] = v.ToDomainReportTenantCount()
	}

	return dtcs, counts, nil
}

// SearchTenantReports
func (r *ReportStore) SearchTenantReports(ctx context.Context, oid string, tid []string, startTime, endTime time.Time, offset, size int32) ([]domain.ReportTenantCountIntf, int64, error) {
	var tcs []ReportTenantCount
	var counts int64
	// GetConnection
	err := r.GetConnection(ctx).Raw(`
    SELECT 
		r1.organization_id,
		r1.tenant_id,
		DATE_FORMAT(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), '%Y-%m') AS date,
		COUNT(r1.report_id) AS monthly_count,
		-- 环比：当前月和前一个月的比较
		IFNULL((
			COUNT(r1.report_id) - 
			(SELECT COUNT(r2.report_id)
			FROM report r2
			WHERE r2.organization_id = r1.organization_id
			AND r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 1 MONTH), '%Y-%m')
			)
		) / NULLIF(
			(SELECT COUNT(r2.report_id)
			FROM report r2
			WHERE r2.organization_id = r1.organization_id
			AND r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 1 MONTH), '%Y-%m')
			), 0) * 100, 2147483647) AS month_on_month,

		-- 同比：当前月和前一年的比较
		IFNULL((
			COUNT(r1.report_id) - 
			(SELECT COUNT(r2.report_id)
			FROM report r2
			WHERE r2.organization_id = r1.organization_id
			AND r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 12 MONTH), '%Y-%m')
			)
		) / NULLIF(
			(SELECT COUNT(r2.report_id)
			FROM report r2
			WHERE r2.organization_id = r1.organization_id
			AND r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 12 MONTH), '%Y-%m')
			), 0) * 100, 2147483647) AS year_on_year

	FROM 
		report r1
	WHERE 
		r1.organization_id  = ?
		AND r1.tenant_id IN (?) 
		AND r1.is_report_complete = 1
		AND r1.created_at between ? and ?
	GROUP BY 
		r1.organization_id,r1.tenant_id, date
	ORDER BY 
		r1.tenant_id, date
	LIMIT ?,?;
    `, oid, tid, startTime, endTime, offset, size).Scan(&tcs).Error
	if err != nil {
		return nil, 0, err
	}
	// GetConnection
	err = r.GetConnection(ctx).Raw(`
    SELECT COUNT(*) FROM
	(SELECT 
		r1.organization_id,
		r1.tenant_id,
		DATE_FORMAT(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), '%Y-%m') AS date,
		COUNT(r1.report_id) AS monthly_count,
		-- 环比：当前月和前一个月的比较
		IFNULL((
			COUNT(r1.report_id) - 
			(SELECT COUNT(r2.report_id)
			FROM report r2
			WHERE r2.organization_id = r1.organization_id
			AND r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 1 MONTH), '%Y-%m')
			)
		) / NULLIF(
			(SELECT COUNT(r2.report_id)
			FROM report r2
			WHERE r2.organization_id = r1.organization_id
			AND r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 1 MONTH), '%Y-%m')
			), 0) * 100, 2147483647) AS month_on_month,

		-- 同比：当前月和前一年的比较
		IFNULL((
			COUNT(r1.report_id) - 
			(SELECT COUNT(r2.report_id)
			FROM report r2
			WHERE r2.organization_id = r1.organization_id
			AND r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 12 MONTH), '%Y-%m')
			)
		) / NULLIF(
			(SELECT COUNT(r2.report_id)
			FROM report r2
			WHERE r2.organization_id = r1.organization_id
			AND r2.tenant_id = r1.tenant_id
			AND DATE_FORMAT(CONVERT_TZ(r2.created_at, '+00:00', '+08:00'), '%Y-%m') = 
				DATE_FORMAT(DATE_SUB(CONVERT_TZ(r1.created_at, '+00:00', '+08:00'), INTERVAL 12 MONTH), '%Y-%m')
			), 0) * 100, 2147483647) AS year_on_year

	FROM 
		report r1
	WHERE 
		r1.organization_id  = ?
		AND r1.tenant_id IN (?) 
		AND r1.is_report_complete = 1 
		AND r1.created_at between ? and ?
	GROUP BY 
		r1.organization_id,r1.tenant_id, date
	ORDER BY 
		r1.tenant_id, date) a
    `, oid, tid, startTime, endTime).Count(&counts).Error
	if err != nil {
		return nil, 0, err
	}

	dtcs := make([]domain.ReportTenantCountIntf, len(tcs))
	for k, v := range tcs {
		dtcs[k] = v.ToDomainReportTenantCount()
	}
	return dtcs, counts, nil
}

// CreateReportAddress
func (store *ReportStore) CreateReportAddress(ctx context.Context, rd domain.ReportAddressIntf) error {
	var r ReportAddress
	r.FromDomainReportAddress(rd)
	return store.GetConnection(ctx).Model(&ReportAddress{}).Create(&r).Error
}

// GetReportAddress
func (r *ReportStore) GetReportAddress(ctx context.Context, reportId string) (domain.ReportAddressIntf, error) {
	db := r.GetConnection(ctx)
	// 查找报告
	var ad ReportAddress
	err := db.Model(&ReportAddress{}).Where("report_id = ?", reportId).First(&ad).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ad.ToDomainReportAddress(), nil
}

// GetInquiryQuestions
func (r *ReportStore) GetInquiryQuestions(ctx context.Context) ([]domain.InquiryDiagnosisIntf, error) {
	db := r.GetConnection(ctx)
	var ids []InquiryDiagnosis
	err := db.Model(&InquiryDiagnosis{}).Find(&ids).Error
	if err != nil {
		return nil, err
	}
	results := make([]domain.InquiryDiagnosisIntf, len(ids))
	for k, v := range ids {
		results[k] = v.ToDomainInquiryDiagnosis()
	}
	return results, nil
}

// ModifyReportInquiryDiagnosis
func (r *ReportStore) ModifyReportInquiryDiagnosis(ctx context.Context, rid, info string) error {
	db := r.GetConnection(ctx)
	return db.Model(&Report{}).Where("report_id = ?", rid).Update("inquiry_diagnosis", info).Error
}
