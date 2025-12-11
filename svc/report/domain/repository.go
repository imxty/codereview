package domain

import (
	"context"
	"time"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
)

type ReportRepository interface {
	dbutils.Tx
	// ListCustomerLastReports 获取常客测量的最后limit份报告
	ListCustomerLastReports(ctx context.Context, tenantID, customerID string, limit int) ([]ReportIntf, error)

	// GetReportIDByToken 通过报告的token获取报告ID
	GetReportIDByToken(ctx context.Context, tenantID, token string) (SharedReportIntf, error)
	// GetReportByID 通过报告ID获取报告
	GetReportByID(ctx context.Context, tenantID, reportID string) (ReportIntf, error)
	// GetReportByReportID 通过报告ID获取报告
	GetReportByReportID(ctx context.Context, reportID string) (ReportIntf, error)
	// UpdateReportTongueFaceReportStatus 更新报告的舌面分析报告状态
	UpdateReportTongueFaceReportStatus(ctx context.Context, faceUrl, tongueUrl, reportID, body string, status int32) error
	// GetCustomerReport 获取常客报告
	GetCustomerReport(ctx context.Context, tenantID, reportID, customerID string) (ReportIntf, error)
	// CreateSharedReport 创建分享报告
	CreateSharedReport(ctx context.Context, sharedReport SharedReportIntf) error
	// ListReports 获取常客/散客报告列表
	ListReports(ctx context.Context, tenantID, customerID string, startTime, endTime time.Time, offset, size int32) ([]ReportIntf, int64, error)
	// ListCustomerReports 获取常客报告列表
	ListCustomerReports(ctx context.Context, tenantId, customerId string, startTime, endTime time.Time) ([]ReportIntf, error)

	// ListReportsWithStaff 获取员工常客/散客报告列表
	ListReportsWithStaff(ctx context.Context, tenantID, customerID, staffId string, reportSharingStatus bool, startTime, endTime time.Time, offset, size int32) ([]ReportIntf, int64, error)

	// ListReportsWithoutTime 获取常客/散客报告列表(没有时间区间)
	ListReportsWithoutTime(ctx context.Context, tenantID, customerID string, offset, size int32) ([]ReportIntf, int64, error)
	// ModifyReportRemark 修改报告的普通备注
	ModifyReportRemark(ctx context.Context, tenantID, reportID, remark string, rev int32) error
	// ModifyReportStaffRemark 修改报告的员工备注
	ModifyReportStaffRemark(ctx context.Context, tenantID, reportID, remark string, rev int32) error
	// CreateReport 创建报告
	CreateReport(ctx context.Context, report ReportIntf) error
	// AddCustomerReport 散客报告变为常客报告
	AddCustomerReport(ctx context.Context, tenantID, reportId, customerID, age string, rev int32) error

	// UpdateReportContent 更新报告的内容
	UpdateReportContent(ctx context.Context, reportID, tenantID string, heartRate int32, isStressState bool, factor *FactorInterpretation, stressStateJson, dirtyDialectJson, diseasesJson, recommendJson, physiqueJson, measurementJson string, chart *MeridianBarChart, dirtyDialectFlag, rev int32) error
	// UpdateReportTreatment 更新报告方案
	UpdateReportTreatment(ctx context.Context, tenantID, reportID, tenantTreatmentRevId string, rev int32) error

	// GetOrganizationReportCount 获取报告统计(当日新增常客测量次数，当日新增散客测量次数，常客累计测量次数,散客累计测量次数)
	GetOrganizationReportCount(ctx context.Context, organizationID string, startTime, endTime time.Time) (int64, int64, int64, int64, error)
	// GetTenantReportCount 获取报告统计(当日新增常客测量次数，当日新增散客测量次数，常客累计测量次数,散客累计测量次数)
	GetTenantReportCount(ctx context.Context, tenantID string, startTime, endTime time.Time) (int64, int64, int64, int64, error)
	// GetStaffMeasurementCount 获取员工测量统计
	GetStaffMeasurementCount(ctx context.Context, tenantID string, staffIds []string, startTime, endTime time.Time) ([]StaffReportRankIntf, error)

	// ListTenantsReportCount 获取商户报告数量
	ListTenantsReportCount(ctx context.Context, tenantID []string, startTime, endTime time.Time) ([]ReportRankIntf, error)

	// GetTenantReportsCount 获取时间范围内商户报告数量
	GetTenantReportsCount(ctx context.Context, tenantID string, startTime, endTime time.Time) (int64, error)
	// GetOrganizationReportsCount 获取时间范围内商户报告数量
	GetOrganizationReportsCount(ctx context.Context, organizationID string, startTime, endTime time.Time) (int64, error)
	// SearchReports 查询报告
	SearchReports(ctx context.Context, reportId string, customerType int, customerId string, tenantId []string, dirtyDialecticsFlag int32, organizationId string, startTime, endTime time.Time, offset, size int32) ([]ReportIntf, int32, error)

	// 查询组织测量报告数量(常客测量数量，散客测量数量，error)
	SearchOrganizationReportsCount(ctx context.Context, oid string, tids []string, startTime, endTime time.Time) (int64, int64, error)
	// SearchOrganizationReports
	SearchOrganizationReports(ctx context.Context, oid string, tids []string, startTime, endTime time.Time, offset, size int32) ([]ReportCountIntf, int64, error)

	// 查询组织测量报告数量(常客测量数量，散客测量数量，error)
	SearchStaffReportsCount(ctx context.Context, tid string, sid []string, startTime, endTime time.Time) (int64, int64, error)
	// SearchStaffReports
	SearchStaffReports(ctx context.Context, tid string, sid []string, startTime, endTime time.Time, offset, size int32) ([]ReportStaffCountIntf, int64, error)

	// SearchTenantReports
	SearchTenantReports(ctx context.Context, oid string, tid []string, startTime, endTime time.Time, offset, size int32) ([]ReportTenantCountIntf, int64, error)
	// SearchTenantReportsDownload
	SearchTenantReportsDownload(ctx context.Context, oid string, tid []string, startTime, endTime time.Time) ([]ReportTenantCountIntf, int64, error)
	// CreateReportAddress
	CreateReportAddress(ctx context.Context, rd ReportAddressIntf) error
	// GetReportAddress
	GetReportAddress(ctx context.Context, reportId string) (ReportAddressIntf, error)
	// GetInquiryQuestions 获取问诊问题
	GetInquiryQuestions(ctx context.Context) ([]InquiryDiagnosisIntf, error)
	// ModifyReportInquiryDiagnosis 修改报告问诊信息
	ModifyReportInquiryDiagnosis(ctx context.Context, rid, info string) error
}
