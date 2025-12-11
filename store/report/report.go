package report

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
	"gorm.io/gorm"
)

const (
	CustomerTypeBoth     = 0
	CustomerTypeCustomer = 1
	CustomerTypeTemp     = 2
)

// Report 报告内容
type Report struct {
	ReportID               string          `gorm:"primary_key;column:report_id"`
	TenantID               string          `gorm:"column:tenant_id"`
	OrganizationID         string          `gorm:"column:organization_id"`
	StaffID                string          `gorm:"column:staff_id"`
	DeviceID               string          `gorm:"column:device_id"`
	IsCustomer             bool            `gorm:"column:is_customer"`
	CustomerID             string          `gorm:"column:customer_id"`
	Age                    string          `gorm:"column:age"`
	Gender                 int32           `gorm:"column:gender"`
	Hand                   int32           `gorm:"column:hand"`
	Remarks                string          `gorm:"column:remarks"`
	StaffRemarks           string          `gorm:"column:staff_remarks"`
	HeartRate              int32           `gorm:"column:heart_rate"`
	Spo                    int32           `gorm:"column:spo"`
	Risks                  string          `gorm:"column:risks"`
	C0                     int32           `gorm:"column:c0"`
	C1                     int32           `gorm:"column:c1"`
	C2                     int32           `gorm:"column:c2"`
	C3                     int32           `gorm:"column:c3"`
	C4                     int32           `gorm:"column:c4"`
	C5                     int32           `gorm:"column:c5"`
	C6                     int32           `gorm:"column:c6"`
	C7                     int32           `gorm:"column:c7"`
	PhysicalDialectics     string          `gorm:"column:physical_dialectics"`
	DirtyDialectic         string          `gorm:"column:dirty_dialectic"`
	DirtyDialecticFlag     int32           `gorm:"column:dirty_dialectic_flag"`
	DietaryAdvice          string          `gorm:"column:dietary_advice"`
	MeasurementJudgment    string          `gorm:"column:measurement_judgment"`
	F0                     int32           `gorm:"column:f0"`
	F1                     int32           `gorm:"column:f1"`
	F2                     int32           `gorm:"column:f2"`
	F3                     int32           `gorm:"column:f3"`
	F4                     int32           `gorm:"column:f4"`
	IsStressState          bool            `gorm:"column:is_stress_state"`
	IsReportComplete       bool            `gorm:"column:is_report_complete"`
	StressState            string          `gorm:"column:stress_state"`
	TenantTreatmentRevID   string          `gorm:"column:tenant_treatment_rev_id"`
	FaceImageURL           string          `gorm:"column:face_image_url"`
	TongueImageURL         string          `gorm:"column:tongue_image_url"`
	TongueFaceReportStatus int32           `gorm:"column:tongue_face_report_status"`
	TongueFaceReport       string          `gorm:"column:tongue_face_report"`
	Rev                    int32           `gorm:"column:rev"`
	CreatedAt              time.Time       // 创建时间
	UpdatedAt              time.Time       // 更新时间
	DeletedAt              *gorm.DeletedAt // 删除时间
}

func (r Report) TableName() string {
	return "report"
}

// domain->db
func (r *Report) FromDomainReport(d domain.ReportIntf) {
	if r == nil || d == nil {
		return
	}

	r.ReportID = d.GetReportID()
	r.TenantID = d.GetTenantID()
	r.OrganizationID = d.GetOrganizationID()
	r.StaffID = d.GetStaffID()
	r.DeviceID = d.GetDeviceID()
	r.IsCustomer = d.GetIsCustomer()
	r.CustomerID = d.GetCustomerID()
	r.Age = d.GetAge()
	r.Gender = d.GetGender()
	r.Hand = d.GetHand()
	r.Remarks = d.GetRemarks()
	r.StaffRemarks = d.GetStaffRemarks()
	r.HeartRate = d.GetHeartRate()
	r.Spo = d.GetSpo()
	r.Risks = d.GetRisks()
	r.C0 = d.GetC0()
	r.C1 = d.GetC1()
	r.C2 = d.GetC2()
	r.C3 = d.GetC3()
	r.C4 = d.GetC4()
	r.C5 = d.GetC5()
	r.C6 = d.GetC6()
	r.C7 = d.GetC7()
	r.PhysicalDialectics = d.GetPhysicalDialectics()
	r.DirtyDialectic = d.GetDirtyDialectic()
	r.DirtyDialecticFlag = d.GetDirtyDialecticFlag()
	r.DietaryAdvice = d.GetDietaryAdvice()
	r.MeasurementJudgment = d.GetMeasurementJudgment()
	r.F0 = d.GetF0()
	r.F1 = d.GetF1()
	r.F2 = d.GetF2()
	r.F3 = d.GetF3()
	r.F4 = d.GetF4()
	r.IsStressState = d.GetIsStressState()
	r.IsReportComplete = d.GetIsReportComplete()
	r.StressState = d.GetStressState()
	r.TenantTreatmentRevID = d.GetTenantTreatmentRevID()
	r.FaceImageURL = d.GetFaceImageURL()
	r.TongueImageURL = d.GetTongueImageURL()
	r.TongueFaceReportStatus = d.GetTongueFaceReportStatus()
	r.TongueFaceReport = d.GetTongueFaceReport()
	r.Rev = d.GetRev()
}

// db->domain
func (r *Report) ToDomainReport() domain.ReportIntf {
	if r == nil {
		return nil
	}

	p := domain.Report{
		ReportID:               r.ReportID,
		TenantID:               r.TenantID,
		OrganizationID:         r.OrganizationID,
		StaffID:                r.StaffID,
		DeviceID:               r.DeviceID,
		IsCustomer:             r.IsCustomer,
		CustomerID:             r.CustomerID,
		Age:                    r.Age,
		Gender:                 r.Gender,
		Hand:                   r.Hand,
		Remarks:                r.Remarks,
		StaffRemarks:           r.StaffRemarks,
		HeartRate:              r.HeartRate,
		Spo:                    r.Spo,
		Risks:                  r.Risks,
		C0:                     r.C0,
		C1:                     r.C1,
		C2:                     r.C2,
		C3:                     r.C3,
		C4:                     r.C4,
		C5:                     r.C5,
		C6:                     r.C6,
		C7:                     r.C7,
		PhysicalDialectics:     r.PhysicalDialectics,
		DirtyDialectic:         r.DirtyDialectic,
		DirtyDialecticFlag:     r.DirtyDialecticFlag,
		DietaryAdvice:          r.DietaryAdvice,
		MeasurementJudgment:    r.MeasurementJudgment,
		F0:                     r.F0,
		F1:                     r.F1,
		F2:                     r.F2,
		F3:                     r.F3,
		F4:                     r.F4,
		IsStressState:          r.IsStressState,
		IsReportComplete:       r.IsReportComplete,
		StressState:            r.StressState,
		TenantTreatmentRevID:   r.TenantTreatmentRevID,
		FaceImageURL:           r.FaceImageURL,
		TongueImageURL:         r.TongueImageURL,
		TongueFaceReportStatus: r.TongueFaceReportStatus,
		TongueFaceReport:       r.TongueFaceReport,
		Rev:                    r.Rev,
		CreatedAt:              r.CreatedAt,
	}
	return &p
}
